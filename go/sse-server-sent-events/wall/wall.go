package wall

import (
	"container/list"
	"context"
	"sync"
)

type Wall struct {
	lastID  int64
	wall    *list.List
	unblock chan struct{}
	mu      *sync.RWMutex
}

func New(initialShift int64) *Wall {
	return &Wall{
		lastID:  initialShift,
		wall:    list.New(),
		unblock: make(chan struct{}),
		mu:      new(sync.RWMutex),
	}
}

func (r *Wall) Pub(m []byte) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastID++
	r.wall.PushFront(m)
	for r.wall.Len() > 1000 {
		r.wall.Remove(r.wall.Back())
	}
	close(r.unblock)
	r.unblock = make(chan struct{})
}

func (r *Wall) Fetch(ctx context.Context, lastID int64) ([][]byte, int64) {
	messages, unblock, id := r.syncFetch(lastID)
	if len(messages) > 0 {
		return messages, id
	}
	select {
	case <-ctx.Done():
		return nil, lastID
	case <-unblock:
		messages, _, id = r.syncFetch(lastID)
		return messages, id
	}
}

func (r *Wall) syncFetch(lastID int64) ([][]byte, chan struct{}, int64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	wall := [][]byte(nil)
	i := r.lastID
	for e := r.wall.Front(); e != nil; e = e.Next() {
		if i <= lastID {
			break
		}
		wall = append(wall, e.Value.([]byte))
		i--
	}
	return wall, r.unblock, r.lastID
}
