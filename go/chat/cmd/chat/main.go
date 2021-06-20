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
	chat.RoomCleaner(minlog.Label(ctx, "tick:trivial"), rooms)
	wrapper := NewWraper("trivial")
	mux.Handle("/api/publish", wrapper(chat.NewPublishingHandler(rooms, trivialValidator)))
	mux.Handle("/api/poll", wrapper(chat.NewPollingHandler(rooms)))
}

func setupSimple(ctx context.Context, mux *http.ServeMux) {
	rooms := chat.New()
	chat.RoomCleaner(minlog.Label(ctx, "tick:simple"), rooms)
	wrapper := NewWraper("small")
	mux.Handle("/api/small/publish", wrapper(chat.NewPublishingHandler(rooms, simpleValidator)))
	mux.Handle("/api/small/poll", wrapper(chat.NewPollingHandler(rooms)))
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
