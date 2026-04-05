package room

import (
	"sync"
	"time"

	"sse/user"
	"sse/wall"
)

type room struct {
	users *user.Users
	wall  *wall.Wall
}

type House struct {
	rooms map[string]*room
	mu    *sync.RWMutex
}

func New() *House {
	return &House{
		rooms: map[string]*room{},
		mu:    new(sync.RWMutex),
	}
}

func (h *House) RoomOrNil(room string) (*wall.Wall, *user.Users) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	r, ok := h.rooms[room]
	if ok {
		return r.wall, r.users
	}
	return nil, nil
}

func (h *House) Room(name string) (*wall.Wall, *user.Users) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[name]
	if ok {
		return r.wall, r.users
	}
	// TODO check len(h.rooms), can we add one more room
	users := user.New() // we will add current user on caller side
	wall := wall.New(time.Now().UnixNano())
	h.rooms[name] = &room{users: users, wall: wall}
	return wall, users
}

func (h *House) List() []string { // for debugging only
	h.mu.RLock()
	defer h.mu.RUnlock()
	r := []string(nil)
	for k := range h.rooms {
		r = append(r, k)
	}
	return r
}
