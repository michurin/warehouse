package user

import (
	"log"
	"sync"
)

type user struct {
	name      string // TODO pack color+nik as []byte. Like wall.Wall do?
	color     string
	lastCheck int64
}

type Users struct {
	users  map[string]*user
	locked bool
	mu     *sync.RWMutex
}

func New() *Users {
	return &Users{
		users:  map[string]*user{},
		locked: false,
		mu:     new(sync.RWMutex),
	}
}

func (u *Users) Touch(userID string, ms int64, name, color string) (bool, bool) {
	u.mu.Lock()
	defer u.mu.Unlock()
	r, ok := u.users[userID]
	if !ok {
		if u.locked {
			return false, false // not allowed, no updates
		}
		u.users[userID] = &user{
			name:      name,
			color:     color,
			lastCheck: ms,
		}
		return true, true // allowed, update
	}
	r.lastCheck = ms
	if name == "" && color == "" {
		return true, false // allowed, no update it is fetch for already existed user
	}
	if r.name == name && r.color == color {
		return true, false // allowed, no update
	}
	r.color = color
	r.name = name
	return true, true // allowed, updated
}

func (u *Users) Lock(v bool) bool {
	u.mu.Lock()
	defer u.mu.Unlock()
	log.Printf("LOCK %v->%v", u.locked, v)
	if u.locked == v {
		return false
	}
	u.locked = v
	return true
}

func (u *Users) Locked() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.locked
}

func (u *Users) List() [][2]string {
	u.mu.RLock()
	defer u.mu.RUnlock()
	r := make([][2]string, 0, len(u.users))
	for _, v := range u.users {
		r = append(r, [2]string{v.name, v.color})
	}
	return r
}

func (u *Users) Get(userID string) (string, string) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	x, ok := u.users[userID]
	if !ok {
		return "", ""
	}
	return x.name, x.color
}
