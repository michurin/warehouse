package stream

import (
	"context"
	"sync"
)

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
func (s *Stream) Put(x []byte) {
	s.mx.Lock()
	s.messages[s.head%s.capacity] = x
	s.head++
	n := s.notifier
	s.notifier = make(chan struct{})
	s.mx.Unlock()
	close(n)
}

// Get obtains bound and returns data, new bound and continuity flag
// If there is no new data the method is waiting for for it or for context.
// The bound is uint64, however keep in mind that JavaScript
// Number.MAX_SAFE_INTEGER = 2**53-1
// Continuity flag is false if
// - server restart detected
// - buffer of messages was exhausted before bound was reached
// The only exception: continuity flag is always true for new clients.
func (s *Stream) Get(ctx context.Context, bound uint64) ([][]byte, uint64, bool) {
	w, t, h, c := s.take(bound)
	if len(t) > 0 {
		return t, h, c
	}
	select {
	case <-ctx.Done():
		return nil, h, true
	case <-w:
		_, t, h, c = s.take(bound)
		return t, h, c
	}
}

func (s *Stream) take(bound uint64) (chan struct{}, [][]byte, uint64, bool) {
	fromFuture := false
	cropped := false
	newSession := bound == 0
	s.mx.RLock()
	h := s.head
	if h < bound { // bound from previous run; server was restarted
		bound = 0
		fromFuture = true
	}
	l := h - bound
	if l > s.capacity {
		l = s.capacity
		cropped = true
	}
	r := make([][]byte, l)
	for i := uint64(0); i < l; i++ {
		r[i] = s.messages[(s.head-l+i)%s.capacity]
	}
	n := s.notifier
	s.mx.RUnlock()
	return n, r, h, newSession || !(cropped || fromFuture)
}
