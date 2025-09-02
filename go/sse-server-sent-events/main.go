package main

import (
	"embed"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

//go:embed index.html
var emgedFS embed.FS

func main() {
	fsh := http.FileServerFS(emgedFS)

	err := http.ListenAndServe(":7011", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/" {
			h := w.Header()
			h.Add("Cache-Control", "no-cache")
			fsh.ServeHTTP(w, r)
			return
		}

		log.Print("Start: " + r.URL.String() + " LEID: " + r.Header.Get("last-event-id"))
		io.Copy(io.Discard, r.Body) // just drop body. We do not need to close it
		h := w.Header()
		h.Add("X-Accel-Buffering", "no")
		h.Add("Content-Type", "text/event-stream")
		h.Add("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		for range 4 {
			n := time.Now().UTC()
			w.Write([]byte("event: message\n"))                          // it MUST be `message` to make e.onmessage be fired
			w.Write([]byte("retry: 3000\n"))                             // server side control for reconnecting delay
			w.Write([]byte("id: " + strconv.Itoa(int(n.Unix())) + "\n")) // it will be `Last-Event-Id: TOKEN` (on request)
			w.Write([]byte("data: " + n.Format(time.RFC3339) + "\n\n"))
			if f, ok := w.(http.Flusher); ok {
				log.Print("flash")
				f.Flush()
			}
			log.Print(r.URL.String() + " (send)")
			time.Sleep(time.Second)
		}
	}))
	if err != nil {
		log.Print(err)
	}
}
