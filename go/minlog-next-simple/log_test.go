package sslog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"sslog"
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
	ctx = sslog.With(ctx, "K", "V")
	ctx = sslog.With(ctx, slog.Group("g", "kk", "vv"))
	log.Info("OK", sslog.Args(ctx)...)

	err := errors.New("message")
	err = sslog.Wrap(err, sslog.Args(ctx)...)
	log.Error("ERR", sslog.ArgsE(err)...)

	// output:
	// level=INFO msg=OK K=V g.kk=vv
	// level=ERROR msg=ERR K=V g.kk=vv
}
