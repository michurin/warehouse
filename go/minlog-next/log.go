package minlog

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
)

type RecordCaller struct {
	File string
	Line int
}

type Record struct {
	Message     string
	IsError     bool
	Context     map[string]any
	Caller      RecordCaller
	ErrorCaller RecordCaller
}

type FieldFunc func(Record) string

type Logger struct {
	printer       func(string)
	fields        []FieldFunc
	argFormatter  func(any) string
	persistFields map[string]string
}

type Option func(*Logger) // dedicated type just to ease code navigation

func New(option ...Option) *Logger {
	l := new(Logger)
	for _, op := range option {
		op(l)
	}
	if l.printer == nil {
		WithStdLogger(log.New(os.Stderr, "", 0))(l)
	}
	if len(l.fields) == 0 {
		WithFields(
			FieldLevel("[INFO]", "[ERRR]"),
			FieldCaller(""),
			FieldErrorCaller(""),
			FieldFallbackKV(),
			FieldMessage(),
		)(l)
	}
	if l.argFormatter == nil {
		WithArgFormatter(formatArg)(l)
	}
	if l.persistFields == nil {
		WithPersistFields()(l)
	}
	return l
}

func (l *Logger) Log(ctx context.Context, a ...any) {
	kv := map[string]any{}
	for k, v := range l.persistFields {
		kv[k] = v
	}
	ctxKvMerge(kv, ctx)

	isErr := false
	errCaller := RecordCaller{}
	msg := make([]string, len(a))
	for i, x := range a {
		if err, ok := x.(error); ok {
			isErr = true
			errCaller = ctxKvMergeError(kv, err)
		}
		msg[i] = l.argFormatter(x)
	}
	message := strings.Join(msg, " ")

	rec := Record{
		Message:     message,
		IsError:     isErr,
		Context:     kv,
		Caller:      caller(2),
		ErrorCaller: errCaller,
	}
	fs := []string(nil)
	for _, f := range l.fields {
		p := f(rec)
		if p != "" {
			fs = append(fs, p)
		}
	}

	l.printer(strings.Join(fs, " ") + "\n")
}

var ctxKey = struct{}{} //nolint:gochecknoglobals

func Ctx(ctx context.Context, kv ...any) context.Context {
	nkv := ctxKv(ctx)
	for i := 0; i < len(kv)-1; i += 2 {
		if k, ok := kv[i].(string); ok {
			nkv[k] = kv[i+1]
		}
	}
	return context.WithValue(ctx, ctxKey, nkv)
}

func ctxKv(ctx context.Context) map[string]any {
	kv := map[string]any{}
	ctxKvMerge(kv, ctx)
	return kv
}

func ctxKvMerge(kv map[string]any, ctx context.Context) { //nolint:revive
	if x, ok := ctx.Value(ctxKey).(map[string]any); ok {
		for k, v := range x {
			kv[k] = v
		}
	}
}

func ctxKvMergeError(kv map[string]any, err error) RecordCaller {
	t := &ctxError{}
	if errors.As(err, &t) {
		for k, v := range t.kv {
			kv[k] = v
		}
		return t.caller
	}
	return RecordCaller{}
}
