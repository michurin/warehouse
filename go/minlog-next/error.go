package minlog

import (
	"context"
	"fmt"
)

type ctxError struct {
	err    error
	kv     map[string]any
	caller RecordCaller
}

func (e *ctxError) Error() string {
	return e.err.Error()
}

func (e *ctxError) Unwrap() error {
	return e.err
}

func Errorf(ctx context.Context, format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	kv := ctxKv(ctx)
	errCaller := ctxKvMergeError(kv, err)
	if errCaller.File == "" {
		errCaller = caller(2)
	}
	return &ctxError{
		err:    err,
		kv:     kv,
		caller: errCaller,
	}
}
