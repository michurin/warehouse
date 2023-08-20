package xslog_test

import (
	"errors"
	"testing"

	"slogplayground/xslog"
)

func TestErrWrap(t *testing.T) {
	specificErr := errors.New("x")
	err := xslog.Errorf("err: %w", specificErr)
	t.Log(errors.Is(err, specificErr))
}
