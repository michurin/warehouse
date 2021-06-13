package chat

import (
	"encoding/json"
	"net/http"

	"github.com/michurin/minlog"
)

type publishRequest struct {
	Message json.RawMessage `json:"message"`
}

type PublishHandler struct {
	Storage *Storage
}

func (h *PublishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := new(publishRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		minlog.Log(r.Context(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	minlog.Log(r.Context(), "Publish:", []byte(req.Message))
	// TODO validate message
	h.Storage.Put(req.Message)
	hdr := w.Header()
	hdr.Set("content-type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`)) // JSON
}
