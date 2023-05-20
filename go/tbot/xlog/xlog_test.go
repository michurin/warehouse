package xlog_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/michurin/warehouse/go/tbot/xlog"
)

type customError struct{}

func (customError) Error() string { return "custom" }

func ExampleNew() {
	_, file, _, _ := runtime.Caller(0)
	pfx := file[:strings.LastIndex(file, "/")+1]

	log := xlog.New(
		xlog.WithStdLogger(log.New(os.Stdout, "", 0)),
		xlog.WithPersistFields("p", "pst"),
		xlog.WithFields(
			xlog.FieldLevel("[INFO]", "[ERRR]"),
			xlog.FieldCaller(pfx),
			xlog.FieldErrorCaller(pfx),
			xlog.FieldNamed("p"),
			xlog.FieldNamed("a"),
			xlog.FieldNamed("b"),
			xlog.FieldNamed("c"),
			xlog.FieldNamed("d"),
			xlog.FieldFallbackKV("a", "b", "c", "d", "p"),
			xlog.FieldMessage(),
		),
	).Log

	ctx := xlog.Ctx(context.Background(), "a", 0)
	log(ctx, "Just message")
	log(ctx, "Test formatting", []byte("valid bytes"))
	log(ctx, "Test formatting", append([]byte("invalid bytes "), 255))
	log(ctx, "Test formatting numbers", 255, 3.14)
	err := error(customError{})
	log(ctx, "Naked error", err)
	log(xlog.Ctx(ctx, "a", 7, "b", "9"), "Naked error with tweaked ctx", err)
	err = xlog.Errorf(xlog.Ctx(context.Background(), "a", 1), "e1: %w", err) // [A]
	err = fmt.Errorf("f1: %w", err)
	err = xlog.Errorf(xlog.Ctx(context.Background(), "a", 2, "b", 2), "e2: %w", err) // [B]
	err = fmt.Errorf("f2: %w", err)
	ctx = xlog.Ctx(ctx, "a", 4, "b", 4, "c", 4, "d", 4)
	log(xlog.Ctx(ctx, "c", 3), "Message", err) // "a" comes from [A], "b" comes from "B"
	log(context.Background(), errors.Is(err, customError{}))

	ctx1 := xlog.Ctx(context.Background(), "a", 1, "b", 1)
	ctx2 := xlog.Ctx(context.Background(), "b", 2, "c", 2)
	ctx = xlog.CloneCtx(ctx2, ctx1) // 2<-1: b=1 replace b=2
	log(ctx, "Show cloning")

	log(xlog.Ctx(context.Background(), "a", 1, "f", 7), "Unknown field")

	// Output:
	// [INFO] xlog_test.go:41 pst 0 Just message
	// [INFO] xlog_test.go:42 pst 0 Test formatting valid bytes
	// [INFO] xlog_test.go:43 pst 0 Test formatting "invalid bytes \xff"
	// [INFO] xlog_test.go:44 pst 0 Test formatting numbers 255 3.14
	// [ERRR] xlog_test.go:46 pst 0 Naked error custom
	// [ERRR] xlog_test.go:47 pst 7 9 Naked error with tweaked ctx custom
	// [ERRR] xlog_test.go:53 xlog_test.go:48 pst 1 2 3 4 Message f2: e2: f1: e1: custom
	// [INFO] xlog_test.go:54 pst true
	// [INFO] xlog_test.go:59 pst 1 1 2 Show cloning
	// [INFO] xlog_test.go:61 pst 1 f=7 Unknown field
}
