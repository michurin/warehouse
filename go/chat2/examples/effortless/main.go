package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/michurin/warehouse/go/chat2/handler"
	"github.com/michurin/warehouse/go/chat2/stream"
	"github.com/michurin/warehouse/go/chat2/text"
)

func validator(raw []byte) ([]byte, error) {
	in := ""
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return nil, err
	}
	cleaned := text.SanitizeText(in, 1000, "")
	if cleaned == "" {
		return nil, errors.New("no text")
	}
	return json.Marshal(cleaned)
}

func bindAddr() string {
	if len(os.Args) == 2 {
		return os.Args[1]
	}
	return ":8080"
}

func main() {
	logger := log.Default()
	addr := bindAddr()
	strm := stream.New(10)
	http.Handle("/", http.FileServer(http.Dir("examples/effortless/htdocs")))
	http.HandleFunc("/kit.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "js/kit.js") })
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "htdocs/favicon.ico") })
	http.Handle("/pub", handler.Pub(logger, strm, validator))
	http.Handle("/sub", handler.Sub(logger, strm, 10*time.Second))
	log.Printf("Listing on %s", addr)
	http.ListenAndServe(addr, nil)
}
