package xlog_test

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/michurin/cnbot/xlog"
)

func ExampleL() {
	// Oh. it global and can ruin other tests
	// However, it could be good idea to setup all tests logging in the same way?
	xlog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})))
	ctx := context.Background()
	xlog.L(ctx, "ok")
	xlog.L(ctx, errors.New("err"))
	// Output:
	// level=INFO msg=ok
	// level=ERROR msg=err error=err
}
