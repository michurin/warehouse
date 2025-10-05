package mxxx

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const (
	envVarStack  = "MXXX_STACK"
	envVarStderr = "MXXX_STDERR"
)

var colors = [2]string{"\033[1;30;102m", "\033[92m"}

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

func writeStack(out []byte) []byte {
	cwd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	pc := make([]uintptr, 20)
	n := runtime.Callers(3, pc)
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
			out = append(out, ' ')
		}
		out = fmt.Appendf(out, "%s%s\033[0m \033[33m%s\033[93m:%d\033[0m", colors[min(n, 1)], fn, file, frame.Line)
		switch frame.Function {
		case "main.main", "runtime.goexit", "testing.tRunner", "testing.runExample":
			break LOOP
		}
		if !more {
			break
		}
	}
	return out
}

func appendUnlessNL(out []byte, c byte) []byte {
	if len(out) > 0 && out[len(out)-1] != '\n' {
		out = append(out, c)
	}
	return out
}

func writeOutput(out []byte) {
	s := os.Stdout
	e, ok := os.LookupEnv(envVarStderr)
	if ok && len(e) > 0 {
		s = os.Stderr
	}
	_, _ = io.Copy(s, bytes.NewReader(out))
}

func P(args ...any) {
	out := writeStack(nil)
	for _, a := range args {
		s := ""
		c := ""
		switch v := a.(type) {
		case string:
			s = strings.TrimSpace(v)
			c = "95"
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
			sp.MaxDepth = 3 // TODO tunable
			sp.ContinueOnMethod = true
			s = strings.TrimSpace(sp.Sdump(v))
			c = "96"
		}
		if strings.Contains(s, "\n") {
			out = appendUnlessNL(out, '\n')
			out = fmt.Appendf(out, "\033[%sm%s\033[0m\n", c, s)
		} else {
			out = appendUnlessNL(out, ' ')
			out = fmt.Appendf(out, "\033[%sm%s\033[0m", c, s)
		}
	}
	out = appendUnlessNL(out, '\n')
	writeOutput(out)
}

func SLEEP(d time.Duration) {
	out := writeStack(nil)
	out = fmt.Appendf(out, "\033[105;1;35m\033[K SLEEP %v\033[0m\n", d)
	writeOutput(out)
	time.Sleep(d)
}
