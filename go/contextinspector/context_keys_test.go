package contextinspector_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	contextinspector "github.com/michurin/warehouse/go/contextinspector"

	"github.com/stretchr/testify/assert"
)

func ExampleCtxKeys() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "K", "VAL")
	fmt.Println(contextinspector.CtxKeys(ctx))
	fmt.Println(ctx.Value(contextinspector.CtxKeys(ctx)[0]))
	// Output:
	// [K]
	// VAL
}

func TestCtxKeys(t *testing.T) {
	t.Run("background", func(t *testing.T) {
		ctx := context.Background()
		r := contextinspector.CtxKeys(ctx)
		if r != nil {
			t.Fail()
		}
	})
	t.Run("todo", func(t *testing.T) {
		ctx := context.TODO()
		r := contextinspector.CtxKeys(ctx)
		assert.Nil(t, r)
	})
	t.Run("with_cancel", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		r := contextinspector.CtxKeys(ctx)
		assert.Nil(t, r)
	})
	t.Run("with_timeout", func(t *testing.T) { // WithDeadline too
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		r := contextinspector.CtxKeys(ctx)
		assert.Nil(t, r)
	})
	t.Run("with_value", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, 1, 2)
		r := contextinspector.CtxKeys(ctx)
		assert.Equal(t, []interface{}{1}, r)
	})
}

type customString string

func TestCtxKeysCounters(t *testing.T) {
	t.Run("with_value", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "A", "A1")
		ctx = context.WithValue(ctx, "A", "A2")
		ctx = context.WithValue(ctx, "B", "B1")
		r := contextinspector.CtxKeysCounters(ctx)
		assert.Equal(t, map[interface{}]int(map[interface{}]int{"A": 2, "B": 1}), r)
	})
	t.Run("with_value_custom_type", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "A", "A1")
		ctx = context.WithValue(ctx, "A", "A2")
		ctx = context.WithValue(ctx, customString("A"), "A3")
		ctx = context.WithValue(ctx, "B", "B1")
		r := contextinspector.CtxKeysCounters(ctx)
		assert.Equal(t, map[interface{}]int(map[interface{}]int{
			"A":               2,
			customString("A"): 1,
			"B":               1}), r)
	})
}

type customCtxAlias bool

func (_ customCtxAlias) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (_ customCtxAlias) Done() <-chan struct{}                   { return nil }
func (_ customCtxAlias) Err() error                              { return nil }
func (_ customCtxAlias) Value(key interface{}) interface{}       { return nil }

type customCtxChain struct {
	privateNext context.Context
	customKey   string
}

func (_ customCtxChain) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (_ customCtxChain) Done() <-chan struct{}                   { return nil }
func (_ customCtxChain) Err() error                              { return nil }
func (_ customCtxChain) Value(key interface{}) interface{}       { return nil }

func TestCtxKeysPanic(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx := customCtxAlias(false)
			_ = contextinspector.CtxKeys(ctx)
		})
	})
}

func TestCtxKeysWithCustom_does_not_work_yet(t *testing.T) {
	ctX := customCtxChain{
		privateNext: customCtxAlias(false),
		customKey:   "one",
	}
	ctx := context.WithValue(ctX, 1, 2)
	_ = contextinspector.CtxKeysWithCustom(ctx, map[string]contextinspector.TypeInfo{
		"github.com/michurin/warehouse/go/contextinspector_test.customCtxChain": {
			Next: "privateNext",
			Key:  "customKey",
		}})
}
