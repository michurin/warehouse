package slogtotext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownPairs(t *testing.T) {
	x := unknowPairs("", map[string]any{
		"b": struct{}{},
		"c": map[string]any{"p": struct{}{}},
		"d": struct{}{},
		"e": map[string]any{"ea": struct{}{}},
	}, map[string]any{
		"a": "A",
		"b": "B", // skip: directly
		"c": "C", // skip: indirect: appears in path
		"d": map[string]any{
			"da": "DA", // doesn't skip: not full path, d only
			"db": "DB", // the same
		},
		"e": map[string]any{
			"ea": "EA", // skip
			"eb": "EB",
		},
	})
	assert.Equal(t, []unknownPair{
		{"a", "A"},
		{"d.da", "DA"},
		{"d.db", "DB"},
		{"e.eb", "EB"},
	}, x)
}
