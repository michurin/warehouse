package main

import (
	"context"
	"net/http"
	"time"

	"github.com/michurin/minlog"

	"github.com/michurin/warehouse/go/chat/pkg/chat"
)

func setupTrivial(ctx context.Context, log chat.Logger, mux *http.ServeMux) {
	rooms := chat.New(log)
	go chat.RoomsCleaner(minlog.Label(ctx, "trivial:tick"), rooms)
	mux.Handle("/api/publish", NewPublishingHandler(rooms, trivialValidator, log, "trivial"))
	mux.Handle("/api/poll", NewPollingHandler(rooms, log, "trivial"))
	mux.Handle("/mon/trivial", NewMonitoringHandler(rooms))
}

func setupSimple(ctx context.Context, log chat.Logger, mux *http.ServeMux) {
	rooms := chat.New(log)
	go chat.RoomsCleaner(minlog.Label(ctx, "simple:tick"), rooms)
	mux.Handle("/api/simple/publish", NewPublishingHandler(rooms, simpleValidator, log, "simple"))
	mux.Handle("/api/simple/poll", NewPollingHandler(rooms, log, "simple"))
	mux.Handle("/mon/simple", NewMonitoringHandler(rooms))
}

func main() {
	ctx := context.Background()
	log := setupLogger()
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public_html")))
	setupTrivial(ctx, log, mux)
	setupSimple(ctx, log, mux)
	log.Log(ctx, "Listening...")
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    300 * time.Second, // 300 is most browsers default
		WriteTimeout:   300 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	err := s.ListenAndServe()
	if err != nil {
		log.Log(ctx, err)
	}
}
