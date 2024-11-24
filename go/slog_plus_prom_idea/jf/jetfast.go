package jf

import (
	"context"
	"errors"
	"log/slog"
)

type ctxHandler struct {
	next     slog.Handler
	fetchers []func(context.Context) []any
}

type source struct {
	source any
	key    string
}

type errKeyT int

const errKey errKeyT = iota

func New(h slog.Handler, f ...func(context.Context) []any) slog.Handler {
	return &ctxHandler{next: h, fetchers: f}
}

func (c *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.next.Enabled(ctx, level)
}

func (c *ctxHandler) Handle(ctx context.Context, record slog.Record) error {
	r := record.Clone()
	for _, f := range c.fetchers {
		r.Add(f(ctx)...)
	}
	return c.next.Handle(ctx, r)
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

func F(field string, key any) func(context.Context) []any {
	return func(ctx context.Context) []any {
		if x := ctx.Value(key); x != nil {
			return []any{field, x}
		}
		return nil
	}
}

func ErrF(f func(context.Context, error) []any) func(context.Context) []any {
	return func(ctx context.Context) []any {
		if e, ok := ctx.Value(errKey).(error); ok && e != nil {
			return f(ctx, e)
		}
		return nil
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
		return context.WithValue(e.ctx, errKey, e.next)
	}
	return ctx
}
