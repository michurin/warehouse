package str_test

import (
	"fmt"
	"testing"

	"anyapi/internal/str"
)

func TestRepr(t *testing.T) { //nolint:paralleltest
	x := str.Repr(map[string]any{
		"x":        7e16,
		"y":        map[string]any{"the": 1.},
		"longlong": "one",
		"longlon":  "two",
		"z": map[string]any{
			"a": []any{map[string]any{"A": "a"}, "B", "C"},
			"b": []any{1., true, false, []any{}, []any(nil), int64(7), str.ValueWithComment{
				Value:   "val",
				Comment: "comment",
			}},
		},
	})
	fmt.Println(x) // TODO assert
}
