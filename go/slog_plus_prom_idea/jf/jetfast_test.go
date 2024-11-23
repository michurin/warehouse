package jf_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"x/jf"
)

type ctxKeyT int

const ctxKey ctxKeyT = iota

func timeReplacer(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey { // remove time just to be reproducible
		return slog.Attr{}
	}
	return a
}

func validate(x int) error {
	if x < 0 {
		return errors.New("negative number")
	}
	return nil
}

func handler(ctx context.Context, x int) error {
	ctx = context.WithValue(ctx, ctxKey, "handler specific")
	err := validate(x) // in real life it could be request to other service or to database
	if err != nil {
		return jf.E(ctx, err)
	}
	return nil // by the way, you can say just jf.E(ctx, err), E considers nil's in proper way
}

func Example_general() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: timeReplacer}))
	h = jf.New(h, "details", ctxKey) // wrap standard handler
	l := slog.New(h)

	ctx := context.Background()

	for _, x := range []int{1, -1} {
		err := handler(ctx, x)
		if err != nil {
			l.ErrorContext(jf.C(ctx, err), "Handler error")
			continue
		}
		l.InfoContext(ctx, "Handler OK") // we will see details=value from context
	}

	// output:
	// level=INFO msg="Handler OK"
	// level=ERROR msg="Handler error" details="handler has been called with argument -1"
}

func ExampleNew_simplest() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: timeReplacer}))
	h = jf.New(h, "details", ctxKey) // wrap standard handler
	l := slog.New(h)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey, "value")

	l.InfoContext(ctx, "Message") // we will see details=value from context

	// output:
	// level=INFO msg=Message details=value
}

func ExampleNew_withGroups() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: timeReplacer}))
	h = jf.New(h, "details", ctxKey) // wrap standard handler
	l := slog.New(h)

	l = l.With("pid", 100) // persistent attribute
	l = l.WithGroup("group")

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey, "value")

	l.InfoContext(ctx, "Message") // details from context will appear in group according to logger

	// output:
	// level=INFO msg=Message pid=100 group.details=value
}

func ExampleC_dealingWithErrors() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: timeReplacer}))
	h = jf.New(h, "details", ctxKey) // wrap standard handler
	l := slog.New(h)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey, "value")

	err := func(ctx context.Context) error {
		ctx = context.WithValue(ctx, ctxKey, "some details")
		err := errors.New("error message")
		return jf.E(ctx, err)
	}(ctx)
	if err != nil {
		l.ErrorContext(jf.C(ctx, err), "Error: "+err.Error())
	}

	// output:
	// level=ERROR msg="Error: error message" details="some details"
}

func ExampleC_itIsSafeToWrapAndUnwrapNilErrors() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: timeReplacer}))
	h = jf.New(h, "details", ctxKey) // wrap standard handler
	l := slog.New(h)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKey, "value")

	err := func(ctx context.Context) error {
		ctx = context.WithValue(ctx, ctxKey, "some details")
		err := error(nil)     // nil error (i.e. no error)
		return jf.E(ctx, err) // it is safe to wrap nil errors
	}(ctx)
	if err != nil {
		panic("it won't be fired")
	}
	l.InfoContext(jf.C(ctx, err), "Message") // it is safe to unwrap nil errors as well

	// output:
	// level=INFO msg=Message details=value
}

func ExampleE_doubleWrapping() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: timeReplacer}))
	h = jf.New(h, "details", ctxKey) // wrap standard handler
	l := slog.New(h)

	err := errors.New("error message")

	ctx := context.WithValue(context.Background(), ctxKey, "a")
	err = jf.E(ctx, err) // wrap number one (initial)

	ctx = context.WithValue(context.Background(), ctxKey, "b")
	err = jf.E(ctx, err) // wrap number two

	l.InfoContext(jf.C(ctx, err), "Message") // we will see context from initial wrap (details=a)

	// output:
	// level=INFO msg=Message details=a
}

func ExampleWrap_panicSafeWrapping() {
	h := slog.Handler(slog.NewTextHandler(os.Stdout, nil))

	_, err := jf.Wrap(h, "wrong", "number of", "arguments")
	if err != nil {
		fmt.Println("Error:", err)
	}

	_, err = jf.Wrap(h, true, true)
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()
	h = jf.New(h, true)

	// output:
	// Error: odd number of key/source pairs: 3
	// Error: key must be a string: bool: true
	// Panic: odd number of key/source pairs: 1
}

func ExampleE_asAndIsWorkAsExpected() {
	ctx := context.Background()

	errx := errors.New("error message")

	err := jf.E(ctx, errx)

	if errors.Is(err, errx) {
		fmt.Println("err is errx")
	}

	// output:
	// err is errx
}
