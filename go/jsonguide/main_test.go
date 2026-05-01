package main_test

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	main "github.com/michurin/warehouse/go/jsonguide"
)

//go:embed testdata
var testdata embed.FS

type mockFileInfo struct{}

func (*mockFileInfo) Name() string       { panic("not implemented") }
func (*mockFileInfo) Size() int64        { panic("not implemented") }
func (*mockFileInfo) Mode() fs.FileMode  { return 0 }
func (*mockFileInfo) ModTime() time.Time { panic("not implemented") }
func (*mockFileInfo) IsDir() bool        { panic("not implemented") }
func (*mockFileInfo) Sys() any           { panic("not implemented") }

type mockFile struct{}

func (*mockFile) Stat() (fs.FileInfo, error) { return new(mockFileInfo), nil }
func (*mockFile) Read([]byte) (int, error)   { panic("not implemented") }
func (*mockFile) Close() error               { panic("not implemented") }

func TestApp(t *testing.T) {
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
			args := []string(nil)
			if strings.HasPrefix(c, "shallow_") {
				args = append(args, "-s")
			}
			if strings.HasPrefix(c, "help_") {
				args = append(args, "-h")
			}
			rc := main.App(bytes.NewReader(data), buf, &mockFile{}, args)
			expErr := strings.HasPrefix(c, "wrong_")
			if (rc != 0 || expErr) && (rc != 1 || !expErr) { // !((rc == 0 && !expErr) || (rc == 1 && expErr))
				t.Error("got: rc=", rc)
			}
			if buf.String() != string(expectedOutput) {
				t.Error("args: " + strings.Join(args, " ") + "\ngot:\n" + buf.String() + "\texpected:\n" + string(expectedOutput))
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
