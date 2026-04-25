package handlerfetch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/michurin/minchat/internal/xdto"
	"github.com/michurin/minchat/internal/xhouse"
	"github.com/michurin/minchat/internal/xuser"
	"github.com/michurin/minchat/internal/xwall"
)

type Handler struct {
	house   *xhouse.House
	timeout time.Duration
}

func New(house *xhouse.House, timeout time.Duration) *Handler {
	return &Handler{house: house, timeout: timeout}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, h.timeout)
	defer cancel()
	hdr := w.Header()
	hdr.Add("X-Accel-Buffering", "no")
	hdr.Add("Content-Type", "text/event-stream")
	hdr.Add("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)

	q := r.URL.Query()
	roomID := q.Get("room")
	userID := q.Get("user")
	if !xdto.CanonicalName(userID) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid user id %q", userID))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !xdto.CanonicalName(roomID) {
		slog.ErrorContext(ctx, fmt.Sprintf("invalid room name %q", roomID))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	leid, err := strconv.ParseInt(r.Header.Get("Last-Event-Id"), 10, 64)
	if err != nil {
		leid = 0
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	messages := [][]byte(nil) // we have to create this var out of the loop, as leid
	wl := (*xwall.Wall)(nil)
	us := (*xuser.Users)(nil)
	for {
		wl, us = h.house.RoomOrNil(roomID)
		if wl == nil {
			slog.Error("Kick. No room", slog.String("user", userID), slog.String("room", roomID))
			writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil, false)}) // reason: no room
			return
		}
		name, _ := us.Get(userID) // check user before feetching // TODO(2) in fact, just check if user exists
		if len(name) == 0 {
			writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil, false)}) // reason: no user
			return
		}
		messages, leid = wl.Fetch(ctx, leid) // it will take a time. So we need to check user after that again
		if ctx.Err() != nil {
			slog.ErrorContext(ctx, ctx.Err().Error())
			return
		}
		name, _ = us.Get(userID) // check user before sending // TODO(2) in fact, just check if user exists
		if len(name) == 0 {
			writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil, false)}) // reason: no user
			return
		}
		writeStreamMessage(w, leid, messages)
	}
}

func writeStreamMessage(w io.Writer, leid int64, messages [][]byte) {
	// TODO(2) check writing errors
	w.Write([]byte("event: message\n")) // message to e.onmessage
	w.Write([]byte("retry: 200\n"))     // server side control for reconnecting delay
	w.Write([]byte("id: "))
	w.Write([]byte(strconv.FormatInt(leid, 10))) // it will be `Last-Event-Id: TOKEN` (on request)
	w.Write([]byte{10})
	for _, m := range messages {
		w.Write([]byte("data: "))
		w.Write(m) // we are storing single line messages only
		w.Write([]byte{10})
	}
	w.Write([]byte{10})
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	} else {
		panic("http.Flusher is not supported")
	}
}
