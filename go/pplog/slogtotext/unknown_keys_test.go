package slogtotext

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnknownPairs_wanishing(t *testing.T) {
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

func TestUnknownPairs_types(t *testing.T) {
	data := any(nil)
	d := json.NewDecoder(strings.NewReader(`{"s":"x", "i":1, "t":true, "f":false, "n":null, "a":["x",1], "o":{"x":1}}`))
	d.UseNumber()
	err := d.Decode(&data)
	require.NoError(t, err)
	x := unknowPairs("", nil, data)
	assert.Equal(t, []unknownPair{
		{K: "a.0", V: "x"},
		{K: "a.1", V: "1"},
		{K: "f", V: "false"},
		{K: "i", V: "1"},
		{K: "n", V: "null"},
		{K: "o.x", V: "1"},
		{K: "s", V: "x"},
		{K: "t", V: "true"},
	}, x)
}

func TestUnknownPairs_invalidType(t *testing.T) {
	assert.Equal(t, []unknownPair{{K: "k", V: "UNKNOWN TYPE int8"}}, unknowPairs("k", nil, int8(1)))
}
