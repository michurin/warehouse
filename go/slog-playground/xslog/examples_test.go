package xslog_test

import (
	"context"
	"log/slog"
	"os"

	"slogplayground/xslog"
)

const thisFileName = "xslog/examples_test.go"

var optsNoTimeNoSourceNoLevel = slog.HandlerOptions{
	AddSource: false,
	Level:     nil,
	ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr { // remove time; just to be reproducible
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	},
}

func ExampleHandler_usualUsecase() {
	// Somewhere you create handler.

	baseHandler := slog.Handler(slog.NewTextHandler(os.Stdout, &optsNoTimeNoSourceNoLevel))

	// You can setup custom attrs for handler. Our wrapper won't manage that attrs.

	baseHandler = baseHandler.WithAttrs([]slog.Attr{slog.Any("app", "one")})

	// Now you are able to setup global logger. You can setup lib-wide or application-wide logger using slog.SetDefault()

	log := slog.New(xslog.Handler(baseHandler, thisFileName))

	// You may have a chain of calls in you apps, let's say next two funcs.

	funcContextFreeLogic := func() error {
		return xslog.Errorf("initial error")
	}

	funcClient := func(ctx context.Context, arg int) error {
		ctx = xslog.Add(ctx, "client", "clientLabel", "arg", arg)
		err := funcContextFreeLogic()
		if err != nil {
			return xslog.Errorfx(ctx, "client error: %w", err)
		}
		return nil
	}

	funcHandler := func(ctx context.Context, input int) error {
		ctx = xslog.Add(ctx, "component", "handlerLabel")
		err := funcClient(ctx, input)
		if err != nil {
			return xslog.Errorfx(ctx, "handler failure: %w", err)
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
	// level=ERROR msg=Error app=one source=xslog/examples_test.go:69 err_source=xslog/examples_test.go:40 err_msg="handler failure: client error: initial error" request_id=deadbeef component=handlerLabel client=clientLabel arg=-1
}

func ExampleHandler_howGroupsAndAttrsDoing() {
	baseHandler := slog.Handler(slog.NewTextHandler(os.Stdout, &optsNoTimeNoSourceNoLevel))

	log := slog.New(xslog.Handler(baseHandler, thisFileName))
	log.Info("Message")
	log.Info("Message-inline-attrs", "P", "Q")
	log.InfoContext(xslog.Add(context.Background(), "V", "W"), "Message-1-ctx-attrs")
	log = log.With("X", "Y")
	log.Info("Message-with-attrs")
	log = log.WithGroup("G")
	log.Info("Message-with-group")

	// output:
	// level=INFO msg=Message source=xslog/examples_test.go:80
	// level=INFO msg=Message-inline-attrs source=xslog/examples_test.go:81 P=Q
	// level=INFO msg=Message-1-ctx-attrs source=xslog/examples_test.go:82 V=W
	// level=INFO msg=Message-with-attrs X=Y source=xslog/examples_test.go:84
	// level=INFO msg=Message-with-group X=Y G.source=xslog/examples_test.go:86
}
