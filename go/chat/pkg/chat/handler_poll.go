package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type pollingRequest struct {
	RoomID string `json:"room"`
	ID     int64  `json:"id"`
}

type PollHandler struct {
	Rooms *Rooms
}

func (h *PollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second) // TODO make it tunable
	defer cancel()
	req := new(pollingRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO validate room id
	// TODO validate id
	mm, lastID := h.Rooms.Fetch(ctx, req.RoomID, req.ID)
	hdr := w.Header()
	hdr.Set("content-type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": mm,
		"lastID":   lastID,
	})
}
