package xdto

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/michurin/minchat/internal/xuser"
)

type UserDTO struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type MessageDTO struct {
	Color      string `json:"color"`
	Message    string `json:"message"`
	Name       string `json:"name"`
	TimeStamep int64  `json:"ts"`
}

type ResponseDTO struct {
	Message *MessageDTO `json:"message,omitempty"`
	Users   *[]UserDTO  `json:"users,omitempty"`
	Locked  *bool       `json:"locked,omitempty"`
}

func BuildResponse(message *MessageDTO, users *xuser.Users) []byte { // TODO do not use *user.Users, use DTOs only
	v := (*[]UserDTO)(nil)
	c := (*bool)(nil)
	if users != nil {
		w := []UserDTO{} // force empty array, not nil
		c = new(users.Locked())
		u := users.List()
		for _, x := range u {
			w = append(w, UserDTO{
				Name:  x[0],
				Color: x[1],
			})
		}
		v = &w
	}
	b, _ := json.Marshal(ResponseDTO{ // TODO err
		Message: message,
		Users:   v,
		Locked:  c,
	})
	return b
}

// BuildRobotMessage
// Robot talks to all. That message is for publishing on the wall.
func BuildRobotMessage(ms int64, m string) *MessageDTO {
	return &MessageDTO{
		Color:      "#990099",
		Message:    m,
		Name:       "#ROBOT",
		TimeStamep: ms,
	}
}

// BuildControlMessage
// Control message is a message for one user only. It can be the
// last message on stream (without publishing for everyone on the wall),
// or the synchronous response.
func BuildControlMessage(m string) *MessageDTO {
	return &MessageDTO{
		Color:      "#333333",
		Message:    m,
		Name:       "#CONTROL",
		TimeStamep: 0, // as it doesn't appear on the wall it doesn't have time stamp
	}
}

func SanitizeMessage(x string) string {
	return strings.Map(func(x rune) rune {
		if unicode.IsControl(x) { // clean up \n as well, useful in JSON sanitizing perspective
			return '\x20'
		}
		return x
	}, x)
}
