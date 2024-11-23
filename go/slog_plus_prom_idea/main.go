package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"x/jf"
)

type keyT int

const requestKey keyT = 0

func x(ctx context.Context) error {
	ctx = context.WithValue(ctx, requestKey, "request from function")
	err := errors.New("just error")
	err = jf.E(ctx, err)
	return err
}

func main() {
	slog.SetDefault(slog.New(jf.New(slog.NewJSONHandler(os.Stdout, nil), "request", requestKey)))

	ctx := context.Background()

	err := x(ctx)
	if err != nil {
		slog.ErrorContext(jf.C(ctx, err), "Error")
	}

	l := slog.New(jf.New(slog.NewJSONHandler(os.Stdout, nil), "request", requestKey)).WithGroup("group")
	l.InfoContext(jf.C(ctx, err), "This is a message")
}
