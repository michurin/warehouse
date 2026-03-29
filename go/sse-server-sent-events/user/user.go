package user

import (
	"sync"
)

type User struct {
	nik       string // TODO
	lastCheck int64
	mu        *sync.RWMutex
}

type Users struct {
	users  map[string]*User
	locked bool // TODO
	mu     *sync.RWMutex
}

func New() *Users {
	return &Users{
		users:  map[string]*User{},
		locked: false,
		mu:     new(sync.RWMutex),
	}
}

func (u *Users) Update(userID string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	// TODO check if locked
	// TODO check if new (unknown) user
	// TODO update user info
	// TODO return if user allowed
}
