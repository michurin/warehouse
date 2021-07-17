package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/michurin/minlog"

	"github.com/michurin/warehouse/go/chat/pkg/chat"
)

func wrap(label string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// You may want to validate session/cookies here
		// or log something
		// or add things like client IP to context for contextual logging.
		//
		// It is just very simple example
		next.ServeHTTP(w, r.WithContext(minlog.Label(r.Context(), label)))
	})
}

func NewPublishingHandler(rooms *chat.Rooms, validator func(json.RawMessage) error, log minlog.Interface, label string) http.Handler {
	return wrap(label+":pub", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rid, msg, err := chat.DecodePublishingRequest(r.Body)
		if err != nil {
			log.Log(ctx, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx = minlog.Label(ctx, "room:"+rid)
		log.Log(ctx, "Publish:", []byte(msg))
		if err = validator(msg); err != nil {
			log.Log(ctx, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rooms.Pub(ctx, rid, msg)
		chat.SetContentType(w.Header())
		err = chat.EncodePublishingResponse(w)
		if err != nil {
			log.Log(ctx, err)
		}
	}))
}

func NewPollingHandler(rooms *chat.Rooms, log minlog.Interface, label string) http.Handler {
	return wrap(label+":poll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		rid, id, err := chat.DecodePollingRequest(r.Body)
		if err != nil {
			log.Log(ctx, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx = minlog.Label(ctx, "room:"+rid)
		mm, lastID := rooms.Fetch(ctx, rid, id)
		chat.SetContentType(w.Header())
		err = chat.EncodePollingResponse(w, mm, lastID)
		if err != nil {
			log.Log(ctx, err)
		}
	}))
}

func NewMonitoringHandler(rooms *chat.Rooms) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		lst := chat.RoomsList(rooms)
		if len(lst) == 0 {
			w.Write([]byte("(no rooms)"))
			return
		}
		sections := make([]string, len(lst))
		i := 0
		for k, v := range lst {
			q := make([]string, len(v))
			for i, x := range v {
				q[i] = string(x)
			}
			sections[i] = k + "\n\n" + strings.Join(q, "\n")
			i++
		}
		sort.Strings(sections)
		w.Write([]byte(strings.Join(sections, "\n\n")))
	})
}
