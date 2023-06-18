package minlog

import (
	"strings"
	"testing"
)

// only cases not covered by other tests

func TestFormatArg_false(t *testing.T) {
	s := formatArg(false)
	if s != "false" {
		t.Fail()
	}
}

func TestFormatArg_longRunes(t *testing.T) {
	s := formatArg(strings.Repeat("Ñ", 1001))
	if len(s) != 2000 || len([]rune(s)) != 1000 {
		t.Fail()
	}
}

func TestFormatArg_longBytes(t *testing.T) {
	s := formatArg(strings.Repeat("N", 1001))
	if len(s) != 1000 || len([]rune(s)) != 1000 {
		t.Fail()
	}
}
