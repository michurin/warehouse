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

func (h *House) RoomOrNil(roomID string) (*wall.Wall, *user.Users) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	r, ok := h.rooms[roomID]
	if ok {
		return r.wall, r.users
	}
	return nil, nil
}

func (h *House) Room(roomID string) (*wall.Wall, *user.Users) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[roomID]
	if ok {
		return r.wall, r.users
	}
	// TODO check len(h.rooms), can we add one more room
	users := user.New() // we will add current user on caller side
	wall := wall.New(time.Now().UnixNano())
	h.rooms[roomID] = &room{users: users, wall: wall}
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

func (h *House) Audit(ms int64) ([]*wall.Wall, []*user.Users) {
	uu := []*room(nil)
	h.mu.RLock()
	for _, v := range h.rooms {
		uu = append(uu, v)
	}
	h.mu.RUnlock()
	walls := []*wall.Wall(nil)
	users := []*user.Users(nil)
	for _, u := range uu {
		if u.users.Audit(ms) {
			walls = append(walls, u.wall)
			users = append(users, u.users)
		}
	}
	return walls, users // TODO remove parallel arrays
}
