package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type pollingRequest struct {
	ID int `json:"id"`
}

type PollHandler struct {
	Storage *Storage
	// TODO add validator
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
	mm, lastID := h.Storage.Get(ctx, req.ID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": mm,
		"lastID":   lastID,
	})
}
