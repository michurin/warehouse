package xlog

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func Init() {
	slog.SetDefault(slog.New(New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02_15:04:05"))
		}
		if a.Key == slog.SourceKey {
			s := a.Value.Any().(*slog.Source)
			i := strings.LastIndexByte(s.File, '/')
			if i > 0 {
				a.Value = slog.StringValue(fmt.Sprintf("%s:%d", s.File[i+1:], s.Line)) // s.Function?
			}
		}
		return a
	}}))))
}
