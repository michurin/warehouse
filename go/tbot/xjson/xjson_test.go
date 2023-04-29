package xjson_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/warehouse/go/tbot/xjson"
)

func TestJsonToEnv(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		x := map[string]any{
			"a": nil,
			"b": false,
			"c": true,
			"d": float64(1),
			"e": "text",
			"f": []any{"element"},
			"g": map[string]any{
				"h": "sub",
			},
		}
		env, err := xjson.JsonToEnv(x)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"tg_b=false",
			"tg_c=true",
			"tg_d=1",
			"tg_e=text",
			"tg_f_0=element",
			"tg_g_h=sub",
		}, env)
	})
	t.Run("invalidType", func(t *testing.T) {
		x := float32(1)
		env, err := xjson.JsonToEnv(x)
		assert.Error(t, err)
		require.Nil(t, env)
	})
	t.Run("invalidTypeInSlice", func(t *testing.T) {
		x := []any{float32(1)}
		env, err := xjson.JsonToEnv(x)
		assert.Error(t, err)
		require.Nil(t, env)
	})
	t.Run("invalidTypeInMap", func(t *testing.T) {
		x := map[string]any{"k": float32(1)}
		env, err := xjson.JsonToEnv(x)
		assert.Error(t, err)
		require.Nil(t, env)
	})
}
