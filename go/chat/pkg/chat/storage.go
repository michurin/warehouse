package chat

import (
	"container/list"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/michurin/minlog"
)

func TimeBasedIDIniter() int64 {
	// We assume, we will have less than 10K messages per second
	// and restarts will take more than one second
	return time.Now().Unix() * 10000
}

const (
	maxCountDown    = 10  // TODO config
	backlogCapacity = 100 // TODO config
)

type Message struct {
	Message json.RawMessage `json:"message"`
}

type slot struct {
	id      int64
	message json.RawMessage
}

type Storage struct { // TODO rename, make it local
	lastID    int64
	bank      *list.List
	coundDown int
	lock      *sync.Mutex // TODO RWMutex?
	signal    chan struct{}
	idIniter  func() int64
}

func New(idIniter func() int64) *Storage {
	return &Storage{
		lastID:    0, // lazy lastID initialization
		bank:      list.New(),
		coundDown: maxCountDown,
		lock:      new(sync.Mutex),
		signal:    make(chan struct{}),
		idIniter:  idIniter,
	}
}

func (s *Storage) Put(message json.RawMessage) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// update counter, chat is touched
	s.coundDown = maxCountDown

	// lazy initialization
	if s.lastID == 0 {
		s.lastID = s.idIniter()
	}

	// publish
	s.lastID++
	s.bank.PushFront(slot{
		id:      s.lastID,
		message: message,
	})

	// cut the tail
	for s.bank.Len() > backlogCapacity {
		s.bank.Remove(s.bank.Back())
	}

	// notify all reader
	close(s.signal)
	s.signal = make(chan struct{})
}

func (s *Storage) get(lastID int64) ([]json.RawMessage, chan struct{}, int64) {
	s.lock.Lock() // TODO RLock?
	defer s.lock.Unlock()
	res := []json.RawMessage(nil)
	for e := s.bank.Front(); e != nil; e = e.Next() {
		s := e.Value.(slot)
		if s.id <= lastID {
			break
		}
		res = append(res, s.message)
	}
	if res == nil {
		return nil, s.signal, lastID // just echo lastID cause we don't have any updates
	}
	return res, nil, s.lastID
}

func (s *Storage) Get(ctx context.Context, lastID int64) ([]json.RawMessage, int64) {
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

func (s *Storage) Tick() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.coundDown--
	return s.coundDown < 0
}

// ----- TODO move to another file?

type Rooms struct {
	Rooms *sync.Map // TODO private
}

func (r *Rooms) room(ctx context.Context, rid string) *Storage {
	si, loaded := r.Rooms.LoadOrStore(rid, New(TimeBasedIDIniter))
	if !loaded {
		minlog.Log(ctx, "New room created:", rid)
	}
	return si.(*Storage)
}

func (r *Rooms) Pub(ctx context.Context, rid string, msg json.RawMessage) {
	r.room(ctx, rid).Put(msg)
}

func (r *Rooms) Fetch(ctx context.Context, rid string, lastID int64) ([]json.RawMessage, int64) {
	return r.room(ctx, rid).Get(ctx, lastID)
}

func RoomCleaner(ctx context.Context, r *Rooms) {
	// TODO move to "constructor"?
	// TODO obtain context.Context and finish on cancel
	// TODO instrumentation?
	ticker := time.NewTicker(time.Minute) // it has to be longer then polling interval
	go func() {
	L:
		for {
			select {
			case <-ticker.C:
			case <-ctx.Done():
				break L
			}
			count := 0
			r.Rooms.Range(func(k, v interface{}) bool {
				if v.(*Storage).Tick() {
					minlog.Log(ctx, "Delete", k)
					r.Rooms.Delete(k)
				}
				count++
				return true
			})
			minlog.Log(ctx, "Tick: ~", count, "rooms")
		}
		ticker.Stop()
	}()
}
