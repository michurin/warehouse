package handler

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"sse/room"
	"sse/static"
	"sse/user"
	"sse/wall"
)

const pollingTimeout = 28 * time.Second

func sanitize(x string) string {
	return strings.Map(func(x rune) rune {
		if unicode.IsControl(x) { // clean up \n as well, useful in JSON sanitizing perspective
			return '\x20'
		}
		return x
	}, x)
}

func strictSanitaze(x string) string {
	return strings.Map(func(x rune) rune {
		if x == '_' || x == '-' || ('A' <= x && x <= 'Z') || ('a' <= x && x <= 'z') || ('0' <= x && x <= '9') {
			return x
		}
		return -1
	}, x)
}

func colorSanitaze(x string) string {
	return strings.Map(func(x rune) rune {
		if x == '#' || ('A' <= x && x <= 'F') || ('a' <= x && x <= 'f') || ('0' <= x && x <= '9') {
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
			writeStreamMessage(w, 0, [][]byte{buildResponse(buildControlMessage(), nil)}) // reason: invalid user or room
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
				writeStreamMessage(w, 0, [][]byte{buildResponse(buildControlMessage(), nil)}) // reason: no room
				return
			}
			name, _ := us.Get(userID) // check user before feetching // TODO in fact, just check if user exists
			if len(name) == 0 {
				writeStreamMessage(w, 0, [][]byte{buildResponse(buildControlMessage(), nil)}) // reason: no user
				return
			}
			messages, leid = wl.Fetch(ctx, leid)
			if ctx.Err() != nil {
				return
			}
			name, _ = us.Get(userID) // check user before sending // TODO in fact, just check if user exists
			if len(name) == 0 {
				writeStreamMessage(w, 0, [][]byte{buildResponse(buildControlMessage(), nil)}) // reason: no user
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
	}
}

func handlerPub(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := readBody(r.Body)
		if req == nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		ms := time.Now().UnixMilli()
		name := strictSanitaze(req.Name)
		color := colorSanitaze(req.Color)
		roomID := strictSanitaze(req.Room)
		userID := strictSanitaze(req.User)
		// TODO check empty
		wall, users := ch.RoomOrNil(roomID)
		if users == nil {
			return
		}
		allowed, updated := users.Touch(userID, ms, name, color)
		if !allowed {
			log.Printf("WARNING: User is not allowed! room=%s, user=%s", roomID, userID)
			http.Error(w, "Not allowed", http.StatusOK) // TODO error
			return
		}
		if updated {
			wall.Pub(buildResponse(buildRobotMessage(ms, "User updated "+name), users))
		}
		wall.Pub(buildResponse(&MessageDTO{
			Color:      color,
			Message:    sanitize(req.Message),
			Name:       name,
			TimeStamep: ms,
		}, nil))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// ----- Transport DTOs -----

type RequestDTO struct {
	Room    string `json:"room"`
	User    string `json:"user"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Lock    bool   `json:"lock"`    // /lock only
	Message string `json:"message"` // /pub only
}

type UserDTO struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type MessageDTO struct {
	Color      string `json:"color"`
	Message    string `json:"message"`
	Name       string `json:"name"`
	TimeStamep int64  `json:"ts"`
}

type ResponseDTO struct {
	Message *MessageDTO `json:"message,omitempty"`
	Users   *[]UserDTO  `json:"users,omitempty"`
	Locked  *bool       `json:"locked,omitempty"`
}

func buildResponse(message *MessageDTO, users *user.Users) []byte { // TODO do not use *user.Users, use DTOs only
	v := (*[]UserDTO)(nil)
	c := (*bool)(nil)
	if users != nil {
		w := []UserDTO{} // force empty array, not nil
		c = ptr(users.Locked())
		u := users.List()
		for _, x := range u {
			w = append(w, UserDTO{
				Name:  x[0],
				Color: x[1],
			})
		}
		v = &w
	}
	b, _ := json.Marshal(ResponseDTO{ // TODO err
		Message: message,
		Users:   v,
		Locked:  c,
	})
	return b
}

func readBody(r io.Reader) *RequestDTO {
	body, err := io.ReadAll(r)
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	dto := new(RequestDTO)
	err = json.Unmarshal(body, dto)
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	// TODO validate, set defaults
	return dto
}

func buildRobotMessage(ms int64, s string) *MessageDTO {
	return &MessageDTO{
		Color:      "#990099",
		Message:    s,
		Name:       "#ROBOT",
		TimeStamep: ms,
	}
}

func buildControlMessage() *MessageDTO {
	ms := time.Now().UnixMilli() // TODO
	return &MessageDTO{
		Color:      "#333333",
		Message:    "",
		Name:       "#CONTROL",
		TimeStamep: ms,
	}
}

// --------------------------

func handlerEnter(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := readBody(r.Body)
		if dto == nil {
			return // TODO http.Error
		}
		ms := time.Now().UnixMilli()
		wall, users := ch.Room(dto.Room)
		allowed, updated := users.Touch(dto.User, ms, dto.Name, dto.Color)
		if !allowed {
			return // TODO http response
		}
		if updated {
			ms := time.Now().UnixMilli()
			wall.Pub(buildResponse(buildRobotMessage(ms, dto.Name+" HERE!"), users))
		}
		body := buildResponse(nil, users)
		log.Print(string(body))
		w.Write(body) // TODO user io.copy, check error
	}
}

func handlerLock(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := readBody(r.Body)
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
			wall.Pub(buildResponse(buildRobotMessage(ms, name+" touched LOCK"), users))
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

func handler(staticFS fs.FS, house *room.House) http.HandlerFunc {
	fsh := http.FileServerFS(staticFS)
	fetchh := handlerFetch(house)
	pubh := handlerPub(house)
	lockh := handlerLock(house)
	enterh := handlerEnter(house)
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
	return http.MaxBytesHandler(handler(static.FS, house), 4096)
}

func ptr[T any](x T) *T { return &x }

// ---------- REVISION ---------- TODO move to package?

func RevisionLoop(ch *room.House) {
	for {
		ms := time.Now().Add(-10 * time.Second).UnixMilli()
		walls, users := ch.Audit(ms)
		for i, w := range walls {
			log.Print("Run: notify")
			w.Pub(buildResponse(buildRobotMessage(ms, "Someone got out"), users[i]))
		}
		time.Sleep(2 * time.Second)
	}
}
