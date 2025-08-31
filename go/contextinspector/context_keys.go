package contextinspector

import (
	"context"
	"reflect"
	"unsafe"
)

type TypeInfo struct {
	Next string
	Key  string
}

var /* const */ stdCtxTypes = map[string]TypeInfo{
	"context.backgroundCtx":    {},
	"context.cancelCtx":        {Next: "Context"},
	"context.emptyCtx":         {},
	"context.stopCtx":          {Next: "Context"}, // go1.21
	"context.timerCtx":         {Next: "Context"},
	"context.todoCtx":          {},
	"context.valueCtx":         {Next: "Context", Key: "key"},
	"context.withoutCancelCtx": {Next: "Context"}, // go1.21
}

func collectKeys(ctx any, acc map[any]int, typeInfos map[string]TypeInfo) {
	v := reflect.ValueOf(ctx)
	t := reflect.TypeOf(ctx)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = v.Type()
	}
	typeName := t.PkgPath() + "." + t.Name()
	info, ok := typeInfos[typeName]
	if !ok {
		panic("Unknown type " + typeName)
	}
	if info.Key != "" {
		k := v.FieldByName(info.Key)
		if k.CanAddr() {
			x := reflect.NewAt(k.Type(), unsafe.Pointer(k.UnsafeAddr())).Elem().Interface() // Black magic: getting value of unexported field
			acc[x] = acc[x] + 1
		} else {
			println("TODO Can not addr") // TODO
		}
	}
	if info.Next != "" {
		if v.CanAddr() { // TODO CanInterface? It can panic
			collectKeys(v.FieldByName(info.Next).Interface(), acc, typeInfos)
		} else {
			println("TODO Can not addr") // TODO
		}
	}
}

func CtxKeysCounters(ctx context.Context) map[any]int {
	acc := map[any]int{}
	collectKeys(ctx, acc, stdCtxTypes)
	return acc
}

func CtxKeys(ctx context.Context) []any {
	m := CtxKeysCounters(ctx)
	r := []any(nil)
	for k := range m {
		r = append(r, k)
	}
	return r
}

func CtxKeysCountersWithCustom(ctx context.Context, tp map[string]TypeInfo) map[any]int {
	acc := map[any]int{}
	t := map[string]TypeInfo{}
	for k, v := range stdCtxTypes {
		t[k] = v
	}
	for k, v := range tp {
		t[k] = v
	}
	collectKeys(ctx, acc, t)
	return acc
}

func CtxKeysWithCustom(ctx context.Context, tp map[string]TypeInfo) []any {
	m := CtxKeysCountersWithCustom(ctx, tp)
	r := []any(nil)
	for k := range m {
		r = append(r, k)
	}
	return r
}
