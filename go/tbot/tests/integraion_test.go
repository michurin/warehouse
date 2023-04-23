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

type apiAct struct {
	isJSON   bool
	request  string
	response []byte
}

func TestLoop(t *testing.T) {
	simpleUpdates := []apiAct{
		{
			true,
			`{"offset":0,"timeout":30}`,
			file(t, "data/get_update.json"),
		},
		{
			true,
			`{"offset":501,"timeout":30}`,
			nil,
		},
	}
	for _, cs := range []struct {
		name   string
		script string
		api    map[string][]apiAct
	}{
		{
			name:   "simple_text",
			script: "scripts/just_ok.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						true,
						fileStr(t, "data/send_message_request.json"),
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
		{
			name:   "media_jpeg",
			script: "scripts/media_jpeg.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendPhoto": {
					{
						false,
						"TODO!", // TODO
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
		{
			name:   "media_mp4",
			script: "scripts/media_mp4.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendVideo": {
					{
						false,
						"TODO!", // TODO
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tgURL, tgClose := botServer(t, cancel, cs.api)
			defer tgClose()

			bot := &xbot.Bot{
				APIOrigin: tgURL,
				Token:     "MORN",
				Client:    http.DefaultClient,
			}

			command := &xproc.Cmd{
				InterruptDelay: time.Second,
				KillDelay:      time.Second,
				Command:        cs.script,
				Cwd:            ".",
			}

			err := app.Loop(ctx, bot, command)
			require.Error(t, err)
			require.Contains(t, err.Error(), "context canceled") // like "api: client: Post \"http://127.0.0.1:34241/botMORN/getUpdates\": context canceled"
		})
	}
}

func botServer(t *testing.T, cancel context.CancelFunc, api map[string][]apiAct) (string, func()) {
	t.Helper()
	testDone := make(chan struct{})
	steps := map[string]int{} // it looks ugly, however we can use it without locks
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		bodyBytes, err := io.ReadAll(r.Body)
		body := string(bodyBytes)
		require.NoError(t, err)

		url := r.URL.String()
		n := steps[url]
		a := api[url][n]
		steps[url] = n + 1
		if a.isJSON {
			require.Equal(t, "application/json", r.Header.Get("content-type"))
			require.JSONEq(t, a.request, body)
		} else {
			require.Contains(t, r.Header.Get("content-type"), "multipart/form-data")
			// TODO check body!
		}
		if a.response == nil {
			cancel()
			<-testDone
		}
		_, err = w.Write(a.response)
		require.NoError(t, err)
	}))
	return tg.URL, func() {
		close(testDone)
		tg.Close()
	}
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
