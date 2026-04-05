package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode"

	"sse/dto"
	"sse/room"
	"sse/static"
	"sse/user"
	"sse/wall"
)

const pollingTimeout = 28 * time.Second

type dtoIn struct {
	Color   string `json:"color"`
	Message string `json:"message"`
	Name    string `json:"name"`
	Room    string `json:"room"`
	User    string `json:"user"`
}

func sanitize(x string) string {
	return strings.Map(func(x rune) rune {
		if unicode.IsControl(x) { // clean up \n as well, useful in JSON sanitizing perspective
			return '\x20'
		}
		return x
	}, x)
}

func roomAndUserFromGet(v url.Values) (string, string) {
	roomID := v.Get("room")
	userID := v.Get("user")
	// TODO validate, set defaults
	return roomID, userID
}

type LockRequisetDTO struct {
	Room string `json:"room"`
	User string `json:"user"`
}

func roomAndUserFromPost(r io.Reader) (string, string) { // TODO legacy
	body, err := io.ReadAll(r)
	if err != nil {
		panic(err) // TODO
	}
	u := new(LockRequisetDTO)
	err = json.Unmarshal(body, u)
	if err != nil {
		panic(err) // TODO
	}
	// TODO validate, set defaults
	return u.Room, u.User
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

		roomID, userID := roomAndUserFromGet(r.URL.Query())
		wall, users, isNew := ch.Room(roomID)
		_ = isNew                                          // TODO?
		allowed, updated := users.Touch(userID, 0, "", "") // TODO in fact: add seeded user
		if !allowed {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
		if updated {
			pub(wall, buildMessage(0, "#ROBOT", "#f00", "Someone enter"), buildRoomStatus(users))
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
		b := new(bytes.Buffer)
		_, err := io.Copy(b, r.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		req := new(dtoIn)
		err = json.Unmarshal(b.Bytes(), req)
		if err != nil {
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
			status := (*dto.RoomStatus)(nil)
			if updated {
				status = buildRoomStatus(users)
			}
			pub(wall, buildMessage(0, name, color, sanitize(req.Message)), status)
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

type ResponseDTO struct {
	Message *string   `json:"message,omitempty"`
	Users   []UserDTO `json:"users,omitempty"`
	Locked  *bool     `json:"locked,omitempty"`
}

func buildResponse(message string, users *user.Users) []byte {
	dto := ResponseDTO{}
	if len(message) > 0 {
		dto.Message = &message
	}
	if users != nil {
		dto.Locked = ptr(users.Locked())
		u := users.List()
		v := make([]UserDTO, len(u))
		for i, x := range u {
			v[i] = UserDTO{
				Name:  x[0],
				Color: x[1],
			}
		}
		dto.Users = v
	}
	b, _ := json.Marshal(dto) // TODO err
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

// --------------------------

func handlerEnter(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dto := readBody(r.Body)
		if dto == nil {
			return // TODO http.Error
		}
		ms := time.Now().UnixMilli()
		wall, users, _ := ch.Room(dto.Room)
		_ = users
		allowed, updated := users.Touch(dto.User, ms, dto.Name, dto.Color)
		if !allowed {
			return // TODO http response
		}
		// TODO if updated
		_ = updated
		_ = wall
		body := buildResponse("", users)
		w.Write(body) // TODO user io.copy, check error
	}
}

func handlerLock(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID, userID := roomAndUserFromPost(r.Body)
		wall, users := ch.RoomOrNil(roomID)
		if wall == nil {
			return
		}
		name, _ := users.Get(userID)
		if len(name) == 0 {
			return
		}
		if users.Lock() {
			pub(wall, buildMessage(0, "#ROBOT", "#ff0000", "Room is locked by "+name), buildRoomStatus(users))
		}
	}
}

func handlerUnlock(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID, userID := roomAndUserFromPost(r.Body)
		wall, users := ch.RoomOrNil(roomID)
		if wall == nil {
			return
		}
		name, _ := users.Get(userID)
		if len(name) == 0 {
			return
		}
		if users.Unlock() {
			pub(wall, buildMessage(0, "#ROBOT", "#ff0000", "Room is UNLOCKED by "+name), buildRoomStatus(users))
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
	unlockh := handlerUnlock(house)
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
			case "/unlock":
				unlockh.ServeHTTP(w, r)
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

func buildRoomStatus(users *user.Users) *dto.RoomStatus { // TODO legacy
	u := users.List()
	v := make([]dto.User, len(u))
	for i, x := range u {
		v[i] = dto.User{
			Name:  x[0],
			Color: x[1],
		}
	}
	return &dto.RoomStatus{
		Locked: users.Locked(),
		Users:  v,
	}
}

func buildMessage(ms int64, name, color, message string) *dto.Message { // TODO LEGACY
	return &dto.Message{
		Color:      color,
		Message:    message,
		Name:       name,
		TimeStamep: ms,
	}
}

func pub(wall *wall.Wall, m *dto.Message, s *dto.RoomStatus) error { // TODO process this error on caller side? or just log this error?
	messageBytes, err := json.Marshal(dto.StreamMessage{Message: m, RoomStatus: s})
	if err != nil {
		return err
	}
	wall.Pub(messageBytes)
	return nil
}

func ptr[T any](x T) *T { return &x }
