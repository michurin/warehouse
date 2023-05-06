package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func File(t *testing.T, f string) []byte {
	t.Helper()
	if f == "" {
		return nil
	}
	data, err := os.ReadFile(f)
	require.NoError(t, err, f)
	return data
}

func FileStr(t *testing.T, f string) string {
	t.Helper()
	return string(File(t, f))
}
