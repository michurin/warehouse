package chat

import (
	"container/list"
	"context"
	"sync"
)

type Message struct {
	Text string `json:"text"`
}

type slot struct {
	id      int
	message Message
}

type Storage struct {
	lastID    int
	bank      *list.List
	dropPoint *list.Element
	lock      *sync.Mutex
	signal    chan struct{}
}

func New() *Storage {
	return &Storage{
		bank:   list.New(),
		lock:   &sync.Mutex{},
		signal: make(chan struct{}),
	}
}

func (s *Storage) Add(message Message) { // TODO return registered-as id
	s.lock.Lock()
	defer s.lock.Unlock()
	s.lastID++
	s.bank.PushFront(slot{s.lastID, message})
	for s.bank.Len() > 10 { // TODO turn 10->s.capacity
		s.bank.Remove(s.bank.Back())
	}
	close(s.signal)
	s.signal = make(chan struct{})
}

func (s *Storage) get(lastID int) ([]Message, chan struct{}, int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	res := []Message(nil)
	actualLastID := 0
	for e := s.bank.Front(); e != nil; e = e.Next() {
		s := e.Value.(slot)
		if s.id <= lastID {
			break
		}
		res = append(res, s.message)
	}
	actualLastID = s.lastID
	if res == nil {
		return nil, s.signal, actualLastID
	}
	return res, nil, actualLastID
}

func (s *Storage) Get(ctx context.Context, lastID int) ([]Message, int) {
	messages, c, actualLastID := s.get(lastID)
	if c != nil {
		select {
		case <-ctx.Done():
		case <-c:
			messages, _, actualLastID = s.get(lastID)
		}
	}
	return messages, actualLastID
}
