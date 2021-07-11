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
	wrapper := NewWraper("trivial")
	mux.Handle("/api/publish", wrapper(chat.NewPublishingHandler(rooms, chat.ValidatorFunc(trivialValidator))))
	mux.Handle("/api/poll", wrapper(chat.NewPollingHandler(rooms)))
	mux.Handle("/mon/trivial", NewMonHandler(rooms))
}

func setupSimple(ctx context.Context, mux *http.ServeMux) {
	rooms := chat.New()
	go chat.RoomsCleaner(minlog.Label(ctx, "tick:simple"), rooms)
	wrapper := NewWraper("small")
	mux.Handle("/api/small/publish", wrapper(chat.NewPublishingHandler(rooms, chat.ValidatorFunc(simpleValidator))))
	mux.Handle("/api/small/poll", wrapper(chat.NewPollingHandler(rooms)))
	mux.Handle("/mon/small", NewMonHandler(rooms))
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
