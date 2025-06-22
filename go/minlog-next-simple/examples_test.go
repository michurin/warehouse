package sslog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"go.uber.org/zap"

	"github.com/michurin/minlog/argctx"
	"github.com/michurin/minlog/argerr"
)

var handlerOptions = &slog.HandlerOptions{
	AddSource: false,
	ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}

func Example() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	ctx = argctx.With(ctx, "K", "V")
	log.Info("OK", argctx.Args(ctx)...)

	err := errors.New("message")
	err = argerr.Wrap(err, argctx.Args(ctx)...) // CONSIDER: helper Wrap(err, ctx)?
	log.Error("ERR", argerr.Args(err)...)

	// output:
	// level=INFO msg=OK K=V
	// level=ERROR msg=ERR K=V
}

func Example_experimentalCtxGroups() {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stdout, handlerOptions))

	ctx = argctx.With(ctx, "R", 1)
	ctx = argctx.With(ctx, func(x []any) any {
		return slog.Group("handler", x...)
	})
	ctx = argctx.With(ctx, "H", 2)
	ctx = argctx.With(ctx, func(x []any) any {
		return slog.Group("adapter", x...)
	})
	ctx = argctx.With(ctx, "A", 3)

	err := errors.New("error message")
	err = argerr.Wrap(err, argctx.Args(ctx)...)

	log.Error(err.Error(), argerr.Args(err)...) // we are obtaining all logging context from error only
	log.Info("OK", argctx.Args(ctx)...)

	// output:
	// level=ERROR msg="error message" R=1 handler.H=2 handler.adapter.A=3
	// level=INFO msg=OK R=1 handler.H=2 handler.adapter.A=3
}

func Example_zapGroups() {
	ctx := context.Background()
	log := zap.NewExample()

	ctx = withInt(ctx, "R", 1)
	ctx = argctx.With(ctx, func(x []any) any {
		return zap.Dict("handler", zapCast(x)...)
	})
	ctx = withInt(ctx, "H", 2)
	ctx = argctx.With(ctx, func(x []any) any {
		return zap.Dict("adapter", zapCast(x)...)
	})
	ctx = withInt(ctx, "A", 3)

	err := errors.New("error message")
	err = argerr.Wrap(err, argctx.Args(ctx)...)

	log.Error(err.Error(), zapFields(err)...) // we are obtaining all logging context from error only
	log.Info("OK", zapFields(ctx)...)

	// output:
	// {"level":"error","msg":"error message","R":1,"handler":{"H":2,"adapter":{"A":3}}}
	// {"level":"info","msg":"OK","R":1,"handler":{"H":2,"adapter":{"A":3}}}
}

func zapFields(x any) []zap.Field { // CONSIDER: split into two well typed functions?
	switch t := x.(type) {
	case error:
		return zapCast(argerr.Args(t))
	case context.Context:
		return zapCast(argctx.Args(t))
	}
	return nil // CONSIDER: this case should not exist
}

func zapCast(a []any) []zap.Field {
	f := []zap.Field(nil)
	for _, x := range a {
		if t, ok := x.(zap.Field); ok {
			f = append(f, t)
		}
	}
	return f
}

func withInt(ctx context.Context, key string, val int) context.Context { // CONSIDER: collection of wrappers? Generics?
	return argctx.With(ctx, zap.Int(key, val))
}
