package main

import (
	"context"
	"net/http"
	"time"

	"github.com/michurin/minlog"

	"github.com/michurin/warehouse/go/chat/pkg/chat"
)

func setupTrivial(mux *http.ServeMux) {
	// TODO validations
	storage := chat.New(chat.InitialLastID())
	wrapper := NewWraper("trivial")
	mux.Handle("/api/publish", wrapper(&chat.PublishHandler{Storage: storage}))
	mux.Handle("/api/poll", wrapper(&chat.PollHandler{Storage: storage}))
}

func setupSmall(mux *http.ServeMux) {
	// TODO validateions
	// TODO multi room
	storage := chat.New(chat.InitialLastID())
	wrapper := NewWraper("small")
	mux.Handle("/api/small/publish", wrapper(&chat.PublishHandler{Storage: storage}))
	mux.Handle("/api/small/poll", wrapper(&chat.PollHandler{Storage: storage}))
}

func main() {
	ctx := context.Background()
	setupLogger()
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public_html")))
	setupTrivial(mux)
	setupSmall(mux)
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
