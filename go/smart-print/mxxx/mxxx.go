package mxxx

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	envVarStack     = "MXXX_STACK"
	envVarStderr    = "MXXX_STDERR"
	envVarDumpDepth = "MXXX_DUMP_DEPTH"
)

type (
	internalString string
	internalInt    int
)

var (
	gcount       = 0
	glock        = new(sync.Mutex)
	colors       = [2]string{"\033[1;30;102m", "\033[1;92m"}
	OutputStream = io.Writer(nil) // just for testing and very extreme things; suppresses MXXX_STDERR effect
)

func dumpDepth() int {
	e, ok := os.LookupEnv(envVarDumpDepth)
	if !ok {
		return 3 // default
	}
	x, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		panic("INVALID VARIABLE " + envVarDumpDepth + ": " + err.Error())
	}
	return int(x)
}

func stackLimit() int {
	e, ok := os.LookupEnv(envVarStack)
	if !ok {
		return 1
	}
	x, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		panic("INVALID VARIABLE " + envVarStack + ": " + err.Error())
	}
	return int(x)
}

func removeLongestCommonPrefix(s, pfx string) (string, bool) {
	// TODO consider UTF and multi-byte chars
	// TODO behavior if len(pfx)==0?
	for i := range len(pfx) {
		if i >= len(s) {
			return "", i > 0
		}
		if pfx[i] != s[i] {
			return s[i:], i > 0
		}
	}
	return s[len(pfx):], len(pfx) > 0
}

func writeStack(out []byte) ([]byte, bool) {
	newLine := false
	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	pc := make([]uintptr, 1024)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
LOOP:
	for n := range stackLimit() {
		frame, more := frames.Next()
		file := frame.File
		file, ok := removeLongestCommonPrefix(file, cwd)
		if ok {
			file, _ = strings.CutPrefix(file, "/")
		}
		file, ok = strings.CutPrefix(file, home)
		if ok {
			file, _ = strings.CutPrefix(file, "/")
		}
		fn := frame.Function
		for range 3 { // TODO configurable
			i := strings.Index(fn, "/")
			if i < 0 {
				break
			}
			fn = fn[i+1:]
		}
		if n > 0 {
			out = append(out, '\n')
			newLine = true
		}
		out = fmt.Appendf(out, "%s%s:%d\033[0m \033[33m%s\033[0m", colors[min(n, 1)], file, frame.Line, fn)
		switch frame.Function {
		case "main.main", "runtime.goexit", "testing.tRunner", "testing.runExample":
			break LOOP
		}
		if !more {
			break
		}
	}
	return out, newLine
}

func appendUnlessNL(out []byte, c byte) []byte {
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, c)
	}
	return out
}

func writeOutput(out []byte) {
	s := OutputStream
	if s == nil {
		s = os.Stdout
		e, ok := os.LookupEnv(envVarStderr)
		if ok && len(e) > 0 {
			s = os.Stderr
		}
	}
	_, _ = io.Copy(s, bytes.NewReader(out))
}

func p(args ...any) {
	out, newLineInStack := writeStack(nil)
	for _, a := range args {
		s := ""
		c := ""
		switch v := a.(type) {
		case wrapper:
			sp := spew.NewDefaultConfig() // TODO DUP
			sp.SortKeys = true
			sp.DisableCapacities = true
			sp.DisablePointerAddresses = true
			sp.DisableMethods = true
			sp.DisablePointerMethods = true
			sp.Indent = "  "
			sp.MaxDepth = dumpDepth()
			sp.ContinueOnMethod = false
			s = strings.TrimSpace(sp.Sdump(v.v))
			c = "33"
		case nil:
			s = "nil"
			c = "105;1;30"
		case internalString:
			s = string(v)
			c = "101;30"
		case internalInt:
			s = fmt.Sprintf("%d", v)
			c = "103;30"
		case string:
			s = strings.TrimSpace(v)
			c = "95"
		case error:
			s = v.Error()
			c = "103;41;1"
		case fmt.Stringer:
			s = v.String()
			c = "92"
		case fmt.GoStringer:
			s = v.GoString()
			c = "94"
		default:
			sp := spew.NewDefaultConfig()
			sp.SortKeys = true
			sp.DisableCapacities = true
			sp.DisablePointerAddresses = true
			sp.DisableMethods = true
			sp.DisablePointerMethods = true
			sp.Indent = "  "
			sp.MaxDepth = dumpDepth()
			sp.ContinueOnMethod = false
			s = strings.TrimSpace(sp.Sdump(v))
			c = "33"
		}
		if newLineInStack || strings.Contains(s, "\n") {
			out = appendUnlessNL(out, '\n')
			x := strings.Split(s, "\n")
			for _, e := range x {
				out = fmt.Appendf(out, "\033[%sm%s\033[0m\n", c, e)
			}
		} else {
			out = appendUnlessNL(out, ' ')
			out = fmt.Appendf(out, "\033[%sm%s\033[0m", c, s)
		}
	}
	out = appendUnlessNL(out, '\n')
	writeOutput(out)
}

type wrapper struct {
	v any
}

// DUMP is a wrapper to force data dumping, ignoring [fmt.Stringer], and [fmt.GoStringer] implementations.
func DUMP(x any) wrapper {
	return wrapper{v: x}
}

// P dumps arguments
func P(args ...any) {
	p(args...)
}

// NOERR dumps and exits if first argument is non-nil error. Otherwise it does nothing.
func NOERR(args ...any) {
	if args[0] == nil {
		return
	}
	p(args...)
	os.Exit(22)
}

// PX wraps call.
func PX[T any](x T, args ...any) T {
	p(append([]any{x}, args...)...)
	return x
}

// EXIT dumps arguments and exits.
func EXIT(args ...any) {
	p(args...)
	os.Exit(22)
}

// SLEEP sleeps.
func SLEEP(d time.Duration, args ...any) {
	out, _ := writeStack(nil)
	out = fmt.Appendf(out, "\n\033[105;1;35m\033[K[SLEEP] %v\033[0m\n", d)
	for _, a := range args {
		out = fmt.Appendf(out, "\033[35;1m[SLEEP] %v\033[0m\n", a)
	}
	writeOutput(out)
	time.Sleep(d)
}

// GO unparallelize goroutines to make logs more readable.
func GO(args ...any) func() {
	glock.Lock()
	gcount++
	p(append([]any{internalString("STARTING GOROUTINE"), internalInt(gcount)}, args...)...)
	return func() {
		p(internalString("STOPPING GOROUTINE"), internalInt(gcount))
		glock.Unlock()
	}
}
