package minlog

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"unicode/utf8"
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
			errCaller = ctxKvOverride(kv, err)
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

func caller(level int) RecordCaller {
	_, file, no, ok := runtime.Caller(level)
	if !ok {
		return RecordCaller{
			File: "nofile",
			Line: 0,
		}
	}
	return RecordCaller{
		File: file,
		Line: no,
	}
}

// ---- TODO split (?)

func Ctx(ctx context.Context, kv ...any) context.Context {
	nkv := ctxKv(ctx)
	for i := 0; i < len(kv)-1; i += 2 {
		if k, ok := kv[i].(string); ok {
			nkv[k] = kv[i+1]
		}
	}
	return context.WithValue(ctx, ctxKey, nkv)
}

type LogPatch struct {
	kv map[string]any
}

// Patch woks with ApplyPatch like this to copy logging context to another go context
func Patch(ctx context.Context) LogPatch {
	return LogPatch{kv: ctxKv(ctx)}
}

func ApplyPatch(ctx context.Context, patch LogPatch) context.Context {
	kv := ctxKv(ctx)
	for k, v := range patch.kv {
		kv[k] = v
	}
	return context.WithValue(ctx, ctxKey, kv)
}

// ---- TODO split

func FieldMessage() FieldFunc {
	return func(r Record) string {
		return r.Message
	}
}

func FieldCaller(pfx string) FieldFunc {
	return func(r Record) string {
		return fmt.Sprintf("%s:%d", strings.TrimPrefix(r.Caller.File, pfx), r.Caller.Line)
	}
}

func FieldErrorCaller(pfx string) FieldFunc {
	return func(r Record) string {
		if r.IsError {
			if r.ErrorCaller.File == "" {
				return ""
			}
			return fmt.Sprintf("%s:%d", strings.TrimPrefix(r.ErrorCaller.File, pfx), r.ErrorCaller.Line)
		}
		return ""
	}
}

func FieldLevel(info, errr string) FieldFunc {
	return func(r Record) string {
		if r.IsError {
			return errr
		}
		return info
	}
}

func FieldFallbackKV(exclude ...string) FieldFunc {
	exc := map[string]struct{}{}
	for _, v := range exclude {
		exc[v] = struct{}{}
	}
	return func(r Record) string {
		fs := []string(nil)
		for k := range r.Context {
			if _, ok := exc[k]; ok {
				continue
			}
			fs = append(fs, k)
		}
		if fs == nil {
			return ""
		}
		sort.Strings(fs)
		pts := make([]string, len(fs))
		for i, k := range fs {
			pts[i] = fmt.Sprintf("%s=%v", k, r.Context[k])
		}
		return strings.Join(pts, " ")
	}
}

func FieldNamed(fieldName string) FieldFunc {
	return func(r Record) string {
		if x, ok := r.Context[fieldName]; ok {
			return fmt.Sprintf("%v", x)
		}
		return ""
	}
}

// ---- TODO split

func WithStdLogger(printer interface{ Print(v ...any) }) Option {
	return func(l *Logger) {
		l.printer = func(x string) {
			printer.Print(x)
		}
	}
}

func WithFields(fields ...FieldFunc) Option {
	return func(l *Logger) {
		l.fields = fields
	}
}

func WithArgFormatter(formatter func(any) string) Option {
	return func(l *Logger) {
		l.argFormatter = formatter
	}
}

func WithPersistFields(kv ...string) Option {
	p := map[string]string{}
	for i := 0; i < len(kv)-1; i++ {
		p[kv[i]] = kv[i+1]
	}
	return func(l *Logger) {
		l.persistFields = p
	}
}

// ---------------------------------

var ctxKey = struct{}{}

type ctxError struct {
	err    error
	kv     map[string]any
	caller RecordCaller
}

func (e *ctxError) Error() string {
	return e.err.Error()
}

func (e *ctxError) Unwrap() error {
	return e.err
}

func Errorf(ctx context.Context, format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	kv := ctxKv(ctx)
	ec := ctxKvOverride(kv, err)
	if ec.File == "" {
		ec = caller(2)
	}
	return &ctxError{
		err:    err,
		kv:     kv,
		caller: ec,
	}
}

func formatArg(x any) string {
	s := ""
	switch t := x.(type) {
	case string:
		s = t
		if s == "" {
			s = "[empty string]"
		}
	case []byte:
		if len(t) > 0 && utf8.Valid(t) {
			s = string(t)
		} else {
			s = fmt.Sprintf("%q", t)
		}
	case bool:
		if t {
			return "true"
		}
		return "false"
	case error:
		return t.Error()
	default:
		s = fmt.Sprintf("%v", t)
	}
	if len([]rune(s)) > 1000 {
		s = s[:400]
	}
	return s
}

func ctxKv(ctx context.Context) map[string]any {
	kv := map[string]any{}
	ctxKvMerge(kv, ctx)
	return kv
}

func ctxKvMerge(kv map[string]any, ctx context.Context) {
	if x, ok := ctx.Value(ctxKey).(map[string]any); ok {
		for k, v := range x {
			kv[k] = v
		}
	}
}

func ctxKvOverride(kv map[string]any, err error) RecordCaller { // TODO rename?
	t := &ctxError{}
	if errors.As(err, &t) {
		for k, v := range t.kv {
			kv[k] = v
		}
		return t.caller
	}
	return RecordCaller{}
}
