package minlog

import (
	"strconv"
	"strings"
)

func Color(next FieldFunc, code ...int) FieldFunc {
	c := make([]string, len(code))
	for i, v := range code {
		c[i] = strconv.Itoa(v)
	}
	s := "\033[" + strings.Join(c, ";") + "m"
	return func(r Record) string {
		t := next(r)
		if t == "" {
			return t
		}
		return s + t + "\033[0m"
	}
}

func Prefix(next FieldFunc, prefix string) FieldFunc {
	return func(r Record) string {
		t := next(r)
		if t == "" {
			return t
		}
		return prefix + t
	}
}
