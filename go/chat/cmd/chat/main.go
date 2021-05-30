package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/michurin/warehouse/go/chat/pkg/chat"
)

func log(x ...interface{}) {
	fmt.Println(x...)
}

func main() {
	storage := chat.New()
	mux := http.NewServeMux()
	//mux.Handle("/", http.StripPrefix("/public_http/", http.FileServer(http.Dir("public_http"))))
	mux.Handle("/", http.FileServer(http.Dir("public_html")))
	mux.Handle("/api/publish", &chat.PublishHandler{Storage: storage})
	mux.Handle("/api/poll", &chat.PollHandler{Storage: storage})
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    300 * time.Second, // 300 is most browsers default
		WriteTimeout:   300 * time.Second,
		MaxHeaderBytes: 1 << 12,
	}
	err := s.ListenAndServe()
	if err != nil {
		log(err)
	}
}
