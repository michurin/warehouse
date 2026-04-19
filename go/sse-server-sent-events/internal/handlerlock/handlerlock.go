package handlerlock

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/michurin/minchat/internal/xdto"
	"github.com/michurin/minchat/internal/xhouse"
)

type Handler struct {
	house *xhouse.House
}

func New(house *xhouse.House) *Handler {
	return &Handler{house: house}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	dto := xdto.ReadBody(r.Body)
	if !xdto.CanonicalName(dto.User) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid user id %q", dto.User))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !xdto.CanonicalName(dto.Room) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid room name %q", dto.Room))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	wall, users := h.house.RoomOrNil(dto.Room)
	if wall == nil {
		slog.ErrorContext(ctx, "lock room: "+dto.Room+" (not found)")
		return
	}
	name, _ := users.Get(dto.User)
	if len(name) == 0 {
		slog.ErrorContext(ctx, "cannot lock room: "+dto.Room+" by user: "+dto.User)
		return
	}
	if users.Lock(dto.Lock) {
		ms := time.Now().UnixMilli()
		wall.Pub(xdto.BuildResponse(xdto.BuildRobotMessage(ms, name+" touched LOCK"), users))
	}
}
