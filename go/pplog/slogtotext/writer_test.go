package slogtotext_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"pplog/slogtotext"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExamplePPLog_justShowingIdea() {
	w := slogtotext.PPLog(os.Stdout, "", `{{.time}} [{{.pid}}] {{.message}}`, nil, nil, 0)
	w.Write([]byte(`{"time":"12:00", "pid":11, "message":"OK"}` + "\n"))
	// output:
	// 12:00 [11] OK
}

// TODO remove fields, using known keys
// TODO invalid json
// TODO format dates
// TODO just example how to avoid <no value>

func ExamplePPLog_dealingWithUnknownKeysAndInvalidData() {
	w := slogtotext.PPLog(
		os.Stdout,
		`INVALID JSON: {{. | printf "%q"}}`,
		`{{.time | tmf "2006-01-02T15:04:05Z" "15:04:05" }} [{{.pid}}] {{.msg}}{{range .UNKNOWN}} {{.K}}={{.V}}{{end}}`,
		nil,
		nil,
		0,
	)
	w.Write([]byte(`[{ "invalid json" ]}` + "\n"))
	w.Write([]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "message a", "unknown": "xx"}` + "\n"))
	w.Write([]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "keys", "xx": {"k": "v", "n": null, "f": 1.2, "a": [1, 2]}}` + "\n"))
	// output:
	// INVALID JSON: "[{ \"invalid json\" ]}"
	// 23:00:00 [11] message a unknown=xx
	// 23:00:00 [11] keys xx.a.0=1 xx.a.1=2 xx.f=1.2 xx.k=v xx.n=nil
}

func nthPerm[T any](n int, a []T) []T {
	idx := make([]int, len(a)) // n%m, n%(m-1),... , n%1; filed from end to message aning
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
		[]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 10, "msg": "xxxxxxxxxxxxxxxxxxxxx"}`), // 76 bytes
		[]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 11, "msg": "message a", "x": "xx"}`),  // 75 bytes
		[]byte(`{"time": "2009-11-10T23:00:00Z", "pid": 12, "msg": "message b"}`),
	}
	outputX := []string{
		`INVALID JSON: "[{ \"invalid json\" ]}"`,
		string(inputX[1]), // too long string wont be formatted
		"23:00:00 [11] message a x=xx",
		"23:00:00 [12] message b",
	}
	for permN := 0; permN < 24; permN++ {
		input := bytes.Join(append(nthPerm(permN, inputX), nil), []byte{'\n'})
		output := strings.Join(append(nthPerm(permN, outputX), ""), "\n")
		for i := range input {
			out := new(bytes.Buffer)
			w := slogtotext.PPLog(
				out,
				`INVALID JSON: {{. | printf "%q"}}`,
				`{{.time | tmf "2006-01-02T15:04:05Z" "15:04:05" }} [{{.pid}}] {{.msg}}{{range .UNKNOWN}} {{.K}}={{.V}}{{end}}`,
				nil,
				nil,
				75, // consider as JSON strings up to 75 bytes long
			)
			n, err := w.Write(input[:i])
			require.NoError(t, err)
			assert.Equal(t, i, n)
			n, err = w.Write(input[i:])
			require.NoError(t, err)
			assert.Equal(t, len(input)-i, n)
			assert.Equal(t, output, out.String())
		}
	}
}
