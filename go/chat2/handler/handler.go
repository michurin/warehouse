package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/michurin/warehouse/go/chat2/httppost"
)

type stream interface {
	Put(x []byte)
	Get(ctx context.Context, bound uint64) ([][]byte, uint64)
}

type logger interface {
	Printf(format string, v ...interface{})
}

type subRequestDTO struct {
	Bound uint64 `json:"bound"`
}

type subResponseDTO struct {
	Bound    uint64            `json:"bound"`
	Messages []json.RawMessage `json:"messages"`
}

func Pub(log logger, strm stream, validator func(raw []byte) ([]byte, error)) http.HandlerFunc {
	return httppost.Handler(log, func(ctx context.Context, body []byte) ([]byte, error) {
		msg, err := validator(body)
		if err != nil {
			return nil, fmt.Errorf("validation error: %w", err)
		}
		strm.Put(msg)
		return nil, nil
	})
}

func Sub(log logger, strm stream, timeout time.Duration) http.HandlerFunc {
	return httppost.Handler(log, func(ctx context.Context, body []byte) ([]byte, error) {
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
			Bound:    b, // In fact JS has limit Number.MAX_SAFE_INTEGER=2*53-1
			Messages: m,
		})
		if err != nil {
			return nil, fmt.Errorf("sub: cannot marshal: %w", err)
		}
		return bodyRes, nil
	})
}
