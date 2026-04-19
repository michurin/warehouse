package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"sse/internal/handlerenter"
	"sse/internal/handlerfetch"
	"sse/internal/handlerlock"
	"sse/internal/handlerpub"
	"sse/internal/handlerstatic"
	"sse/internal/xdto"
	"sse/room"
)

const pollingTimeout = 28 * time.Second

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
	fetchh := handlerfetch.New(house, pollingTimeout)
	lockh := handlerlock.New(house)
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
