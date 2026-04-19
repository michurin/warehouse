package xlog_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/michurin/minchat/internal/xlog"
)

func Example() {
	ctx := context.Background()
	h := xlog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     nil,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue("10:00")
			}
			return a
		},
	}))
	l := slog.New(h)
	l.InfoContext(ctx, "OK")
	ctx = xlog.WithAddr(ctx, "127.0.0.1")
	ctx = xlog.WithRoom(ctx, "hole")
	ctx = xlog.WithUser(ctx, "person")
	l.InfoContext(ctx, "OK")
	// output:
	// time=10:00 level=INFO msg=OK
	// time=10:00 level=INFO msg=OK addr=127.0.0.1 room=hole user=person
}
