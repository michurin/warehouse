package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"sse/handler"
	"sse/internal/xlog"
	"sse/loggingmw"
	"sse/room"
)

func main() {
	slog.SetDefault(slog.New(xlog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
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
	house := room.New()
	go handler.RevisionLoop(house)
	err := http.ListenAndServe(":7011", loggingmw.MW(handler.Handler(house)))
	if err != nil {
		log.Printf("Listener error: %s", err.Error())
	}
}
