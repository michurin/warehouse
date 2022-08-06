package stream

import (
	"context"
	"sync"
)

// reusable closed channel
var closedChan = make(chan struct{})

func init() { close(closedChan) }

// The head (and therefore bound too) is uint64,
// however keep in mind that JavaScript
// Number.MAX_SAFE_INTEGER = 2**53-1
type Stream struct {
	notifier chan (struct{})
	mx       *sync.RWMutex
	capacity uint64
	head     uint64 // increasing infinitely; (head % capacity) points to next slot to write
	messages [][]byte
}

func New(capacity uint64) *Stream {
	return &Stream{
		notifier: make(chan struct{}),
		mx:       new(sync.RWMutex),
		capacity: capacity,
		head:     0,
		messages: make([][]byte, capacity),
	}
}

// Put puts data to storage and unlock all waiting Get calls.
// You MUST NOT to mutate input slice after putting it.
func (s *Stream) Put(x []byte) {
	s.mx.Lock()
	s.messages[s.head%s.capacity] = x
	s.head++
	n := s.notifier
	s.notifier = make(chan struct{})
	s.mx.Unlock()
	close(n)
}

// Get obtains bound and returns data and new bound.
// It is literally Waiter+Updates.
// You man want to use Waiter and Updates separately if
// you have several streams.
func (s *Stream) Get(ctx context.Context, bound uint64) ([][]byte, uint64) {
	select {
	case <-ctx.Done():
		return nil, bound
	case <-s.Waiter(bound):
		return s.Updates(bound)
	}
}

// Waiter returns channel that is getting closed when new
// messages appears. Method Updates returns nonempty result
// after waiter channel gets closes.
func (s *Stream) Waiter(bound uint64) (c chan struct{}) {
	s.mx.RLock()
	if s.head == bound {
		c = s.notifier
	} else {
		c = closedChan
	}
	s.mx.RUnlock()
	return
}

// Updates returns all messages older then bound and new bound.
func (s *Stream) Updates(bound uint64) ([][]byte, uint64) {
	s.mx.RLock()
	h := s.head
	if h < bound { // bound from previous run; server was restarted
		bound = 0
	}
	l := h - bound
	if l > s.capacity {
		l = s.capacity
	}
	r := make([][]byte, l)
	for i := uint64(0); i < l; i++ {
		r[i] = s.messages[(s.head-l+i)%s.capacity]
	}
	s.mx.RUnlock()
	return r, h
}
