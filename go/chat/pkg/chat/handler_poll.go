package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/michurin/minlog"
)

type pollingRequest struct {
	RoomID string `json:"room"`
	ID     int64  `json:"id"`
}

type pollingResponse struct {
	Messages []json.RawMessage `json:"messages"`
	LastID   int64             `json:"lastID"`
}

type PollingHandler struct {
	rooms *Rooms
}

func NewPollingHandler(r *Rooms) *PollingHandler {
	return &PollingHandler{
		rooms: r,
	}
}

func (h *PollingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second) // TODO make it tunable
	defer cancel()
	req := new(pollingRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		minlog.Log(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rid := req.RoomID
	id := req.ID
	ctx = minlog.Label(ctx, "room:"+rid)
	if err = validateRoomID(rid); err != nil {
		minlog.Log(ctx, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = validateID(id); err != nil {
		minlog.Log(ctx, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mm, lastID := h.rooms.Fetch(ctx, rid, id)
	hdr := w.Header()
	hdr.Set("content-type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(pollingResponse{
		Messages: mm,
		LastID:   lastID,
	})
	if err != nil {
		minlog.Log(ctx, err)
	}
}
