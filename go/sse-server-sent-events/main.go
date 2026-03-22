package main

import (
	"bytes"
	"container/list"
	"context"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

//go:embed static
var emgedFS embed.FS

func handleStatic(fsh http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Add("Cache-Control", "no-cache")
		fsh.ServeHTTP(w, r)
	}
}

func handleFetch(ch *room) http.HandlerFunc {
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
		for {
			messages, leid = ch.fetch(ctx, leid)
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

func handleSend(ch *room) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)
		_, err := io.Copy(b, r.Body)
		if err != nil {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
		for e := range bytes.FieldsFuncSeq(b.Bytes(), func(r rune) bool { return r == 10 || r == 13 }) {
			ch.pub(e)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	fsst, err := fs.Sub(emgedFS, "static")
	if err != nil {
		log.Panic(err)
	}
	fsh := http.FileServerFS(fsst)

	ch := &room{
		lastID: time.Now().UnixNano(), // let lastID grow between restart (naive)
		wall:   list.New(),
		lock:   new(sync.Mutex),
		signal: make(chan struct{}),
	}

	http.HandleFunc("/", handleStatic(fsh))
	http.HandleFunc("/fetch", handleFetch(ch))
	http.HandleFunc("/send", handleSend(ch))
	err := http.ListenAndServe(":7011", nil)
	if err != nil {
		log.Printf("Listener error: %s", err.Error())
	}
}

// --- Room

type room struct {
	lastID int64
	wall   *list.List
	lock   *sync.Mutex
	signal chan struct{}
}

func (r *room) pub(m []byte) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.lastID++
	r.wall.PushFront(m)
	for r.wall.Len() > 10 {
		r.wall.Remove(r.wall.Back())
	}
	close(r.signal)
	r.signal = make(chan struct{})
}

func (r *room) syncFetch(lastID int64) ([][]byte, chan struct{}, int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	w := [][]byte(nil)
	i := r.lastID
	l := lastID
	for e := r.wall.Front(); e != nil; e = e.Next() {
		if i <= lastID {
			break
		}
		if len(w) == 0 {
			l = i
		}
		w = append(w, e.Value.([]byte))
		i--
	}
	return w, r.signal, l
}

func (r *room) fetch(ctx context.Context, lastID int64) ([][]byte, int64) {
	m, c, id := r.syncFetch(lastID)
	if len(m) > 0 {
		return m, id
	}
	select {
	case <-ctx.Done():
		return nil, lastID
	case <-c:
		m, _, id := r.syncFetch(lastID)
		return m, id
	}
}
