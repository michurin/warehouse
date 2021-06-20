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

type PublishingHandler struct {
	rooms     *Rooms
	validator func(json.RawMessage) error // TODO it would be nice to have interface here
}

func NewPublishingHandler(r *Rooms, v func(json.RawMessage) error) *PublishingHandler {
	return &PublishingHandler{
		rooms:     r,
		validator: v,
	}
}

func (h *PublishingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := new(publishRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		minlog.Log(r.Context(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx := minlog.Label(r.Context(), "room:"+req.RoomID)
	minlog.Log(ctx, "Publish:", []byte(req.Message))
	rid := req.RoomID
	msg := req.Message
	if err = validateRoomID(rid); err != nil {
		minlog.Log(ctx, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = h.validator(msg); err != nil {
		minlog.Log(ctx, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.rooms.Pub(ctx, rid, msg)
	hdr := w.Header()
	hdr.Set("content-type", "application/json; charset=UTF-8")
	_, err = w.Write([]byte(`{}`)) // JSON
	if err != nil {
		minlog.Log(ctx, err)
	}
}
