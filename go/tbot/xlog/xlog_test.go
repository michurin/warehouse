package xlog_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/michurin/warehouse/go/tbot/xlog"
)

type customError struct{}

func (customError) Error() string { return "custom" }

func ExampleLog() {
	xlog.Fields = []string{"a", "b", "c", "d"}
	ctx := xlog.Ctx(context.Background(), "a", 0)
	xlog.Log(ctx, "Just message")
	err := error(customError{})
	xlog.Log(ctx, "Naked error", err)
	xlog.Log(xlog.Ctx(ctx, "a", 7, "b", "9"), "Naked error with tweaked ctx", err)
	err = xlog.Errorf(xlog.Ctx(ctx, "a", 1), "e1: %w", err) // [A]
	err = fmt.Errorf("f1: %w", err)
	err = xlog.Errorf(xlog.Ctx(ctx, "a", 2, "b", 2), "e2: %w", err) // [B]
	err = fmt.Errorf("f2: %w", err)
	ctx = xlog.Ctx(ctx, "a", 4, "b", 4, "c", 4, "d", 4)
	xlog.Log(xlog.Ctx(ctx, "c", 3), "Message", err) // "a" comes from [A], "b" comes from "B"
	fmt.Println(errors.Is(err, customError{}))
	// Output:
	// [info] 0 Just message
	// [error] 0 Naked error custom
	// [error] 7 9 Naked error with tweaked ctx custom
	// [error] 1 2 3 4 Message f2: e2: f1: e1: custom
	// true
}
