package sslog

import (
	"context"
	"errors"
)

type Error struct {
	next error
	args []any
}

func (e Error) Error() string { return e.next.Error() }

func (e Error) Unwrap() error { return e.next }

func Wrap(ctx context.Context, err error) error {
	if err == nil {
		return err
	}
	if ctx == nil {
		return err
	}
	t := new(Error)
	if errors.As(err, t) {
		return err // already wrapped
	}
	return Error{
		next: err,
		args: Args(ctx),
	}
}

func ArgsE(err error) []any {
	t := new(Error)
	if errors.As(err, t) { // allows err=nil
		return t.args
	}
	return nil
}
