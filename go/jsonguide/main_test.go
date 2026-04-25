package main_test

import (
	"bytes"
	"embed"
	"os"
	"strings"
	"testing"

	main "github.com/michurin/warehouse/go/jsonguide"
)

//go:embed testdata
var testdata embed.FS

func TestMain(t *testing.T) {
	t.Parallel()
	fs, err := testdata.ReadDir("testdata")
	noerr(t, err)
	cases := []string(nil)
	for _, f := range fs {
		n := f.Name()
		if strings.HasSuffix(n, ".json") {
			cases = append(cases, n[:len(n)-5])
		}
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			t.Parallel()
			data, err := os.ReadFile("testdata/" + c + ".json")
			noerr(t, err)
			expectedOutput, err := os.ReadFile("testdata/" + c + ".out")
			noerr(t, err)
			buf := new(strings.Builder)
			rc := main.App(bytes.NewReader(data), buf, false)
			expectError := strings.HasPrefix(c, "wrong_")
			if !((rc == 0 && !expectError) || (rc == 1 && expectError)) {
				t.Error("got: rc=", rc)
			}
			if buf.String() != string(expectedOutput) {
				t.Error("got:\n" + buf.String() + "\texpected:\n" + string(expectedOutput))
			}
		})
	}
}

func noerr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err.Error())
	}
}
