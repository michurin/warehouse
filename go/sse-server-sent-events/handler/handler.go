package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sse/internal/handlerenter"
	"sse/internal/handlerpub"
	"sse/internal/handlerstatic"
	"sse/internal/xdto"
	"sse/room"
	"sse/user"
	"sse/wall"
)

const pollingTimeout = 28 * time.Second

func strictSanitaze(x string) string {
	return strings.Map(func(x rune) rune {
		if x == '_' || x == '-' || ('A' <= x && x <= 'Z') || ('a' <= x && x <= 'z') || ('0' <= x && x <= '9') {
			return x
		}
		return -1
	}, x)
}

func handlerFetch(ch *room.House) http.HandlerFunc {
	// TODO process errors, TODO use Copy
	return func(w http.ResponseWriter, r *http.Request) {
		// io.Copy(io.Discard, r.Body) // just drop body. We do not need to close it. Oh. It works without ctx
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, pollingTimeout)
		defer cancel()
		h := w.Header()
		h.Add("X-Accel-Buffering", "no")
		h.Add("Content-Type", "text/event-stream")
		h.Add("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)

		q := r.URL.Query()
		roomID := strictSanitaze(q.Get("room"))
		userID := strictSanitaze(q.Get("user"))
		if len(roomID) == 0 {
			roomID = "main"
		}
		if len(roomID) > 50 || len(userID) == 0 || len(userID) > 30 {
			log.Printf("ERROR roomID/userID=%q/%q", roomID, userID)
			writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil)}) // reason: invalid user or room
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
		wl := (*wall.Wall)(nil)
		us := (*user.Users)(nil)
		for {
			wl, us = ch.RoomOrNil(roomID)
			if wl == nil {
				slog.Error("Kick. No room", slog.String("user", userID), slog.String("room", roomID))
				writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil)}) // reason: no room
				return
			}
			name, _ := us.Get(userID) // check user before feetching // TODO in fact, just check if user exists
			if len(name) == 0 {
				writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil)}) // reason: no user
				return
			}
			messages, leid = wl.Fetch(ctx, leid)
			if ctx.Err() != nil {
				slog.ErrorContext(ctx, ctx.Err().Error())
				return
			}
			name, _ = us.Get(userID) // check user before sending // TODO in fact, just check if user exists
			if len(name) == 0 {
				writeStreamMessage(w, 0, [][]byte{xdto.BuildResponse(xdto.BuildControlMessage(""), nil)}) // reason: no user
				return
			}
			writeStreamMessage(w, leid, messages)
		}
	}
}

func writeStreamMessage(w io.Writer, leid int64, messages [][]byte) {
	// TODO check writing errors
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

func handlerLock(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := xdto.ReadBody(r.Body)
		wall, users := ch.RoomOrNil(req.Room)
		if wall == nil {
			log.Print("lock room: " + req.Room + " (not found)")
			return
		}
		name, _ := users.Get(req.User)
		if len(name) == 0 {
			log.Print("cannot lock room: " + req.Room + " by user: " + req.User)
			return
		}
		if users.Lock(req.Lock) {
			ms := time.Now().UnixMilli()
			wall.Pub(xdto.BuildResponse(xdto.BuildRobotMessage(ms, name+" touched LOCK"), users))
		}
	}
}

func handlerDump(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := map[string]any{}
		for _, room := range ch.List() {
			wall, users := ch.RoomOrNil(room)
			if wall == nil {
				continue
			}
			res[room] = map[string]any{
				"users": users.List(),
				"lock":  users.Locked(),
			}
		}
		j := json.NewEncoder(w)
		j.SetIndent("", "  ")
		j.Encode(res) // TODO err
	}
}

func handler(house *room.House) http.HandlerFunc {
	fsh := handlerstatic.New()
	enterh := handlerenter.New(house)
	pubh := handlerpub.New(house)
	fetchh := handlerFetch(house)
	lockh := handlerLock(house)
	dumph := handlerDump(house)
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		if binPath, ok := strings.CutPrefix(path, "/bin/"); ok {
			switch r.Method {
			case http.MethodGet:
				switch binPath {
				case "fetch":
					fetchh.ServeHTTP(w, r)
					return
				case "dump":
					dumph.ServeHTTP(w, r)
					return
				}
			case http.MethodPost:
				switch binPath {
				case "pub":
					pubh.ServeHTTP(w, r)
					return
				case "enter":
					enterh.ServeHTTP(w, r)
					return
				case "lock":
					lockh.ServeHTTP(w, r)
					return
				}
			}
		}
		if r.Method == http.MethodGet {
			if path == "/" || path == "/favicon.ico" {
				w.Header().Set("Cache-Control", "no-cache")
				fsh.ServeHTTP(w, r)
				return
			}
			if docPath, ok := strings.CutPrefix(path, "/s/"); ok {
				w.Header().Set("Cache-Control", "no-cache")
				r.URL.Path = "/" + docPath
				fsh.ServeHTTP(w, r)
				return
			}
			tail, _ := strings.CutPrefix(path, "/")
			key := strings.Map(func(x rune) rune {
				if ('a' <= x && x <= 'z') || ('A' <= x && x <= 'Z') || ('0' <= x && x <= '9') || x == '_' || x == '-' {
					return x
				}
				return -1
			}, tail)
			if key == tail {
				w.Header().Set("Cache-Control", "no-cache")
				r.URL.Path = "/chat.html"
				fsh.ServeHTTP(w, r)
				return
			} else {
				http.Redirect(w, r, "/"+key, http.StatusPermanentRedirect)
				return
			}
		}
		http.Error(w, "405 not allowed", http.StatusMethodNotAllowed)
	}
}

func Handler(house *room.House) http.Handler {
	return http.MaxBytesHandler(handler(house), 4096)
}

// ---------- REVISION ---------- TODO move to package?

func RevisionLoop(ch *room.House) {
	for {
		ms := time.Now().Add(-10 * time.Second).UnixMilli()
		walls, users := ch.Audit(ms)
		for i, w := range walls {
			log.Print("Run: notify")
			w.Pub(xdto.BuildResponse(xdto.BuildRobotMessage(ms, "Someone got out"), users[i]))
		}
		time.Sleep(2 * time.Second)
	}
}
