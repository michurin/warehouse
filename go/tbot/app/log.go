package app

import (
	"os"
	"runtime"
	"strings"

	"github.com/michurin/minlog"

	"github.com/michurin/cnbot/app/aw"
)

func color(next minlog.FieldFunc, colorCode string) minlog.FieldFunc { // TODO move it to minlog package?
	return func(r minlog.Record) string {
		t := next(r)
		if t == "" {
			return t
		}
		return "\033[" + colorCode + "m" + next(r) + "\033[0m"
	}
}

func prefix(next minlog.FieldFunc, prefix string) minlog.FieldFunc { // TODO move it to minlog package?
	return func(r minlog.Record) string {
		t := next(r)
		if t == "" {
			return t
		}
		return prefix + next(r)
	}
}

func SetupLogging() {
	_, file, _, _ := runtime.Caller(0)
	pfx := strings.TrimSuffix(file, "app/log.go")
	opts := []minlog.Option(nil)
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		opts = []minlog.Option{
			minlog.WithFields(
				minlog.FieldLevel("", "\033[1;33;41m ERROR \033[0m"),
				color(minlog.FieldCaller(pfx), "1;34"),
				color(minlog.FieldErrorCaller(pfx), "1;31"),
				color(minlog.FieldNamed("comp"), "32"),
				color(minlog.FieldNamed("bot"), "35"),
				color(minlog.FieldNamed("api"), "1;35"),
				color(minlog.FieldNamed("user"), "1;32"),
				prefix(color(minlog.FieldNamed("pid"), "33"), "PID:"),
				minlog.FieldFallbackKV("api", "bot", "comp", "pid", "user"),
				minlog.FieldMessage()),
		}
	} else {
		opts = []minlog.Option{
			minlog.WithFields(
				minlog.FieldLevel("[I]", "[E]"),
				minlog.FieldCaller(pfx),
				minlog.FieldErrorCaller(pfx),
				minlog.FieldFallbackKV(),
				minlog.FieldMessage()),
		}
	}
	aw.Log = minlog.New(opts...).Log
}
