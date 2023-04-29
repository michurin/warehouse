package xlog

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// TODO it has to be DefaultLogger and methods

var (
	Fields     = []string(nil) // TODO it has to be []struct{name string; formatter: func(string) string}?
	LabelInfo  = "[info]"      // TODO move labels to fields, to be able to order them? caller etc to fields too?
	LabelError = "[error]"
)

var ctxKey = struct{}{}

type ctxError struct {
	err error
	kv  map[string]any
}

func (e *ctxError) Error() string {
	return e.err.Error()
}

func (e *ctxError) Unwrap() error {
	return e.err
}

func Errorf(ctx context.Context, format string, a ...any) error {
	// TODO original caller
	err := fmt.Errorf(format, a...)
	kv := ctxKv(ctx)
	ctxKvOverride(kv, err)
	return &ctxError{
		err: err,
		kv:  kv,
	}
}

func Ctx(ctx context.Context, kv ...any) context.Context {
	nkv := ctxKv(ctx)
	for i := 0; i < len(kv)-1; i += 2 {
		if k, ok := kv[i].(string); ok {
			nkv[k] = kv[i+1]
		}
	}
	return context.WithValue(ctx, ctxKey, nkv)
}

func Log(ctx context.Context, a ...any) {
	// TODO caller
	fkv := ctxKv(ctx)

	errorLevel := false
	for _, x := range a {
		if err, ok := x.(error); ok {
			errorLevel = true
			ctxKvOverride(fkv, err)
		}
	}
	label := LabelInfo
	if errorLevel {
		label = LabelError
	}

	fs := []string(nil)
	for _, f := range Fields {
		if x, ok := fkv[f]; ok {
			fs = append(fs, fmt.Sprintf("%+v", x))
		}
	}

	msg := []string(nil)
	for _, m := range a {
		msg = append(msg, formatArg(m))
	}

	// TODO ------v--v-- remove this spaces if %s substitutions are empty strings
	fmt.Printf("%s %s %s\n", label, strings.Join(fs, " "), strings.Join(msg, " "))
}

func formatArg(x any) string { // TODO has to be method, has to be tunable
	s := ""
	switch t := x.(type) {
	case string:
		s = t
	case []byte:
		s = fmt.Sprintf("%q", string(t)) // TODO check correct UTF
	default:
		s = fmt.Sprintf("%v", t)
	}
	if len(s) > 1000 { // TODO rune!
		s = s[:400]
	}
	return s
}

func ctxKv(ctx context.Context) map[string]any {
	a := map[string]any{}
	if kv, ok := ctx.Value(ctxKey).(map[string]any); ok {
		for k, v := range kv {
			a[k] = v
		}
	}
	return a
}

func ctxKvOverride(kv map[string]any, err error) {
	t := &ctxError{}
	if errors.As(err, &t) {
		for k, v := range t.kv {
			kv[k] = v
		}
	}
}
