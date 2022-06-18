package stream

import (
	"context"
	"sync"
)

type Stream struct {
	notifier chan (struct{})
	mx       *sync.RWMutex
	capacity int
	head     int // increasing infinitely; (head % capacity) points to newest element
	messages [][]byte
}

func New(capacity int) *Stream {
	return &Stream{
		notifier: make(chan struct{}),
		mx:       new(sync.RWMutex),
		capacity: capacity,
		head:     -1, // 0-1
		messages: make([][]byte, capacity),
	}
}

func (s *Stream) Put(x []byte) {
	s.mx.Lock()
	s.head++
	s.messages[s.head%s.capacity] = x
	n := s.notifier
	s.notifier = make(chan struct{})
	s.mx.Unlock()
	close(n)
}

func (s *Stream) Get(ctx context.Context, bound int) ([][]byte, int) {
	w, t, h := s.take(bound)
	if len(t) > 0 {
		return t, h
	}
	select {
	case <-ctx.Done():
		return nil, h // len=0, h=bound?
	case <-w:
		_, t, h = s.take(bound)
		return t, h
	}
}

func (s *Stream) take(bound int) (chan struct{}, [][]byte, int) {
	s.mx.RLock()
	h := s.head
	b := h - s.capacity
	if b < 0 {
		b = -1
	}
	if bound <= s.head && b < bound { // consider bound if it is valid only
		b = bound
	}
	r := [][]byte(nil)
	for i := s.head; i > b; i-- {
		r = append(r, s.messages[i%s.capacity])
	}
	n := s.notifier
	s.mx.RUnlock()
	return n, r, h
}
