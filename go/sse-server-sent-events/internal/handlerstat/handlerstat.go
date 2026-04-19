package handlerstat

import (
	"encoding/json"
	"net/http"

	"github.com/michurin/minchat/internal/xhouse"
)

type Handler struct {
	house *xhouse.House
}

func New(house *xhouse.House) *Handler {
	return &Handler{house: house}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := map[string]any{}
	for _, room := range h.house.List() {
		wall, users := h.house.RoomOrNil(room)
		if wall == nil {
			continue
		}
		res[room] = map[string]any{
			"users": users.List(),
			"lock":  users.Locked(),
		}
	}
	j := json.NewEncoder(w)
	j.SetIndent("", "  ")
	j.Encode(res) // TODO err
}
