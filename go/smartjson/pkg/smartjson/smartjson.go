package smartjson

import (
	"fmt"
	"sort"
	"strings"
)

type Opts struct {
	Width  int
	Indent int
	Theme  Theme
}

var reprStringReplacer = strings.NewReplacer(`"`, `\"`, `/`, `\/`, "\b", `\b`, "\f", `\f`, "\n", `\n`, "\r", `\r`, "\t", `\t`, `\`, `\\`)

func reprString(s string, q, body Str) Str {
	return Concat(q, Wrap(New(reprStringReplacer.Replace(s)), body), q)
}

func marshal(level int, inp any, opts *Opts) (Str, string) {
	ppfx := "\n" + strings.Repeat(" ", opts.Indent*level)
	level++
	pfx := Repeat(" ", opts.Indent*level)
	switch val := inp.(type) {
	case map[string]any:
		kk := make([]string, 0, len(val))
		for k := range val {
			kk = append(kk, k)
		}
		sort.Strings(kk)
		ss := make([]Str, len(kk))
		mm := make([]string, len(kk))
		f := true
		for i, k := range kk {
			s, m := marshal(level, val[k], opts)
			r := reprString(k, opts.Theme.StrKeyQuo, opts.Theme.StrKeyBody)
			mm[i] = pfx.String() + r.String() + opts.Theme.MapMuPa.String() + m
			if s.Len() != 0 {
				ss[i] = Concat(r, opts.Theme.MapSiPa, s)
				t := Concat(pfx, r, opts.Theme.MapMuPa, s)
				if t.Len() < opts.Width {
					mm[i] = t.String()
				}
			} else {
				f = false
			}
		}
		s := New("")
		if f {
			s = Concat(opts.Theme.MapSiBo, Join(ss, opts.Theme.MapSiSe), opts.Theme.MapSiBc)
		}
		return s, opts.Theme.MapMuBo.String() + "\n" + strings.Join(mm, opts.Theme.MapMuSe.String()) + ppfx + opts.Theme.MapMuBc.String()
	case []any:
		ss := make([]Str, len(val))
		mm := make([]string, len(val))
		f := true
		for i, v := range val {
			s, m := marshal(level, v, opts)
			mm[i] = pfx.String() + m
			if s.Len() != 0 {
				ss[i] = s
				t := Concat(pfx, s)
				if t.Len() < opts.Width {
					mm[i] = t.String()
				}
			} else {
				f = false
			}
		}
		s := New("")
		if f {
			s = Concat(opts.Theme.LstSiBo, Join(ss, opts.Theme.LstSiSe), opts.Theme.LstSiBc)
		}
		return s, opts.Theme.LstMuBo.String() + "\n" + strings.Join(mm, opts.Theme.LstMuSe.String()+"\n") + ppfx + opts.Theme.LstMuBc.String()
	case string:
		o := reprString(val, opts.Theme.StrValQuo, opts.Theme.StrValBody)
		return o, o.String()
	case float64:
		o := Wrap(New(fmt.Sprintf("%g", val)), opts.Theme.Flo)
		return o, o.String()
	case bool:
		if val {
			return opts.Theme.Tru, opts.Theme.Tru.String()
		}
		return opts.Theme.Fal, opts.Theme.Fal.String()
	case nil:
		return opts.Theme.Nul, opts.Theme.Nul.String()
	default:
		o := reprString(fmt.Sprintf("%v", val), opts.Theme.ErrorQuo, opts.Theme.ErrorBody)
		return o, o.String()
	}
}

func Marshal(inp any, opts *Opts) string {
	s, m := marshal(0, inp, opts)
	if s.Len() < opts.Width {
		return s.String()
	}
	return m
}
