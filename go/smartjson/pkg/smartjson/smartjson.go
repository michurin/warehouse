package smartjson

import (
	"fmt"
	"sort"
	"strings"
)

type Opts struct { // TODO: colors, separators; colors have to be considered on width calculation(!)
	Width  int
	Indent int
}

var reprStringReplacer = strings.NewReplacer(`"`, `\"`, `/`, `\/`, "\b", `\b`, "\f", `\f`, "\n", `\n`, "\r", `\r`, "\t", `\t`, `\`, `\\`)

func reprString(s string) string {
	return `"` + reprStringReplacer.Replace(s) + `"`
}

func max(a int, s string) int {
	b := len(s)
	if b > a {
		return b
	}
	return a
}

func marshal(level int, inp any, opts *Opts) (string, string) {
	ppfx := strings.Repeat(" ", opts.Indent*level)
	level++
	pfx := strings.Repeat(" ", opts.Indent*level)
	switch val := inp.(type) {
	case map[string]any:
		kk := make([]string, 0, len(val))
		for k := range val {
			kk = append(kk, k)
		}
		sort.Strings(kk)
		aa := make([]string, len(kk)) // aa — strings with single line representation for short single line relust
		bb := make([]string, len(kk)) // bb — multi line result with multi line parts
		cc := make([]string, len(kk)) // cc — multi line result with single line parts
		f := true                     // flag — is full set
		w := 0
		for i, k := range kk {
			s, m := marshal(level, val[k], opts)
			f = f && s != ""
			r := reprString(k)
			aa[i] = r + ": " + s
			bb[i] = pfx + r + ": " + m
			cc[i] = pfx + r + ": " + s
			w = max(w, cc[i])
		}
		t := bb
		s := ""
		if f {
			if w < opts.Width {
				t = cc
			}
			s = "{" + strings.Join(aa, ", ") + "}"
		}
		return s, "{\n" + strings.Join(t, ", ") + "\n" + ppfx + "}"
	case []any:
		aa := make([]string, len(val)) // the same naming as above
		bb := make([]string, len(val))
		cc := make([]string, len(val))
		f := true
		w := 0
		for i, v := range val {
			s, m := marshal(level, v, opts)
			f = f && s != ""
			aa[i] = s
			bb[i] = pfx + m
			cc[i] = pfx + s
			w = max(w, cc[i])
		}
		t := bb
		s := ""
		if f {
			if w < opts.Width {
				t = cc
			}
			s = "[" + strings.Join(aa, ", ") + "]"
		}
		return s, "[\n" + strings.Join(t, ",\n") + "\n" + ppfx + "]"
	case string:
		o := reprString(val)
		return o, o
	case float64:
		o := fmt.Sprintf("%g", val)
		return o, o
	case bool:
		o := "false"
		if val {
			o = "true"
		}
		return o, o
	case nil:
		return "null", "null"
	default:
		o := reprString(fmt.Sprintf("%v", val))
		return o, o
	}
}

func Marshal(inp any, opts *Opts) string {
	s, m := marshal(0, inp, opts)
	if len(s) < opts.Width {
		return s
	}
	return m
}
