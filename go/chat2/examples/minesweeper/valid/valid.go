package valid

import (
	"fmt"
	"strings"

	"github.com/michurin/warehouse/go/chat2/text"
)

func Open(w, h, x, y int, cid, name, clr string) error {
	m := []string(nil)
	if x < 0 {
		m = append(m, fmt.Sprintf("invalid x: %d", x))
	}
	if x >= w {
		m = append(m, fmt.Sprintf("invalid x: %d", x))
	}
	if y < 0 {
		m = append(m, fmt.Sprintf("invalid y: %d", y))
	}
	if y >= h {
		m = append(m, fmt.Sprintf("invalid y: %d", y))
	}
	if err := simple(cid, 33, 127, 24); err != nil {
		m = append(m, fmt.Sprintf("invalid CID: %s", err.Error()))
	}
	if err := Color(clr); err != nil {
		m = append(m, err.Error())
	}
	if text.SanitizeText(name, 10, "") == "" { // TODO 10 — constant
		m = append(m, fmt.Sprintf("invalid name"))
	}
	if m != nil {
		return fmt.Errorf(strings.Join(m, ", "))
	}
	return nil
}

func Color(c string) error {
	t := []rune(c)
	if len(t) != 7 {
		return fmt.Errorf("invalid color: %s", c)
	}
	if t[0] != '#' {
		return fmt.Errorf("invalid color: %s", c)
	}
	for _, e := range t[1:] {
		if !((e >= '0' && e <= '9') || (e >= 'a' && e <= 'f') || (e >= 'A' && e <= 'F')) {
			return fmt.Errorf("invalid color: %s", c)
		}
	}
	return nil
}

func simple(x string, r1, r2 rune, l int) error {
	t := []rune(x)
	if len(t) != l {
		return fmt.Errorf("invalid len: %d: %s", len(t), x)
	}
	for i, c := range t {
		if c < r1 || c >= r2 {
			return fmt.Errorf("invalid char at %d: %s", i, x)
		}
	}
	return nil
}
