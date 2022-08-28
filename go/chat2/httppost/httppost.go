package httppost

import (
	"context"
	"errors"
	"io"
	"net/http"
	"runtime/debug"
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
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("%s %s: Panic: %v\n%s\n", r.Method, r.URL.String(), rec, debug.Stack())
				return
			}
			if err != nil {
				log.Printf("%s %s: Error: %s", r.Method, r.URL.String(), err)
				return
			}
			log.Printf("%s %s: %s -> %s", r.Method, r.URL.String(), string(body), string(resp))
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
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
		w.Write([]byte{13, 10}) // just to be curl and command line friendly
	}
}
