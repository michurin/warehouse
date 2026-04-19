package handlerenter

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
	if dto == nil {
		slog.ErrorContext(ctx, "body reading")
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !xdto.CanonicalName(dto.User) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid user id %q", dto.User))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !xdto.CanonicalName(dto.Name) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid user name %q", dto.Name))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !xdto.CanonicalName(dto.Room) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid room name %q", dto.Room))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !xdto.ValidColor(dto.Color) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid color %q", dto.Color))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	wall, users := h.house.Room(dto.Room)
	ms := time.Now().UnixMilli()
	allowed, updated := users.Touch(dto.User, ms, dto.Name, dto.Color)
	if !allowed {
		w.Write(xdto.BuildResponse(xdto.BuildControlMessage("locked"), nil))
		return
	}
	if updated {
		wall.Pub(xdto.BuildResponse(xdto.BuildRobotMessage(ms, dto.Name+" HERE!"), users)) // TODO(2) template and localization
	}
	w.Write(xdto.BuildResponse(nil, users)) // TODO(2) user io.copy, check error
}
