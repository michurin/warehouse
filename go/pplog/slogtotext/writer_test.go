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

//nolint:errcheck
func ExamplePPLog_justShowingIdea() {
	w := slogtotext.PPLog(os.Stdout, "", `{{.time}} [{{.pid}}] {{.message}}`, nil, nil, 0)
	w.Write([]byte(`{"time":"12:00", "pid":11, "message":"OK"}` + "\n"))
	// output:
	// 12:00 [11] OK
}

//nolint:errcheck
func ExamplePPLog_includeUnknownFields() {
	w := slogtotext.PPLog(os.Stdout, "", `{{.time}} [{{.pid}}] {{.message}}{{range .UNKNOWN}} {{.K}}={{.V}}{{end}}`, nil, nil, 0)
	w.Write([]byte(`{"time":"12:00", "pid":11, "message":"OK", "request_id": "xx", "g": {"a": "A", "b": "B"}}` + "\n"))
	// output:
	// 12:00 [11] OK g.a=A g.b=B request_id=xx
}

//nolint:errcheck
func ExamplePPLog_includeUnknownFieldsButSkipSomeOfThem() {
	w := slogtotext.PPLog(os.Stdout, "", `{{.time}} [{{.pid}}] {{.message}}{{range .UNKNOWN}} {{.K}}={{.V}}{{end}}`, map[string]any{
		"time":    struct{}{},
		"pid":     struct{}{},
		"message": struct{}{},
		"host":    struct{}{}, // skip field that is not showing up in template
		"g": map[string]any{
			"a": struct{}{}, // skip g.a
		},
	}, nil, 0)
	w.Write([]byte(`{"time":"12:00", "pid":11, "message":"OK", "request_id": "xx", "host": "sun", "g": {"a": "A", "b": "B"}}` + "\n"))
	// output:
	// 12:00 [11] OK g.b=B request_id=xx
}

//nolint:errcheck
func ExamplePPLog_invalidJSON() {
	w := slogtotext.PPLog(os.Stdout, `INVALID JSON: {{. | printf "%q"}}`, "", nil, nil, 0)
	w.Write([]byte(`{[} invalid json` + "\n"))
	// output:
	// INVALID JSON: "{[} invalid json"
}

//nolint:errcheck
func ExamplePPLog_formatTimestapm() {
	w := slogtotext.PPLog(os.Stdout, "", `{{.time | tmf "2006-01-02T15:04:05Z" "15:04:05"}} {{.message}}`, nil, nil, 0)
	w.Write([]byte(`{"time":"2009-11-10T23:00:00Z", "message":"OK"}` + "\n"))
	// output:
	// 23:00:00 OK
}

//nolint:errcheck
func ExamplePPLog_noValues() {
	w := slogtotext.PPLog(os.Stdout, "", `[{{.message}}]`, nil, nil, 0)
	w.Write([]byte("{}\n"))
	w = slogtotext.PPLog(os.Stdout, "", `{{if .message}}[{{.message}}]{{else}}[-]{{end}}`, nil, nil, 0)
	w.Write([]byte("{}\n"))
	// output:
	// [<no value>]
	// [-]
}

func nthPermutation[T any](n int, a []T) []T {
	idx := make([]int, len(a)) // n%m, n%(m-1),... , n%1; filed from end to begin
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

func TestNthPermutation(t *testing.T) { // just check
	a := []int{1, 2, 3}
	assert.Equal(t, []int{1, 2, 3}, nthPermutation(0, a))
	assert.Equal(t, []int{1, 3, 2}, nthPermutation(1, a))
	assert.Equal(t, []int{2, 1, 3}, nthPermutation(2, a))
	assert.Equal(t, []int{2, 3, 1}, nthPermutation(3, a))
	assert.Equal(t, []int{3, 1, 2}, nthPermutation(4, a))
	assert.Equal(t, []int{3, 2, 1}, nthPermutation(5, a))
	assert.Equal(t, []int{1, 2, 3}, nthPermutation(6, a)) // 0
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
		input := bytes.Join(append(nthPermutation(permN, inputX), nil), []byte{'\n'})
		output := strings.Join(append(nthPermutation(permN, outputX), ""), "\n")
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
