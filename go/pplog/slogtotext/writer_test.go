package slogtotext_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"pplog/slogtotext"

	"github.com/stretchr/testify/assert"
)

// TODO just the simplest case
// TODO just example how to avoid <no value>

func ExamplePPLog_dealingWithUnknownKeysAndInvalidData() {
	w := slogtotext.PPLog(
		os.Stdout,
		`INVALID JSON: {{. | printf "%q"}}`,
		`{{.time | tmf "2006-01-02T15:04:05Z" "15:04:05" }} [{{.pid}}] {{.msg}}{{range .UNKNOWN}} {{.K}}={{.V}}{{end}}`,
		map[string]any{
			"msg":  struct{}{},
			"pid":  struct{}{},
			"time": struct{}{},
		},
		nil,
		0,
	)
	w.Write([]byte(`
[{ "invalid json" ]}
{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "begin", "unknown": "xx"}
{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "keys", "xx": {"k": "v", "n": null, "f": 1.2, "a": [1, 2]}}
{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "end"}
`))
	// output:
	// INVALID JSON: ""
	// INVALID JSON: "[{ \"invalid json\" ]}"
	// 23:00:00 [11] begin unknown=xx
	// 23:00:00 [11] keys xx.a.0=1 xx.a.1=2 xx.f=1.2 xx.k=v xx.n=nil
	// 23:00:00 [11] end
}

func nthPerm[T any](n int, a []T) []T {
	idx := make([]int, len(a)) // n%m, n%(m-1),... , n%1; filed from end to beginning
	for i := range a {
		j := i + 1
		idx[len(a)-i-1] = n % j
		n /= j
	}
	e := make([]T, len(a)) // copy
	r := make([]T, len(a)) // result
	copy(e, a)
	for i := range a {
		m := idx[i]
		r[i] = e[m]
		e = append(e[:m], e[m+1:]...)
	}
	return r
}

func TestPPLog_parts(t *testing.T) {
	inputX := [][]byte{
		[]byte(`[{ "invalid json" ]}`),
		[]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "xxxxxxxxxxxxxxxxxxxxxxxxxxx"}`), // 80 bytes
		[]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "begin", "unknown": "xx"}`),      // 75 bytes
		[]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "end"}`),
	}
	outputX := []string{
		`INVALID JSON: "[{ \"invalid json\" ]}"`,
		string(inputX[1]), // too long string wont be formatted
		"23:00:00 [11] begin unknown=xx",
		"23:00:00 [11] end",
	}
	for n := 0; n < 24; n++ {
		input := bytes.Join(append(nthPerm(n, inputX), nil), []byte{'\n'})
		output := strings.Join(append(nthPerm(n, outputX), ""), "\n")
		for i := range input {
			out := new(bytes.Buffer)
			w := slogtotext.PPLog(
				out,
				`INVALID JSON: {{. | printf "%q"}}`,
				`{{.time | tmf "2006-01-02T15:04:05Z" "15:04:05" }} [{{.pid}}] {{.msg}}{{range .UNKNOWN}} {{.K}}={{.V}}{{end}}`,
				map[string]any{"msg": struct{}{}, "pid": struct{}{}, "time": struct{}{}},
				nil,
				77, // consider as JSON strings up to 77 bytes long
			)
			w.Write(input[:i])
			w.Write(input[i:])
			assert.Equal(t, output, out.String())
		}
	}
}
