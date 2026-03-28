package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"sse/loggingmw"
	"sse/static"
	"sse/wall"
)

func handleStatic(fsh http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Add("Cache-Control", "no-cache")
		fsh.ServeHTTP(w, r)
	}
}

func handleFetch(ch *wall.Wall) http.HandlerFunc {
	// TODO process errors, TODO use Copy
	return func(w http.ResponseWriter, r *http.Request) {
		leid, err := strconv.ParseInt(r.Header.Get("Last-Event-Id"), 10, 64)
		if err != nil {
			leid = 0
		}
		io.Copy(io.Discard, r.Body) // just drop body. We do not need to close it. Oh. It works without ctx
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
			messages, leid = ch.Fetch(ctx, leid)
			if ctx.Err() != nil {
				return
			}
			w.Write([]byte("event: message\n"))                          // it MUST be `message` to make e.onmessage be fired
			w.Write([]byte("retry: 3000\n"))                             // server side control for reconnecting delay
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

func handleSend(ch *wall.Wall) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)
		_, err := io.Copy(b, r.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		for e := range bytes.FieldsFuncSeq(b.Bytes(), func(r rune) bool { return r == 10 || r == 13 }) {
			ch.Pub(e)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	fsh := http.FileServerFS(static.FS)

	ch := wall.New(time.Now().UnixNano())

	http.HandleFunc("/", handleStatic(fsh))
	http.HandleFunc("/fetch", handleFetch(ch))
	http.HandleFunc("/send", handleSend(ch))
	err := http.ListenAndServe(":7011", http.MaxBytesHandler(loggingmw.MW(http.DefaultServeMux), 4000))
	if err != nil {
		log.Printf("Listener error: %s", err.Error())
	}
}
