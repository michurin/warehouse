package argctx

import (
	"context"
)

type ctxKeyT int

const ctxKey = ctxKeyT(0)

func With(ctx context.Context, args ...any) context.Context { // With mimics [slog.With]
	if len(args) == 0 {
		return ctx
	}
	a, _ := ctx.Value(ctxKey).([]any) // a can be nil, it is not crime
	b := make([]any, len(a), len(a)+len(args))
	copy(b, a)
	b = append(b, args...)
	return context.WithValue(ctx, ctxKey, b)
}

func Args(ctx context.Context) []any {
	a, _ := ctx.Value(ctxKey).([]any)
	return grouping(a)
}

func grouping(a []any) []any {
	r := []any(nil)
	for i, v := range a {
		if f, ok := v.(func([]any) any); ok {
			r = append(r, f(grouping(a[i+1:])))
			break
		}
		r = append(r, v)
	}
	return r
}
