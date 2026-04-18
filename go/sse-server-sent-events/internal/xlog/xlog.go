package xlog

import (
	"context"
	"log/slog"
	"math/rand"
)

type xlogKeyT int

const xlogKey xlogKeyT = 0

type Handler struct {
	next slog.Handler
}

func New(next slog.Handler) *Handler {
	return &Handler{next: next}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, rec slog.Record) error {
	v, ok := ctx.Value(xlogKey).([]any)
	if ok {
		rec.Add(v...)
	}
	return h.next.Handle(ctx, rec)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("not implemented")
}

func (h *Handler) WithGroup(name string) slog.Handler {
	panic("not implemented")
}

func with(ctx context.Context, args ...any) context.Context {
	if len(args) == 0 {
		return ctx
	}
	v, ok := ctx.Value(xlogKey).([]any)
	if ok {
		v = append(v, args...)
	} else {
		v = args
	}
	return context.WithValue(ctx, xlogKey, v)
}

func WithRequestID(ctx context.Context) context.Context {
	return with(ctx, slog.Int("request_id", rand.Intn(10000)))
}

func WithAddr(ctx context.Context, addr string) context.Context {
	return with(ctx, slog.String("addr", addr))
}

func WithMethod(ctx context.Context, method string) context.Context {
	return with(ctx, slog.String("method", method))
}

func WithURL(ctx context.Context, url string) context.Context {
	return with(ctx, slog.String("url", url))
}

func WithLocation(ctx context.Context, location string) context.Context {
	return with(ctx, slog.String("location", location))
}

func WithStatus(ctx context.Context, status int) context.Context {
	return with(ctx, slog.Int("status", status))
}

func WithRoom(ctx context.Context, room string) context.Context {
	return with(ctx, slog.String("room", room))
}

func WithUser(ctx context.Context, user string) context.Context {
	return with(ctx, slog.String("user", user))
}
