package smartjson

import (
	"fmt"
	"sort"
	"strings"
)

type Opts struct {
	Width  int
	Indent int
}

var reprStringReplacer = strings.NewReplacer(`"`, `\"`, `/`, `\/`, "\b", `\b`, "\f", `\f`, "\n", `\n`, "\r", `\r`, "\t", `\t`, `\`, `\\`)

func reprString(s string) string {
	return `"` + reprStringReplacer.Replace(s) + `"`
}

func marshal(level int, inp any, opts *Opts) (string, string) {
	ppfx := "\n" + strings.Repeat(" ", opts.Indent*level)
	level++
	pfx := strings.Repeat(" ", opts.Indent*level)
	switch val := inp.(type) {
	case map[string]any:
		kk := make([]string, 0, len(val))
		for k := range val {
			kk = append(kk, k)
		}
		sort.Strings(kk)
		ss := make([]string, len(kk))
		mm := make([]string, len(kk))
		f := true
		for i, k := range kk {
			s, m := marshal(level, val[k], opts)
			r := reprString(k)
			mm[i] = pfx + r + ": " + m
			if s != "" {
				ss[i] = r + ": " + s
				t := pfx + r + ": " + s
				if len(t) < opts.Width {
					mm[i] = t
				}
			} else {
				f = false
			}
		}
		s := ""
		if f {
			s = "{" + strings.Join(ss, ", ") + "}"
		}
		return s, "{\n" + strings.Join(mm, ", ") + ppfx + "}"
	case []any:
		ss := make([]string, len(val))
		mm := make([]string, len(val))
		f := true
		for i, v := range val {
			s, m := marshal(level, v, opts)
			mm[i] = pfx + m
			if s != "" {
				ss[i] = s
				t := pfx + s
				if len(t) < opts.Width {
					mm[i] = t
				}
			} else {
				f = false
			}
		}
		s := ""
		if f {
			s = "[" + strings.Join(ss, ", ") + "]"
		}
		return s, "[\n" + strings.Join(mm, ",\n") + ppfx + "]"
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
