package xlog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

// TODO it has to be DefaultLogger and methods

var pathReplacer = func() *strings.Replacer {
	// WARNING: this naive approach wont work in dedicated package
	_, fileName, _, _ := runtime.Caller(0)
	sep := string(os.PathSeparator)
	p := strings.Split(fileName, sep)
	s := strings.Join(p[:len(p)-2], sep) + sep
	return strings.NewReplacer(s, "")
}()

func relativeCaller(level int) string {
	_, file, no, ok := runtime.Caller(level)
	if ok {
		return fmt.Sprintf("%s:%d", pathReplacer.Replace(file), no)
	}
	return "nocaller"
}

type Field struct {
	Name string
	Proc func(any) string
}

func ProcFuncCaller(any) string {
	return relativeCaller(3)
}

var StdFieldTime = Field{
	Name: "log_time",
	Proc: func(any) string {
		return time.Now().Format("2006-01-02 15:04:05.000")
	},
}

var StdFieldLevel = Field{
	Name: "log_level",
	Proc: func(x any) string {
		if x.(int) == LevelError {
			return "[error]"
		}
		return "[info]"
	},
}

var StdFieldCaller = Field{
	Name: "log_caller",
	Proc: ProcFuncCaller,
}

var StdFieldOCaller = Field{
	Name: "log_ocaller",
	Proc: func(x any) string {
		return x.(string)
	},
}

var StdFieldMessage = Field{
	Name: "log_message",
	Proc: func(x any) string {
		return x.(string)
	},
}

var Fields = []Field{
	StdFieldTime,
	StdFieldLevel,
	StdFieldCaller,
	StdFieldOCaller,
	StdFieldMessage,
}

const (
	LevelInfo = iota
	LevelError
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
	kv["log_ocaller"] = relativeCaller(2)
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
	fkv := ctxKv(ctx)

	errorLevel := LevelInfo
	msg := make([]string, len(a))
	for i, x := range a {
		if err, ok := x.(error); ok {
			errorLevel = LevelError
			ctxKvOverride(fkv, err)
		}
		msg[i] = formatArg(x)
	}
	fkv["log_level"] = errorLevel
	fkv["log_message"] = strings.Join(msg, " ")

	fkv["log_time"] = nil   // value doesn't matter, looks slightly hackish
	fkv["log_caller"] = nil // value doesn't matter, looks slightly hackish

	fs := []string(nil)
	for _, f := range Fields {
		if x, ok := fkv[f.Name]; ok {
			p := ""
			if f.Proc != nil {
				p = f.Proc(x)
			} else {
				p = fmt.Sprintf("%v", x) // TODO some conversion
			}
			if p != "" {
				fs = append(fs, p)
			}
		}
	}

	fmt.Printf(strings.Join(fs, " ") + "\n") // TODO log.New(os.Stdout, "", 0).Printf() // log.LstdFlags?
}

func formatArg(x any) string { // TODO has to be method, has to be tunable
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
	default:
		s = fmt.Sprintf("%v", t)
	}
	if len([]rune(s)) > 1000 {
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
