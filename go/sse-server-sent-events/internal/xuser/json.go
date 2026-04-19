package xuser

import (
	"encoding/json"
	"strconv"
)

type marshaler struct {
	users *Users
}

type dtoUser struct {
	Name      string `json:"nik"`
	Color     string `json:"color"`
	LastCheck string `json:"last_check"`
}

func JSON(u *Users) interface {
	json.Marshaler
	json.Unmarshaler
} {
	return &marshaler{users: u}
}

func (r *marshaler) MarshalJSON() ([]byte, error) {
	r.users.mu.RLock()
	defer r.users.mu.RUnlock()
	u := map[string]dtoUser{}
	for k, v := range r.users.users {
		u[k] = dtoUser{
			Name:      v.name,
			Color:     v.color,
			LastCheck: strconv.FormatInt(v.lastCheck, 10),
		}
	}
	dto := struct {
		Locked bool               `json:"locked"`
		Users  map[string]dtoUser `json:"users"`
	}{
		Locked: r.users.locked,
		Users:  u,
	}
	return json.Marshal(dto)
}

func (r *marshaler) UnmarshalJSON(data []byte) error {
	r.users.mu.Lock()
	defer r.users.mu.Unlock()
	// TODO load
	return nil
}
