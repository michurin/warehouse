package argctx_test

import (
	"context"
	"fmt"

	"github.com/michurin/minlog/argctx"
)

func Example() {
	ctx := context.Background()
	ctx = argctx.With(ctx, 1, "a", "b")
	fmt.Println(argctx.Args(ctx))
	// output:
	// [1 a b]
}

func ExampleWith() {
	ctx := context.Background()
	ctx = argctx.With(ctx, 1)
	ctx = argctx.With(ctx, "a", "b") // With can be called several times
	ctx = argctx.With(ctx)           // no args are legal
	fmt.Println(argctx.Args(ctx))
	// output:
	// [1 a b]
}

func ExampleWith_group() {
	ctx := context.Background()
	ctx = argctx.With(ctx, 1)
	ctx = argctx.With(ctx, func(x []any) any {
		return fmt.Sprintf("group:%v", x) // in real life it can be something like: return slog.Group("group", x...)
	})
	ctx = argctx.With(ctx, 2)
	ctx = argctx.With(ctx, 3)
	fmt.Println(argctx.Args(ctx))
	// output:
	// [1 group:[2 3]]
}

func ExampleArgs() {
	ctx := context.Background()
	fmt.Printf("%#v\n", argctx.Args(ctx)) // it is safe to call Args without With
	ctx = argctx.With(ctx, 1, "a", "b")
	fmt.Println(argctx.Args(ctx))
	// output:
	// []interface {}(nil)
	// [1 a b]
}
