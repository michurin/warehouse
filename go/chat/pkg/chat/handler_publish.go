package chat

import (
	"encoding/json"
	"net/http"

	"github.com/michurin/minlog"
)

type publishRequest struct {
	RoomID  string          `json:"room"`
	Message json.RawMessage `json:"message"`
}

type PublishHandler struct {
	Rooms *Rooms // TODO private
}

func (h *PublishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := new(publishRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		minlog.Log(r.Context(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := minlog.Label(r.Context(), "room:"+req.RoomID)
	minlog.Log(ctx, "Publish:", []byte(req.Message))
	// TODO validate message
	// TODO validate room id
	h.Rooms.Pub(ctx, req.RoomID, req.Message)
	hdr := w.Header()
	hdr.Set("content-type", "application/json; charset=UTF-8")
	w.Write([]byte(`{}`)) // JSON
}
