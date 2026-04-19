package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/michurin/minchat/internal/xlog"
)

type ResponseWriterFlusher interface {
	http.Flusher // in fact, Flasher is only necessary for SSE handler
	http.ResponseWriter
}

type wr struct {
	status   int
	location string
	next     ResponseWriterFlusher
}

var _ ResponseWriterFlusher = (*wr)(nil)

func (w *wr) Header() http.Header {
	return w.next.Header()
}

func (w *wr) Write(b []byte) (int, error) {
	return w.next.Write(b)
}

func (w *wr) Flush() {
	w.next.Flush()
}

func (w *wr) WriteHeader(statusCode int) {
	w.status = statusCode
	w.location = w.next.Header().Get("Location")
	w.next.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = xlog.WithRequestID(ctx)
		ctx = xlog.WithAddr(ctx, r.RemoteAddr)
		ctx = xlog.WithMethod(ctx, r.Method)
		ctx = xlog.WithURL(ctx, r.URL.String())
		wo := &wr{status: http.StatusOK, next: w.(ResponseWriterFlusher)}
		start := time.Now()
		defer func() {
			ctx = xlog.WithStatus(ctx, wo.status)
			if len(wo.location) > 0 {
				ctx = xlog.WithLocation(ctx, wo.location)
			}
			slog.InfoContext(ctx, fmt.Sprintf("%v", time.Since(start)))
		}()
		next.ServeHTTP(wo, r.WithContext(ctx))
	})
}
