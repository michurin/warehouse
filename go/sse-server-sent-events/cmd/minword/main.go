package main

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"sse/loggingmw"
	"sse/room"
	"sse/static"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func handleStatic(fsh http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Add("Cache-Control", "no-cache")
		fsh.ServeHTTP(w, r)
	}
}

func handleFetch(ch *room.House) http.HandlerFunc {
	// TODO process errors, TODO use Copy
	return func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body) // just drop body. We do not need to close it. Oh. It works without ctx
		roomID, userID := roomAndUser(r.URL.Query())
		log.Printf("room/user = %s/%s", roomID, userID)
		leid, err := strconv.ParseInt(r.Header.Get("Last-Event-Id"), 10, 64)
		if err != nil {
			leid = 0
		}
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 28*time.Second)
		defer cancel()
		h := w.Header()
		h.Add("X-Accel-Buffering", "no")
		h.Add("Content-Type", "text/event-stream")
		h.Add("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		messages := [][]byte(nil) // we have to create this var out of the loop, as leid
		for {                     // TODO check writing errors
			messages, leid = ch.Fetch(ctx, roomID, userID, leid)
			if ctx.Err() != nil {
				return
			}
			w.Write([]byte("event: message\n"))                          // it MUST be `message` to make e.onmessage be fired
			w.Write([]byte("retry: 200\n"))                              // server side control for reconnecting delay
			w.Write([]byte("id: " + strconv.FormatInt(leid, 10) + "\n")) // it will be `Last-Event-Id: TOKEN` (on request)
			for _, m := range messages {
				w.Write([]byte("data: "))
				w.Write(m) // we are storing single line messages only
				w.Write([]byte{10})
			}
			w.Write([]byte{10})
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func handlePub(ch *room.House) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)
		_, err := io.Copy(b, r.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		text := strings.Map(func(x rune) rune {
			if unicode.IsControl(x) { // clean up \n as well, useful in JSON sanitizing perspective
				return '\x20'
			}
			return x
		}, strings.TrimSpace(b.String()))
		if len(text) < 2 {
			http.Error(w, "Error", http.StatusBadRequest)
			return
		}
		text = text[:len(text)-1] + `,"ts":` + strconv.FormatInt(time.Now().UnixMilli(), 10) + "}"
		roomID, userID := roomAndUser(r.URL.Query())
		ch.Pub(roomID, userID, []byte(text))
		w.WriteHeader(http.StatusOK)
	}
}

func handler(staticFS fs.FS, house *room.House) http.HandlerFunc {
	fsh := http.FileServerFS(staticFS)
	fetchh := handleFetch(house)
	pubh := handlePub(house)
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.EscapedPath()
		switch r.Method {
		case http.MethodGet:
			switch path {
			case "/fetch":
				fetchh.ServeHTTP(w, r)
				return
			}
			fsh.ServeHTTP(w, r)
			return
		case http.MethodPost:
			switch path {
			case "/pub":
				pubh.ServeHTTP(w, r)
				return
			case "/lock":
				http.Error(w, "OK", http.StatusOK) // TODO
				return
			case "/unlock":
				http.Error(w, "OK", http.StatusOK) // TODO
				return
			}
		default:
			http.Error(w, "not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func roomAndUser(v url.Values) (string, string) {
	roomID := v.Get("room")
	userID := v.Get("user")
	// TODO validate, set defaults
	return roomID, userID
}

func main() {
	house := room.New()
	err := http.ListenAndServe(":7011", http.MaxBytesHandler(loggingmw.MW(handler(static.FS, house)), 4000))
	if err != nil {
		log.Printf("Listener error: %s", err.Error())
	}
}
