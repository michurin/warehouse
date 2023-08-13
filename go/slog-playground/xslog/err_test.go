package xslog_test

import (
	"context"
	"errors"
	"testing"

	"slogplayground/xslog"
)

func TestErrWrap(t *testing.T) {
	specificErr := errors.New("x")
	err := xslog.Errorf(context.Background(), "err: %w", specificErr)
	t.Log(errors.Is(err, specificErr))
}
