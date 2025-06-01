package sslog

import (
	"errors"
)

type Error struct {
	next error
	args []any
}

func (e Error) Error() string { return e.next.Error() }

func (e Error) Unwrap() error { return e.next }

func Wrap(err error, args []any) error {
	if err == nil {
		return err
	}
	if len(args) == 0 {
		return err
	}
	t := new(Error)
	if errors.As(err, t) {
		return err // already wrapped
	}
	return Error{
		next: err,
		args: args,
	}
}

func ArgsE(err error) []any {
	t := new(Error)
	if errors.As(err, t) { // allows err=nil
		return t.args
	}
	return nil
}
