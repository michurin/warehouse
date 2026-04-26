package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/michurin/minchat/internal/handlerenter"
	"github.com/michurin/minchat/internal/handlerfetch"
	"github.com/michurin/minchat/internal/handlerlock"
	"github.com/michurin/minchat/internal/handlerpub"
	"github.com/michurin/minchat/internal/handlerstat"
	"github.com/michurin/minchat/internal/handlerstatic"
	"github.com/michurin/minchat/internal/xdto"
	"github.com/michurin/minchat/internal/xhouse"
)

func handler(house *xhouse.House, pollingTimeout time.Duration) http.HandlerFunc {
	fsh := handlerstatic.New()
	enterh := handlerenter.New(house)
	pubh := handlerpub.New(house)
	fetchh := handlerfetch.New(house, pollingTimeout)
	lockh := handlerlock.New(house)
	dumph := handlerstat.New(house)
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
			room, _ := strings.CutPrefix(path, "/")
			if xdto.CanonicalName(room) {
				w.Header().Set("Cache-Control", "no-cache")
				r.URL.Path = "/chat.html"
				fsh.ServeHTTP(w, r)
				return
			} else {
				http.Redirect(w, r, "/main", http.StatusPermanentRedirect)
				return
			}
		}
		http.Error(w, "405 not allowed", http.StatusMethodNotAllowed)
	}
}

func Handler(house *xhouse.House, pollingTimeout time.Duration) http.Handler {
	return http.MaxBytesHandler(handler(house, pollingTimeout), 4096)
}
