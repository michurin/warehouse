package xdto

import (
	"encoding/json"
	"io"
	"log/slog"
)

type RequestDTO struct {
	Room    string `json:"room"`
	User    string `json:"user"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Lock    bool   `json:"lock"`    // /lock only
	Message string `json:"message"` // /pub only
}

func ReadBody(r io.Reader) *RequestDTO { // TODO(2) return error
	body, err := io.ReadAll(r)
	if err != nil {
		slog.Error("Read body: " + err.Error())
		return nil
	}
	dto := new(RequestDTO)
	err = json.Unmarshal(body, dto)
	if err != nil {
		slog.Error("Unmarshal body: " + err.Error())
		return nil
	}
	return dto
}

func CanonicalName(n string) bool {
	if len(n) == 0 {
		return false
	}
	for i, c := range n {
		if i > 32 || !(('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

func ValidColor(c string) bool {
	if len(c) != 7 {
		return false
	}
	if c[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		b := c[i]
		if !(('a' <= b && b <= 'f') || ('A' <= b && b <= 'F') || ('0' <= b && b <= '9')) {
			return false
		}
	}
	return true
}
