package sslog_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"sslog"
)

func Test(t *testing.T) {
	ctx := t.Context()
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx = sslog.With(ctx, "K", "V")
	ctx = sslog.With(ctx, slog.Group("g", "kk", "vv"))
	log.Info("OK", sslog.Args(ctx)...)

	err := errors.New("message")
	err = sslog.Wrap(sslog.With(ctx, "E", "X"), err)
	log.Error("ERR", sslog.ArgsE(err)...)

	// TODO ASSERTS
}
