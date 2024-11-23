package jf

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type ctxHandler struct {
	next     slog.Handler
	fetchers []source
}

type source struct {
	source any
	key    string
}

func New(h slog.Handler, pairs ...any) slog.Handler {
	h, err := Wrap(h, pairs...)
	if err != nil {
		panic(err)
	}
	return h
}

func Wrap(h slog.Handler, pairs ...any) (slog.Handler, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("odd number of key/source pairs: %d", len(pairs))
	}
	f := make([]source, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		x, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("key must be a string: %[1]T: %[1]v", pairs[i])
		}
		f = append(f, source{source: pairs[i+1], key: x})
	}
	return &ctxHandler{next: h, fetchers: f}, nil
}

func (c *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.next.Enabled(ctx, level)
}

func (c *ctxHandler) Handle(ctx context.Context, record slog.Record) error {
	a := []any(nil)
	for _, f := range c.fetchers {
		v := ctx.Value(f.source)
		if v != nil {
			a = append(a, f.key, v)
		}
	}
	record.Add(a...)
	return c.next.Handle(ctx, record)
}

func (c *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ctxHandler{
		next:     c.next.WithAttrs(attrs),
		fetchers: c.fetchers,
	}
}

func (c *ctxHandler) WithGroup(name string) slog.Handler {
	return &ctxHandler{
		next:     c.next.WithGroup(name),
		fetchers: c.fetchers,
	}
}

type ctxError struct {
	next error
	ctx  context.Context
}

func (c *ctxError) Error() string {
	return c.next.Error()
}

func (c *ctxError) Unwrap() error {
	return c.next
}

func E(ctx context.Context, err error) error {
	if err == nil {
		return nil // safe to use with nil's
	}
	if e := new(ctxError); errors.As(err, &e) {
		return e // prevent double wrapping
	}
	return &ctxError{next: err, ctx: ctx}
}

func C(ctx context.Context, err error) context.Context {
	if e := new(ctxError); errors.As(err, &e) {
		return e.ctx // TODO: fuse t.ctx and ctx together?
	}
	return ctx
}
