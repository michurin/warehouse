package tests_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/michurin/warehouse/go/tbot/app"
	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xproc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_justCall(t *testing.T) {
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/botMORN/xMorn", r.URL.String())
		require.Equal(t, "application/x-morn", r.Header.Get("content-type"))
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
		ContentType: "application/x-morn",
		Body:        []byte("{request}"),
	})

	require.NoError(t, err)
	assert.Equal(t, []byte("{response}"), body)
}

func TestLoop_simpleMessage(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testDone := make(chan struct{})

	updateCount := 0 // we can use it without locks
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "application/json", r.Header.Get("content-type"))
		bodyBytes, err := io.ReadAll(r.Body)
		body := string(bodyBytes)
		require.NoError(t, err)
		respFile := ""
		switch r.URL.String() {
		case "/botMORN/getUpdates":
			switch updateCount {
			case 0:
				require.JSONEq(t, `{"offset":0,"timeout":30}`, body)
				respFile = "data/get_update.json"
			case 1:
				require.JSONEq(t, `{"offset":501,"timeout":30}`, body)
				cancel()
				<-testDone
			default:
				t.Fatal("unexpected getUpdates", updateCount)
			}
			updateCount++
		case "/botMORN/sendMessage":
			require.JSONEq(t, fileStr(t, "data/send_message_request.json"), body)
		default:
			t.Fatal("unexpected api Method", r.URL.String())
		}
		_, err = w.Write(file(t, respFile))
		require.NoError(t, err)
	}))
	defer tg.Close()

	bot := &xbot.Bot{
		APIOrigin: tg.URL,
		Token:     "MORN",
		Client:    http.DefaultClient,
	}

	command := &xproc.Cmd{
		InterruptDelay: time.Second,
		KillDelay:      time.Second,
		Command:        "scripts/just_ok.sh",
		Cwd:            ".",
	}

	err := app.Loop(ctx, bot, command)
	close(testDone)
	require.Error(t, err) // TODO check "context canceled"
}

func file(t *testing.T, f string) []byte {
	t.Helper()
	if f == "" {
		return nil
	}
	data, err := os.ReadFile(f)
	require.NoError(t, err, f)
	return data
}

func fileStr(t *testing.T, f string) string {
	return string(file(t, f))
}
