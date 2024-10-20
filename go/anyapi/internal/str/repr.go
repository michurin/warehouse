package str

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type ColorScheme struct { // TODO move to file
	key,
	keyClose,
	index,
	indexClose,
	xtrue,
	xfalse,
	xnil,
	padding,
	comment,
	commentClose string
}

type ValueWithComment struct { // только для однострочных значений
	Value,
	Comment string
}

func reprStep(x any, prefix string, cs ColorScheme) (string, bool) { //nolint:funlen,gocognit,cyclop // long switch is ok for parsers
	nextPrefix := "  " + prefix
	switch v := x.(type) {
	case string:
		return v, false // TODO multiline
	case float64, float32:
		return fmt.Sprintf("%.20g", v), false
	case bool:
		if v {
			return cs.xtrue, false
		} else {
			return cs.xfalse, false
		}
	case nil:
		return cs.xnil, false
	case []any:
		if v == nil {
			return "nil", false
		}
		if len(v) == 0 {
			return "[]", false
		}
		res := []string{}
		prec := int(math.Ceil(math.Log10(float64(len(v)))))
		for i, e := range v {
			str, mul := reprStep(e, nextPrefix, cs)
			sep := " "
			if mul {
				sep = "\n"
			}
			res = append(res, fmt.Sprintf("%s%s%*d%s%s%s", prefix, cs.index, prec, i, cs.indexClose, sep, str))
		}
		return strings.Join(res, "\n"), true
	case map[string]any:
		if v == nil {
			return "nil", false
		}
		if len(v) == 0 {
			return "{}", false
		}
		kk := []string{}
		mk := 0
		for k := range v {
			kk = append(kk, k)
			cl := len([]rune(k))
			if mk < cl {
				mk = cl
			}
		}
		sort.Strings(kk)
		res := []string{}
		for i, k := range kk {
			str, mul := reprStep(v[k], nextPrefix, cs)
			sep := "\n"
			if !mul {
				pd := mk - len([]rune(k))
				switch pd {
				case 0:
					sep = " "
				case 1:
					sep = "  "
				default:
					dot := " "
					if i%2 == 0 {
						dot = cs.padding + "·" + cs.commentClose
					}
					sep = " " + strings.Repeat(dot, mk-len([]rune(k))-1) + " "
				}
			}
			res = append(res, fmt.Sprintf("%s%s%s%s%s%s", prefix, cs.key, k, cs.keyClose, sep, str))
		}
		return strings.Join(res, "\n"), true
	case ValueWithComment:
		return v.Value + " " + cs.comment + v.Comment + cs.commentClose, false
	default:
		return fmt.Sprintf("%T", v), false
	}
}

func Repr(x any) string {
	cs := ColorScheme{
		key:          "\033[32m",
		keyClose:     "\033[0m",
		index:        "\033[33m",
		indexClose:   "\033[0m",
		xtrue:        "\033[1;92mtrue\033[0m",
		xfalse:       "\033[1;91mfalse\033[0m",
		xnil:         "\033[1;93mnull\033[0m",
		padding:      "\033[2;33m",
		comment:      "\033[1;33m",
		commentClose: "\033[0m",
	}
	str, _ := reprStep(x, "", cs)
	return str
}
