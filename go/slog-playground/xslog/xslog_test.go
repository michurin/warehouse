package xslog_test

import (
	"context"
	"log/slog"
	"os"

	"slogplayground/xslog"
)

func ExampleHandler_usualUsecase() {
	// Somewhere you create handler.

	baseHandler := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr { // remove time; just to be reproducible
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	// You can setup custom attrs for handler. Our wrapper won't manage that attrs.

	baseHandler = baseHandler.WithAttrs([]slog.Attr{slog.Any("app", "one")})

	// Now you are able to setup global logger. You can setup lib-wide or application-wide logger using slog.SetDefault()

	log := slog.New(xslog.NewHandler(baseHandler, "xslog/xslog_test.go"))

	// You may have a chain of calls in you apps, let's say next two funcs.

	funcClient := func(ctx context.Context, arg int) error {
		ctx = xslog.Add(ctx, "client", "clientLabel", "arg", arg)
		if arg < 0 {
			return xslog.Errorf(ctx, "client error: invalid arg")
		}
		return nil
	}

	funcHandler := func(ctx context.Context, input int) error {
		ctx = xslog.Add(ctx, "component", "handlerLabel")
		err := funcClient(ctx, input)
		if err != nil {
			return xslog.Errorf(ctx, "handler failure: %w", err)
		}
		return nil
	}

	// You instrumentation is able to setup context and call the chain

	ctx := context.Background()

	ctx = xslog.Add(ctx, "request_id", "deadbeef")

	err := funcHandler(ctx, -1) // -1 will cause error
	if err != nil {
		log.Error("Error", err)
	}

	// output:
	// level=ERROR msg=Error app=one source=xslog/xslog_test.go:60 err_source=xslog/xslog_test.go:38 err_msg="handler failure: client error: invalid arg" request_id=deadbeef component=handlerLabel client=clientLabel arg=-1
}
