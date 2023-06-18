package minlog_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/michurin/minlog"
)

type customError struct{}

func (customError) Error() string { return "custom" }

func Example_almostAllFeaturesOverview() {
	_, file, _, _ := runtime.Caller(0)               //nolint:dogsled
	pfx := path.Dir(file) + string(os.PathSeparator) // prefix of file, including last separator

	log := minlog.New(
		minlog.WithStdLogger(log.New(os.Stdout, "", 0)),
		minlog.WithPersistFields("p", "pst"),
		minlog.WithFields(
			minlog.FieldLevel("[INFO]", "[ERRR]"),
			minlog.FieldCaller(pfx),
			minlog.FieldErrorCaller(pfx),
			minlog.FieldNamed("p"),
			minlog.FieldNamed("a"),
			minlog.FieldNamed("b"),
			minlog.FieldNamed("c"),
			minlog.FieldNamed("d"),
			minlog.FieldFallbackKV("a", "b", "c", "d", "p"),
			minlog.FieldMessage(),
		),
	).Log

	ctx := minlog.Ctx(context.Background(), "a", 0)
	log(ctx, "Just message")
	log(ctx, "Test formatting", []byte("valid bytes"))
	log(ctx, "Test formatting", append([]byte("invalid bytes "), 255))
	log(ctx, "Test formatting numbers", 255, 3.14)
	err := error(customError{})
	log(ctx, "Naked error", err)
	log(minlog.Ctx(ctx, "a", 7, "b", "9"), "Naked error with tweaked ctx", err)
	err = minlog.Errorf(minlog.Ctx(context.Background(), "a", 1), "e1: %w", err) // [A]
	err = fmt.Errorf("f1: %w", err)
	err = minlog.Errorf(minlog.Ctx(context.Background(), "a", 2, "b", 2), "e2: %w", err) // [B]
	err = fmt.Errorf("f2: %w", err)
	ctx = minlog.Ctx(ctx, "a", 4, "b", 4, "c", 4, "d", 4)
	log(minlog.Ctx(ctx, "c", 3), "Message", err) // "a" comes from [A], "b" comes from "B"
	log(context.Background(), errors.Is(err, customError{}))

	ctx1 := minlog.Ctx(context.Background(), "a", 1, "b", 1)
	ctx2 := minlog.Ctx(context.Background(), "b", 2, "c", 2)
	ctx = minlog.ApplyPatch(ctx2, minlog.TakePatch(ctx1)) // 2<-1: b=1 replace b=2
	log(ctx, "Show cloning")

	log(minlog.Ctx(context.Background(), "a", 1, "f", 7), "Unknown field")

	// Output:
	// [INFO] log_test.go:42 pst 0 Just message
	// [INFO] log_test.go:43 pst 0 Test formatting valid bytes
	// [INFO] log_test.go:44 pst 0 Test formatting "invalid bytes \xff"
	// [INFO] log_test.go:45 pst 0 Test formatting numbers 255 3.14
	// [ERRR] log_test.go:47 pst 0 Naked error custom
	// [ERRR] log_test.go:48 pst 7 9 Naked error with tweaked ctx custom
	// [ERRR] log_test.go:54 log_test.go:49 pst 1 2 3 4 Message f2: e2: f1: e1: custom
	// [INFO] log_test.go:55 pst true
	// [INFO] log_test.go:60 pst 1 1 2 Show cloning
	// [INFO] log_test.go:62 pst 1 f=7 Unknown field
}

func Example_structuredLogging() {
	log := minlog.New(
		minlog.WithStdLogger(log.New(os.Stdout, "", 0)),
		minlog.WithFields(minlog.FieldJSON()),
	)
	ctx := context.Background()
	log.Log(ctx, "OK")
	ctx = minlog.Ctx(ctx, "k", "val")
	log.Log(ctx, "Error:", errors.New("it's error"))
	// Output:
	// {"caller":"log_test.go:83","context":{},"level":"info","message":"OK"}
	// {"caller":"log_test.go:85","context":{"k":"val"},"level":"error","message":"Error: it's error"}
}

func TestFieldJSON_error(t *testing.T) {
	f := minlog.FieldJSON()
	s := f(minlog.Record{Context: map[string]any{"x": immarshalable{}}})
	if s != `{"marshaller_error":"FieldJSON error: json: error calling MarshalJSON for type minlog_test.immarshalable: err"}` {
		t.Log(s)
		t.Fail()
	}
}

func TestNew_withoutArgs(_ *testing.T) { // assert nothing, we just ensure there are no panics with defaults
	log := minlog.New()
	log.Log(context.Background(), "")
}

type immarshalable struct{}

func (immarshalable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("err")
}
