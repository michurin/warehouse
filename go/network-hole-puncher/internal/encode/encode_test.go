package encode_test

import (
	"testing"

	"github.com/michurin/warehouse/go/network-hole-puncher/internal/encode"
)

func TestPackUnpack(t *testing.T) { // TODO Naive. BTW the good place for fuzzing
	a := "OK"
	d, err := encode.Pack([]byte(a))
	if err != nil {
		t.Fail()
	}
	b, err := encode.Unpack(d)
	if err != nil {
		t.Fail()
	}
	if string(b) != a {
		t.Fail()
	}
}
