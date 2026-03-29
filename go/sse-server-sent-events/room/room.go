package room

import (
	"context"
	"sse/wall"
	"sync"
	"time"
)

type Room struct {
	users  map[string]struct{} // TODO: name and count of connections
	locked bool
	mu     *sync.Mutex
	wall   *wall.Wall
}

func (r *Room) Pub(user string, message []byte) {
	// TODO check user if locked
	r.wall.Pub(message)
}

func (r *Room) Fetch(ctx context.Context, user string, lastID int64) ([][]byte, int64) {
	// TODO check user if locked
	// TODO manage users list and connections counters
	return r.wall.Fetch(ctx, lastID)
}

type House struct {
	rooms map[string]*Room
	mu    *sync.RWMutex
}

func New() *House {
	return &House{
		rooms: map[string]*Room{},
		mu:    new(sync.RWMutex), // TODO we need just Mutex?
	}
}

func (h *House) room(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[name]
	if !ok {
		// TODO check len
		r = &Room{
			users:  map[string]struct{}{},
			locked: false,
			mu:     new(sync.Mutex),
			wall:   wall.New(time.Now().UnixNano()),
		}
		h.rooms[name] = r
	}
	return r
}

func (h *House) Pub(room, user string, message []byte) {
	h.room(room).Pub(user, message)
}

func (h *House) Fetch(ctx context.Context, room, user string, lastID int64) ([][]byte, int64) {
	return h.room(room).Fetch(ctx, user, lastID)
}

// Lock(room, user), lock and send broadcast notification
// Unlock(room, user)
