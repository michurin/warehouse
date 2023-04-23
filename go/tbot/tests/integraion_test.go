package tests_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_justCall(t *testing.T) {
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/botMORN/xMorn", r.URL.String())
		require.Equal(t, "applicatoin/x-morn", r.Header.Get("content-type"))
		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		require.Equal(t, []byte("{request}"), body)
		_, err = w.Write([]byte("{response}"))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer tg.Close()

	ctx := context.Background()

	bot := xbot.Bot{
		APIOrigin: tg.URL,
		Token:     "MORN",
		Client:    http.DefaultClient,
	}
	body, err := bot.API(ctx, &xbot.Request{
		Method:      "xMorn",
		ContentType: "applicatoin/x-morn",
		Body:        []byte("{request}"),
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("{response}"), body)
}
