package argerr

import (
	"errors"
)

type xerr struct {
	next error
	args []any
}

func (e xerr) Error() string { return e.next.Error() }

func (e xerr) Unwrap() error { return e.next }

func Wrap(err error, args ...any) error {
	if err == nil {
		return err
	}
	if len(args) == 0 {
		return err
	}
	t := new(xerr)
	if errors.As(err, t) {
		return err // already wrapped
	}
	return xerr{
		next: err,
		args: args,
	}
}

func Args(err error) []any {
	t := new(xerr)
	if errors.As(err, t) { // allows err=nil
		return t.args
	}
	return nil
}
