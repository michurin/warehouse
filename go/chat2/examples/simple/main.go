package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/michurin/warehouse/go/chat2/handler"
	"github.com/michurin/warehouse/go/chat2/stream"
	"github.com/michurin/warehouse/go/chat2/text"
)

var reColorStr = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

func validator(raw []byte) ([]byte, error) { // slightly oversimplified approach
	in := map[string]string{}
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return nil, err
	}
	color := in["color"]
	if !reColorStr.MatchString(color) {
		return nil, errors.New("invalid color")
	}
	return json.Marshal(map[string]string{
		"name":  text.SanitizeText(in["name"], 10, "[noname]"),
		"text":  text.SanitizeText(in["text"], 1000, "[nomessage]"),
		"color": color,
	})
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
	http.Handle("/", http.FileServer(http.Dir("examples/simple/htdocs")))
	http.HandleFunc("/kit.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "js/kit.js") })
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "htdocs/favicon.ico") })
	http.Handle("/pub", handler.Pub(logger, strm, validator))
	http.Handle("/sub", handler.Sub(logger, strm, 10*time.Second))
	log.Printf("Listing on %s", addr)
	http.ListenAndServe(addr, nil)
}
