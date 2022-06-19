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
	Get(ctx context.Context, bound int) ([][]byte, int)
}

type logger interface {
	Printf(format string, v ...interface{})
}

var errNotAllowed = errors.New("Method not allowed")

func handler(log logger, f func(context.Context, []byte) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK // default
		var resp []byte
		var body []byte
		var err error
		defer func() {
			if err == nil {
				log.Printf("%s %s %d %s -> %s", r.Method, r.URL.String(), status, string(body), string(resp))
			} else {
				log.Printf("%s %s %d %s: %s", r.Method, r.URL.String(), status, string(body), err)
			}
		}()
		if r.Method != http.MethodPost {
			err = errNotAllowed
			goto fin
		}
		body, err = io.ReadAll(r.Body)
		if err != nil {
			goto fin
		}
		resp, err = f(r.Context(), body)
		if err != nil {
			goto fin
		}
	fin:
		if err == errNotAllowed {
			status = http.StatusMethodNotAllowed
			resp = []byte("Method not allowed")
		} else if err != nil {
			status = http.StatusInternalServerError
			resp = []byte("Internal server error")
		}
		w.WriteHeader(status)
		w.Write(resp)
		w.Write([]byte{13, 10})
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
		a, b := strm.Get(ctx, req.Bound)
		m := make([]json.RawMessage, len(a))
		for i, v := range a {
			m[i] = json.RawMessage(v)
		}
		bodyRes, err := json.Marshal(subResponseDTO{
			Bound:    b,
			Messages: m,
		})
		if err != nil {
			return nil, fmt.Errorf("sub: cannot marshal: %w", err)
		}
		return bodyRes, nil
	})
}
