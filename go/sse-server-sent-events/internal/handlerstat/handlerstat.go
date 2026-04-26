package handlerstat

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"

	"github.com/michurin/minchat/internal/xhouse"
)

type Handler struct {
	house *xhouse.House
}

func New(house *xhouse.House) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
	mux.Handle("GET /stat", &Handler{house: house})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<ul>
<li><a href="/stat">/stat</a></li>
<li><a href="/debug/pprof/">/debug/pprof/</a></li>
</ul>`))
	})
	return mux
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := map[string]any{}
	for _, room := range h.house.List() {
		wall, users := h.house.RoomOrNil(room)
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
