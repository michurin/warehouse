package app

import (
	"os"
	"runtime"
	"strings"

	"github.com/michurin/minlog"

	"github.com/michurin/cnbot/app/aw"
)

func SetupLogging() {
	_, file, _, _ := runtime.Caller(0)
	pfx := strings.TrimSuffix(file, "app/log.go")
	opts := []minlog.Option(nil)
	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		opts = []minlog.Option{
			minlog.WithFields(
				minlog.Color(minlog.FieldLevel("", " ERR "), minlog.HiYellow, minlog.BgRed, minlog.Bold),
				minlog.Color(minlog.FieldCaller(pfx), minlog.HiBlue),
				minlog.Color(minlog.FieldErrorCaller(pfx), minlog.HiRed),
				minlog.Color(minlog.FieldNamed("comp"), minlog.HiGreen),
				minlog.Color(minlog.FieldNamed("bot"), minlog.Magenta),
				minlog.Color(minlog.FieldNamed("api"), minlog.Cyan),
				minlog.Color(minlog.FieldNamed("user"), minlog.HiGreen),
				minlog.Prefix(minlog.Color(minlog.FieldNamed("pid"), minlog.HiYellow), "PID:"),
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
