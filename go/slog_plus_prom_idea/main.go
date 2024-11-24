package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"x/jf"

	"github.com/mdobak/go-xerrors"
)

type keyT int

const requestKey keyT = 0

func f(_ context.Context, e error) []any {
	r := []any{"error", e.Error()}
	tr := xerrors.StackTrace(e)
	if tr != nil {
		f := tr.Frames()
		s := make([]map[string]any, len(f))
		for i, v := range f {
			s[i] = map[string]any{
				"source": filepath.Join(
					filepath.Base(filepath.Dir(v.File)),
					filepath.Base(v.File),
				),
				"line": v.Line,
				"func": filepath.Base(v.Function),
			}
		}
		r = append(r, "trace", s)
	}
	return r
}

func x(ctx context.Context) error {
	ctx = context.WithValue(ctx, requestKey, "request from function")
	err := errors.New("just error")
	err = xerrors.WithStackTrace(err, 0)
	err = jf.E(ctx, err)
	return err
}

func main() {
	slog.SetDefault(slog.New(jf.New(
		slog.NewJSONHandler(os.Stdout, nil),
		jf.F("request", requestKey),
		jf.ErrF(f))))

	ctx := context.Background()
	ctx = context.WithValue(ctx, requestKey, "OK?")

	err := x(ctx)
	if err != nil {
		slog.ErrorContext(jf.C(ctx, err), "Error")
	}

	l := slog.New(jf.New(slog.NewJSONHandler(os.Stdout, nil), jf.F("request", requestKey))).WithGroup("group")
	l.InfoContext(jf.C(ctx, err), "This is a message")
}
