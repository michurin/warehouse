package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/michurin/minlog"

	"github.com/michurin/warehouse/go/chat/pkg/chat"
)

func main() {
	ctx := context.Background()
	minlog.SetDefaultLogger(minlog.New(
		minlog.WithLabelPlaceholder("-"),
		minlog.WithLineFormatter(func(tm, level, label, caller, msg string) string {
			c := "\033[32;1m"
			if level != minlog.DefaultInfoLabel {
				c = "\033[31;1m"
			}
			return fmt.Sprintf("%s %s%s\033[0m %s \033[33m%s\033[0m %s", tm, c, level, label, caller, msg)
		})))
	storage := chat.New(chat.InitialLastID())
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public_html")))
	mux.Handle("/api/publish", &chat.PublishHandler{Storage: storage})
	mux.Handle("/api/poll", &chat.PollHandler{Storage: storage})
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
