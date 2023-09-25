package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/michurin/cnbot/app/aw"
	"github.com/michurin/cnbot/ctxlog"
)

func SetupLogging() {
	l := slog.New(ctxlog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}), "app/log.go"))
	aw.L = func(ctx context.Context, a any) {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip
		r := slog.Record{}
		switch v := a.(type) {
		case error:
			r = slog.NewRecord(time.Now(), slog.LevelError, "Error", pcs[0])
			r.Add(v)
		case string:
			r = slog.NewRecord(time.Now(), slog.LevelInfo, v, pcs[0])
		default:
			r = slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf("%[1]T: %#[1]v", a), pcs[0])
		}
		_ = l.Handler().Handle(ctx, r)
	}
}
