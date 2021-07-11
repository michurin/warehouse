package main

import (
	"context"
	"net/http"
	"time"

	"github.com/michurin/minlog"

	"github.com/michurin/warehouse/go/chat/pkg/chat"
)

func setupTrivial(ctx context.Context, mux *http.ServeMux) {
	rooms := chat.New()
	go chat.RoomsCleaner(minlog.Label(ctx, "tick:trivial"), rooms)
	mux.Handle("/api/publish", NewPublishingHandler(rooms, trivialValidator, "trivial"))
	mux.Handle("/api/poll", NewPollingHandler(rooms, "trivial"))
	mux.Handle("/mon/trivial", NewMonitoringHandler(rooms)) // TODO HTML: links to this
}

func setupSimple(ctx context.Context, mux *http.ServeMux) {
	rooms := chat.New()
	go chat.RoomsCleaner(minlog.Label(ctx, "tick:simple"), rooms)
	mux.Handle("/api/small/publish", NewPublishingHandler(rooms, simpleValidator, "small"))
	mux.Handle("/api/small/poll", NewPollingHandler(rooms, "small"))
	mux.Handle("/mon/small", NewMonitoringHandler(rooms))
}

func main() {
	ctx := context.Background()
	setupLogger()
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public_html")))
	setupTrivial(ctx, mux)
	setupSimple(ctx, mux)
	minlog.Log(ctx, "Listening...")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    300 * time.Second, // 300 is most browsers default
		WriteTimeout:   300 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	err := s.ListenAndServe()
	if err != nil {
		minlog.Log(ctx, err)
	}
}
