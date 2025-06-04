package sslog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"

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
	ctx = argctx.With(ctx, slog.Group("g", slog.String("kk", "vv")))
	log.Info("OK", argctx.Args(ctx)...)

	err := errors.New("message")
	err = argerr.Wrap(err, argctx.Args(ctx)...)
	log.Error("ERR", argerr.Args(err)...)

	// output:
	// level=INFO msg=OK K=V g.kk=vv
	// level=ERROR msg=ERR K=V g.kk=vv
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
