package minlog

import (
	"path"
	"testing"
)

func TestCalleri_ok(t *testing.T) {
	r := caller(1)
	if r.Line != 9 || path.Base(r.File) != "caller_test.go" {
		t.Fail()
	}
}

func TestCaller_error(t *testing.T) {
	r := caller(100)
	if r.Line != 0 || r.File != "nofile" {
		t.Fail()
	}
}
