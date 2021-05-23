package chat

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type PollHandler struct {
	Storage *Storage
}

func (h *PollHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second) // TODO make it tunable
	defer cancel()
	err := r.ParseForm()
	id, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mm, lastID := h.Storage.Get(ctx, id)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": mm,
		"lastID":   lastID,
	})
}
