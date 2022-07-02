package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type stream interface {
	Put(x []byte)
	Get(ctx context.Context, bound uint64) ([][]byte, uint64, bool)
}

type logger interface {
	Printf(format string, v ...interface{})
}

var errorMethodNotAllowed = errors.New("Method not allowed")

func handler(log logger, f func(context.Context, []byte) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp []byte
		var body []byte
		var err error
		defer func() {
			if err == nil {
				log.Printf("%s %s: %s -> %s", r.Method, r.URL.String(), string(body), string(resp))
			} else {
				log.Printf("%s %s: Error: %s", r.Method, r.URL.String(), err)
			}
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

type subRequestDTO struct {
	Bound int `json:"bound"`
}

type subResponseDTO struct {
	Bound    int               `json:"bound"`
	Messages []json.RawMessage `json:"messages"`
}

func Pub(log logger, strm stream, validator func(raw []byte) ([]byte, error)) http.HandlerFunc {
	return handler(log, func(ctx context.Context, body []byte) ([]byte, error) {
		msg, err := validator(body)
		if err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}
		strm.Put(msg)
		return nil, nil
	})
}

func Sub(log logger, strm stream, timeout time.Duration) http.HandlerFunc {
	return handler(log, func(ctx context.Context, body []byte) ([]byte, error) {
		req := subRequestDTO{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			return nil, fmt.Errorf("sub: cannot unmarshal: %w", err)
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		a, b, _ := strm.Get(ctx, uint64(req.Bound)) // We just ignore continuity flag
		m := make([]json.RawMessage, len(a))
		for i, v := range a {
			m[i] = json.RawMessage(v)
		}
		bodyRes, err := json.Marshal(subResponseDTO{
			Bound:    int(b), // In fact JS has limit Number.MAX_SAFE_INTEGER=2*53-1
			Messages: m,
		})
		if err != nil {
			return nil, fmt.Errorf("sub: cannot marshal: %w", err)
		}
		return bodyRes, nil
	})
}
