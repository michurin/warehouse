package handlerpub

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"sse/internal/xdto"
	"sse/room"
)

type Handler struct {
	house *room.House
}

func New(house *room.House) *Handler {
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
	wall, users := h.house.RoomOrNil(dto.Room)
	if users == nil {
		w.Write(xdto.BuildResponse(xdto.BuildControlMessage("noroom"), nil))
		return
	}
	ms := time.Now().UnixMilli()
	allowed, updated := users.Touch(dto.User, ms, dto.Name, dto.Color)
	if !allowed {
		slog.ErrorContext(ctx, fmt.Sprintf("WARNING: User is not allowed! room=%q, user=%q, name=%q", dto.Room, dto.User, dto.Name))
		w.Write(xdto.BuildResponse(xdto.BuildControlMessage("locked"), nil))
		return
	}
	if updated {
		wall.Pub(xdto.BuildResponse(xdto.BuildRobotMessage(ms, "User updated "+dto.Name), users))
	}
	text := xdto.SanitizeMessage(dto.Message)
	if len(text) > 0 {
		wall.Pub(xdto.BuildResponse(&xdto.MessageDTO{
			Color:      dto.Color,
			Message:    xdto.SanitizeMessage(dto.Message),
			Name:       dto.Name,
			TimeStamep: ms,
		}, nil))
	}
	w.WriteHeader(http.StatusOK)
}
