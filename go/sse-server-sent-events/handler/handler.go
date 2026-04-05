package handler

import (
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"sse/room"
	"sse/static"
	"sse/user"
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
		roomID := q.Get("room")
		userID := q.Get("user")
		// TODO validate, set defaults
		wall, users, isNew := ch.Room(roomID)
		_ = isNew                   // TODO?
		nik, _ := users.Get(userID) // TODO in fact, just check if user exists
		if len(nik) == 0 {
			http.Error(w, "Forbidden", http.StatusForbidden) // TODO reset; TODO superfluous response.WriteHeader call from
		}
		leid, err := strconv.ParseInt(r.Header.Get("Last-Event-Id"), 10, 64)
		if err != nil {
			leid = 0
		}
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		messages := [][]byte(nil) // we have to create this var out of the loop, as leid
		for {                     // TODO check writing errors
			messages, leid = wall.Fetch(ctx, leid)
			if ctx.Err() != nil {
				return
			}
			// TODO user update? touch only
			w.Write([]byte("event: message\n"))                          // it MUST be `message` to make e.onmessage be fired
			w.Write([]byte("retry: 200\n"))                              // server side control for reconnecting delay
			w.Write([]byte("id: " + strconv.FormatInt(leid, 10) + "\n")) // it will be `Last-Event-Id: TOKEN` (on request)
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
		name := sanitize(req.Name)   // TODO validate
		color := sanitize(req.Color) // TODO validate
		roomID := req.Room           // TODO validate
		userID := req.User           // TODO validate
		wall, users := ch.RoomOrNil(roomID)
		if users == nil {
			return
		}
		allowed, updated := users.Touch(userID, ms, name, color)
		if allowed {
			u := (*user.Users)(nil)
			if updated {
				u = users
			}
			wall.Pub(buildResponse(&MessageDTO{
				Color:      color,
				Message:    sanitize(req.Message),
				Name:       name,
				TimeStamep: 0, // TODO
			}, u))
		} else {
			log.Printf("WARNING: User is not allowed! room=%s, user=%s", roomID, userID)
		}
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
	Users   []UserDTO   `json:"users,omitempty"`
	Locked  *bool       `json:"locked,omitempty"`
}

func buildResponse(message *MessageDTO, users *user.Users) []byte { // TODO do not use *user.Users, use DTOs only
	v := []UserDTO(nil)
	c := (*bool)(nil)
	if users != nil {
		c = ptr(users.Locked())
		u := users.List()
		for _, x := range u {
			v = append(v, UserDTO{
				Name:  x[0],
				Color: x[1],
			})
		}
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
		// TODO log
		return nil
	}
	dto := new(RequestDTO)
	err = json.Unmarshal(body, dto)
	if err != nil {
		// TODO log
		return nil
	}
	// TODO validate, set defaults
	return dto
}

func buildRobotMessage(s string) *MessageDTO {
	return &MessageDTO{
		Color:      "#990099",
		Message:    s,
		Name:       "#ROBOT",
		TimeStamep: 0, // TODO
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
		wall, users, _ := ch.Room(dto.Room)
		allowed, updated := users.Touch(dto.User, ms, dto.Name, dto.Color)
		if !allowed {
			return // TODO http response
		}
		if updated {
			wall.Pub(buildResponse(buildRobotMessage(dto.Name+" HERE!"), users))
		}
		body := buildResponse(nil, users)
		w.Write(body) // TODO user io.copy, check error
	}
}

func handlerLock(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := readBody(r.Body)
		wall, users := ch.RoomOrNil(req.Room)
		if wall == nil {
			return
		}
		name, _ := users.Get(req.User)
		if len(name) == 0 {
			return
		}
		if users.Lock(req.Lock) {
			wall.Pub(buildResponse(buildRobotMessage("Someone enters"), users))
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
		switch r.Method {
		case http.MethodGet:
			switch path {
			case "/fetch":
				fetchh.ServeHTTP(w, r)
				return
			case "/dump":
				dumph.ServeHTTP(w, r)
				return
			}
			handleStatic(fsh).ServeHTTP(w, r)
			return
		case http.MethodPost:
			switch path {
			case "/pub":
				pubh.ServeHTTP(w, r)
				return
			case "/enter":
				enterh.ServeHTTP(w, r)
				return
			case "/lock":
				lockh.ServeHTTP(w, r)
				return
			}
		default:
			http.Error(w, "not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func Handler(house *room.House) http.Handler {
	return handler(static.FS, house)
}

func handleStatic(fsh http.Handler) http.HandlerFunc { // TODO move to MW package?
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Add("Cache-Control", "no-cache")
		fsh.ServeHTTP(w, r)
	}
}

func ptr[T any](x T) *T { return &x }
