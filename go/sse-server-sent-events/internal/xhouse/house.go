package xhouse

import (
	"sync"
	"time"

	"github.com/michurin/minchat/internal/xuser"
	"github.com/michurin/minchat/internal/xwall"
)

type room struct {
	users *xuser.Users
	wall  *xwall.Wall
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

func (h *House) RoomOrNil(roomID string) (*xwall.Wall, *xuser.Users) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	r, ok := h.rooms[roomID]
	if ok {
		return r.wall, r.users
	}
	return nil, nil
}

func (h *House) Room(roomID string) (*xwall.Wall, *xuser.Users) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[roomID]
	if ok {
		return r.wall, r.users
	}
	// TODO check len(h.rooms), can we add one more room
	users := xuser.New() // we will add current user on caller side
	wall := xwall.New(time.Now().UnixNano())
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

func (h *House) Audit(ms int64) ([]*xwall.Wall, []*xuser.Users) {
	uu := []*room(nil)
	h.mu.RLock()
	for _, v := range h.rooms {
		uu = append(uu, v)
	}
	h.mu.RUnlock()
	walls := []*xwall.Wall(nil)
	users := []*xuser.Users(nil)
	for _, u := range uu {
		if u.users.Audit(ms) {
			walls = append(walls, u.wall)
			users = append(users, u.users)
		}
	}
	return walls, users // TODO remove parallel arrays
}
