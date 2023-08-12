package xslog

import (
	"context"
)

type ctxKeyT int

const ctxKey = ctxKeyT(0)

func Add(ctx context.Context, x ...any) context.Context {
	rx := ctx.Value(ctxKey)
	if ox, ok := rx.([][]any); ok {
		return context.WithValue(ctx, ctxKey, append(ox, x))
	}
	return context.WithValue(ctx, ctxKey, [][]any{x})
}
