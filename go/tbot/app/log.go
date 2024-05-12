package app

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"sort"
	"time"

	"github.com/michurin/cnbot/app/aw"
	"github.com/michurin/cnbot/ctxlog"
)

// logHandler implements interface slog.Handler
// it is drop-in replacement for slog.NewTextHandler, but more human friendly
type logHandler struct{}

func (logHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (logHandler) Handle(_ context.Context, r slog.Record) error {
	kv := map[string]any{} // not thread safe, however r.Attrs works consequently
	r.Attrs(func(a slog.Attr) bool {
		kv[a.Key] = a.Value.Any()
		return true
	})
	std := ""                                                    // std attributes
	for _, a := range []string{"bot", "comp", "api", "source"} { // order significant
		if v, ok := kv[a]; ok {
			std = std + " [" + v.(string) + "]"
			delete(kv, a)
		}
	}
	ekeys := []string(nil) // extra keys
	for k := range kv {
		ekeys = append(ekeys, k)
	}
	sort.Strings(ekeys)
	nstd := ""
	for _, a := range ekeys {
		nstd += fmt.Sprintf(" %s=%v", a, kv[a])
	}
	fmt.Printf("%s%s%s %s\n", r.Time.Format("2006-01-02 15:04:05"), std, nstd, r.Message)
	return nil
}

func (h logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("NOT IMPLEMENTED")
}

func (logHandler) WithGroup(name string) slog.Handler {
	panic("NOT IMPLEMENTED")
}

func SetupLogging() {
	l := slog.New(ctxlog.Handler(logHandler{}, "app/log.go"))
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
