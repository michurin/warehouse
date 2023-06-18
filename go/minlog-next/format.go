package minlog

import (
	"fmt"
	"unicode/utf8"
)

func formatArg(x any) string {
	var s string
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
	r := []rune(s)
	if len(r) > 1000 { // length is limited in runes not in bytes
		s = string(r[:1000])
	}
	return s
}
