package httppost

import (
	"context"
	"errors"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

type logger interface {
	Printf(format string, v ...interface{})
}

var errorMethodNotAllowed = errors.New("Method not allowed")

func Handler(log logger, f func(context.Context, []byte) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp []byte
		var body []byte
		var err error
		start := time.Now()
		clientIP := r.Header.Get("X-Real-IP") // TODO all logging has to be in middleware; or logging has to be smarter and consider ctx
		if clientIP == "" {
			clientIP = "-"
		}
		defer func() {
			dt := time.Now().Sub(start).Milliseconds()
			if rec := recover(); rec != nil {
				log.Printf("%s %dms %s %s: Panic: %v\n%s\n", clientIP, dt, r.Method, r.URL.String(), rec, debug.Stack())
				return
			}
			if err != nil {
				log.Printf("%s %dms %s %s: Error: %s", clientIP, dt, r.Method, r.URL.String(), err)
				return
			}
			log.Printf("%s %dms %s %s: %s -> %s", clientIP, dt, r.Method, r.URL.String(), string(body), string(resp))
		}()
		if r.Method != http.MethodPost {
			err = errorMethodNotAllowed // for logging in defer
			http.Error(w, err.Error(), http.StatusMethodNotAllowed)
			return
		}
		body, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err = f(r.Context(), body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if resp == nil {
			resp = []byte(`{}`) // force valid JSON, hackish, in perfect world we wouldn't change data here
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		w.Write([]byte{13, 10}) // just to be curl and command line friendly
	}
}
