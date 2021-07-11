package chat

import (
	"container/list"
	"context"
	"encoding/json"
	"sync"
	"time"
)

func timeBasedIDIniter() int64 {
	// We assume, we will have less than 10K messages per second
	// and restarts will take more than one second
	return time.Now().Unix() * 10000
}

const (
	maxCountDown    = 10  // TODO config
	backlogCapacity = 100 // TODO config
)

type slot struct {
	id      int64
	message json.RawMessage
}

type oneRoom struct {
	lastID    int64
	bank      *list.List
	coundDown int
	lock      *sync.Mutex // TODO RWMutex?
	signal    chan struct{}
	idIniter  func() int64
}

func newRoom(idIniter func() int64) *oneRoom {
	return &oneRoom{
		lastID:    0, // lazy lastID initialization
		bank:      list.New(),
		coundDown: maxCountDown,
		lock:      new(sync.Mutex),
		signal:    make(chan struct{}),
		idIniter:  idIniter,
	}
}

func (r *oneRoom) pub(message json.RawMessage) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// update counter, chat is touched
	r.coundDown = maxCountDown

	// lazy initialization
	if r.lastID == 0 {
		r.lastID = r.idIniter()
	}

	// publish
	r.lastID++
	r.bank.PushFront(slot{
		id:      r.lastID,
		message: message,
	})

	// cut the tail
	for r.bank.Len() > backlogCapacity {
		r.bank.Remove(r.bank.Back())
	}

	// notify all reader
	close(r.signal)
	r.signal = make(chan struct{})
}

func (r *oneRoom) get(lastID int64) ([]json.RawMessage, chan struct{}, int64) {
	r.lock.Lock() // TODO RLock?
	defer r.lock.Unlock()
	res := []json.RawMessage(nil)
	for e := r.bank.Front(); e != nil; e = e.Next() {
		s := e.Value.(slot)
		if s.id <= lastID {
			break
		}
		res = append(res, s.message)
	}
	if res == nil {
		return nil, r.signal, lastID // just echo lastID cause we don't have any updates
	}
	return res, nil, r.lastID
}

func (r *oneRoom) fetch(ctx context.Context, lastID int64) ([]json.RawMessage, int64) {
	messages, c, actualLastID := r.get(lastID)
	if c != nil {
		select {
		case <-ctx.Done():
		case <-c:
			messages, _, actualLastID = r.get(lastID)
		}
	}
	return messages, actualLastID
}

func (s *oneRoom) tick() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.coundDown--
	return s.coundDown < 0
}

type Logger interface {
	Log(ctx context.Context, message ...interface{})
}

type Rooms struct {
	rooms  *sync.Map
	logger Logger
}

func New(l Logger) *Rooms {
	if l == nil { // TODO default logger
		panic("Logger must to be not nil")
	}
	return &Rooms{
		rooms:  new(sync.Map),
		logger: l,
	}
}

func (r *Rooms) room(ctx context.Context, rid string) *oneRoom {
	si, loaded := r.rooms.LoadOrStore(rid, newRoom(timeBasedIDIniter))
	if !loaded {
		r.logger.Log(ctx, "New room created:", rid)
	}
	return si.(*oneRoom)
}

func (r *Rooms) Pub(ctx context.Context, rid string, msg json.RawMessage) {
	r.room(ctx, rid).pub(msg)
}

func (r *Rooms) Fetch(ctx context.Context, rid string, lastID int64) ([]json.RawMessage, int64) {
	return r.room(ctx, rid).fetch(ctx, lastID)
}

func RoomsCleaner(ctx context.Context, r *Rooms) {
	ticker := time.NewTicker(time.Minute) // it has to be longer then polling interval
L:
	for {
		select {
		case <-ticker.C:
			count := 0
			r.rooms.Range(func(k, v interface{}) bool {
				if v.(*oneRoom).tick() {
					r.logger.Log(ctx, "Delete", k)
					r.rooms.Delete(k)
				}
				count++
				return true
			})
			r.logger.Log(ctx, "Tick: ~", count, "rooms")
		case <-ctx.Done():
			break L
		}
	}
	r.logger.Log(ctx, "Ticker stopping")
	ticker.Stop()
}

func RoomsList(r *Rooms) map[string][]json.RawMessage {
	lst := map[string][]json.RawMessage{}
	r.rooms.Range(func(k, v interface{}) bool {
		m, _, _ := v.(*oneRoom).get(0)
		lst[k.(string)] = m
		return true
	})
	return lst
}
