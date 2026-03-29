package room

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"sse/dto"
	"sse/user"
	"sse/wall"
)

type Room struct {
	users  *user.Users
	locked bool
	mu     *sync.Mutex
	wall   *wall.Wall
}

func (r *Room) Pub(userID string, message []byte) {
	r.users.Update(userID)
	r.wall.Pub(message)
}

func (r *Room) Fetch(ctx context.Context, user string, lastID int64) ([][]byte, int64) {
	// TODO check user if locked
	// TODO manage users list and connections counters
	// TODO user.Inc defer user.Dec
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
		// TODO check len, can we add one more room
		u := user.New()
		// TODO and current user
		w := wall.New(time.Now().UnixNano())
		b, err := json.Marshal(dto.StreamMessage{
			RoomStatus: &dto.RoomStatus{
				Locked: false,
				Users:  []string{}, // get from u?
			},
		})
		if err != nil {
			panic(err) // it is impossible
		}
		w.Pub(b)
		r = &Room{
			users:  u,
			locked: false,
			mu:     new(sync.Mutex),
			wall:   w,
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
