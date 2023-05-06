package xlog_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/michurin/warehouse/go/tbot/xlog"
)

// TODO cannot be run in parallel

type customError struct{}

func (customError) Error() string { return "custom" }

func ExampleLog_message() {
	xlog.Fields = []xlog.Field{
		xlog.StdFieldLevel,
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
		{Name: "d"},
		xlog.StdFieldMessage,
	}
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

func ExampleLog_caller() {
	xlog.Fields = []xlog.Field{
		xlog.StdFieldLevel,
		xlog.StdFieldCaller,
		xlog.StdFieldOCaller,
		xlog.StdFieldMessage,
	}
	ctx := context.Background()
	err := xlog.Errorf(ctx, "Just message")    // we will see this line in logs
	err = xlog.Errorf(ctx, "Wrapped: %w", err) // not this; but log message will be wrapped correctly
	xlog.Log(ctx, err)
	// Output:
	// [error] xlog/xlog_test.go:56 xlog/xlog_test.go:54 Wrapped: Just message
}

func ExampleLog_formatting() {
	xlog.Fields = []xlog.Field{
		xlog.StdFieldLevel,
		xlog.StdFieldMessage,
	}
	ctx := context.Background()
	xlog.Log(ctx, []byte(nil))       // nil
	xlog.Log(ctx, []byte{})          // len=0
	xlog.Log(ctx, []byte("ok"))      // bytes
	xlog.Log(ctx, []byte{255})       // wrong char
	xlog.Log(ctx, error(nil))        // nil error
	xlog.Log(ctx, errors.New("err")) // true error
	// Output:
	// [info] ""
	// [info] ""
	// [info] ok
	// [info] "\xff"
	// [info] <nil>
	// [error] err
}

func ExampleLog_cloneContext() {
	xlog.Fields = []xlog.Field{
		xlog.StdFieldLevel,
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
		xlog.StdFieldMessage,
	}
	ctx1 := xlog.Ctx(context.Background(), "a", "A1", "b", "B1")
	xlog.Log(ctx1, "Original context")
	ctx2 := xlog.Ctx(context.Background(), "a", "A2", "c", "C2")
	xlog.Log(ctx2, "Second context")
	ctx3 := xlog.CloneCtx(ctx2, ctx1) // copy ctx1 -> ctx2
	xlog.Log(ctx3, "Meted together context")
	// Output:
	// [info] A1 B1 Original context
	// [info] A2 C2 Second context
	// [info] A1 B1 C2 Meted together context
}
