package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"unicode"

	"github.com/michurin/warehouse/go/chat2/handler"
	"github.com/michurin/warehouse/go/chat2/stream"
)

func validator(raw []byte) ([]byte, error) {
	in := map[string]string{}
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return nil, err
	}
	name, ok := in["name"]
	if !ok {
		return nil, errors.New("no name")
	}
	text, ok := in["text"]
	if !ok {
		return nil, errors.New("no text")
	}
	name = sanitize(name, 10)
	if name == "" {
		return nil, errors.New("name is empty")
	}
	text = sanitize(text, 1000)
	if text == "" {
		return nil, errors.New("text is empty")
	}
	return json.Marshal(map[string]string{
		"name": name,
		"text": text,
	})
}

// truncate, remove ctrl chars and collapse spaces
func sanitize(a string, n int) string {
	r := []rune(nil)
	s := false
	for _, c := range a {
		if !unicode.IsPrint(c) {
			continue
		}
		if unicode.IsSpace(c) {
			s = len(r) > 0 // skip spaces at the begging
			continue
		}
		if s {
			if len(r) >= n-1 { // drop spaces at last position
				break
			}
			r = append(r, '\u0020')
			s = false
		}
		r = append(r, c)
		if len(r) >= n {
			break
		}
	}
	return string(r)
}

func main() {
	logger := log.Default()
	strm := stream.New(10)
	http.Handle("/", http.FileServer(http.Dir("htdocs")))
	http.HandleFunc("/pub", handler.Pub(logger, strm, validator))
	http.HandleFunc("/sub", handler.Sub(logger, strm, 10*time.Second))
	http.ListenAndServe(":8080", nil)
}
