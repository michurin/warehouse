package xlog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sync/atomic"
	"time"
)

var defaultLogger atomic.Pointer[slog.Logger]

func init() { //nolint:gochecknoinits
	defaultLogger.Store(slog.Default())
}

// SetDefault mimics slog.SetDefault
func SetDefault(l *slog.Logger) {
	defaultLogger.Store(l)
}

// L is botwide logging function, it could be private
func L(ctx context.Context, a any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip
	r := slog.Record{}
	switch v := a.(type) {
	case error:
		r = slog.NewRecord(time.Now(), slog.LevelError, v.Error(), pcs[0])
		r.Add("error", v) // it will be skipped in ctxlog.Handler
	case string:
		r = slog.NewRecord(time.Now(), slog.LevelInfo, v, pcs[0])
	case []byte:
		r = slog.NewRecord(time.Now(), slog.LevelInfo, safeString(v), pcs[0])
	case nil:
		r = slog.NewRecord(time.Now(), slog.LevelInfo, "<nil>", pcs[0])
	default:
		r = slog.NewRecord(time.Now(), slog.LevelWarn, fmt.Sprintf("%[1]T: %#[1]v", a), pcs[0])
	}
	_ = defaultLogger.Load().Handler().Handle(ctx, r)
}
