package handlerstatic

import (
	"net/http"

	"sse/static"
)

type Handler struct {
	next http.Handler
}

func New() *Handler {
	return &Handler{next: http.FileServerFS(static.FS)}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.next.ServeHTTP(w, r)
}
