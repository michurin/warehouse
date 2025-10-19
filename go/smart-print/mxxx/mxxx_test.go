package mxxx_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"mxxx/mxxx"
)

func eqFile(t *testing.T, actual []byte, filename string) {
	t.Helper()
	os.WriteFile(filename+"-", actual, 0o644)
	buff, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(buff, actual) {
		t.Error("not equal")
	}
}

func TestP(t *testing.T) {
	out := new(bytes.Buffer)
	mxxx.OutputStream = io.MultiWriter(out, os.Stderr)
	os.Setenv("MXXX_STACK", "10")

	out.Reset()
	mxxx.P(1)
	eqFile(t, out.Bytes(), "test_p10.json")

	out.Reset()
	mxxx.P(context.TODO())
	eqFile(t, out.Bytes(), "test_context.json")

	out.Reset()
	mxxx.P(mxxx.DUMP(context.TODO()))
	eqFile(t, out.Bytes(), "test_context_dump.json")
}

// for i in *.json-; do cat $i >${i%-}; done
