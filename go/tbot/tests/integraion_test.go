package tests_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/warehouse/go/tbot/app"
	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xproc"
)

// TODO setup xlog

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
						"--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"photo\"; filename=\"image.jpeg\"\r\nContent-Type: image/jpeg\r\n\r\n\xff\xd8\xff\r\n--BOUND--\r\n",
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
		{
			name:   "media_png",
			script: "scripts/media_png.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendPhoto": {
					{
						false,
						"--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"photo\"; filename=\"image.png\"\r\nContent-Type: image/png\r\n\r\n\x89PNG\r\n\x1a\n\r\n--BOUND--\r\n",
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
		{
			name:   "media_mp3",
			script: "scripts/media_mp3.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendAudio": {
					{
						false,
						"--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"audio\"; filename=\"audio.mpeg\"\r\nContent-Type: audio/mpeg\r\n\r\nID3\r\n--BOUND--\r\n",
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
		{
			name:   "media_ogg",
			script: "scripts/media_ogg.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": { // consider ogg as document, it seems it's not fully supported
					{
						false,
						"--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document\"\r\nContent-Type: application/ogg\r\n\r\nOggS\x00\r\n--BOUND--\r\n",
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
						"--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"video\"; filename=\"video.mp4\"\r\nContent-Type: video/mp4\r\n\r\n\x00\x00\x00\fftypmp4_\r\n--BOUND--\r\n",
						file(t, "data/send_message_response.json"),
					},
				},
			},
		},
		{
			name:   "media_pdf",
			script: "scripts/media_pdf.sh",
			api: map[string][]apiAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": {
					{
						false,
						"--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document\"\r\nContent-Type: application/pdf\r\n\r\n%PDF-\r\n--BOUND--\r\n",
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

// ---- move to tooling

type apiAct struct {
	isJSON   bool
	request  string
	response []byte
}

func botServer(t *testing.T, cancel context.CancelFunc, api map[string][]apiAct) (string, func()) {
	t.Helper()
	testDone := make(chan struct{})
	steps := map[string]int{} // it looks ugly, however we can use it without locks
	tg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		bodyBytes, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		body := string(bodyBytes)

		url := r.URL.String()
		n := steps[url]
		a := api[url][n] // TODO this panic is caught by server! so test wont fail!
		steps[url] = n + 1
		if a.isJSON {
			require.Equal(t, "application/json", r.Header.Get("content-type"))
			require.JSONEq(t, a.request, body)
		} else {
			ctype := r.Header.Get("content-type")
			require.Contains(t, ctype, "multipart/form-data")
			idx := strings.Index(ctype, "boundary=")
			assert.Greater(t, idx, -1)
			universal := strings.ReplaceAll(body, ctype[idx+9:], "BOUND")
			assert.Equal(t, a.request, universal)
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

// ----

func TestHttp(t *testing.T) { // curl -F works transparently as is
	for _, cs := range []struct {
		name string
		curl []string
		qs   string
		api  map[string][]apiAct
	}{
		{
			name: "curl_F",
			curl: []string{"-q", "-s", "-F", "user_id=10", "-F", "text=ok"},
			qs:   "",
			api: map[string][]apiAct{
				"/botMORN/someMethod": {{
					isJSON:   false,
					request:  "--BOUND\r\nContent-Disposition: form-data; name=\"user_id\"\r\n\r\n10\r\n--BOUND\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nok\r\n--BOUND--\r\n",
					response: []byte("done."),
				}},
			},
		},
		{
			name: "curl_d",
			curl: []string{"-q", "-s", "-d", "ok"},
			qs:   "?to=111",
			api: map[string][]apiAct{
				"/botMORN/sendMessage": {{
					isJSON:   true,
					request:  `{"chat_id":111, "text":"ok"}`,
					response: []byte("done."),
				}},
			},
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			tgURL, tgClose := botServer(t, nil, cs.api)
			defer tgClose()

			bot := &xbot.Bot{
				APIOrigin: tgURL,
				Token:     "MORN",
				Client:    http.DefaultClient,
			}

			h := app.Handler(bot)

			s := httptest.NewServer(h)

			cmd := exec.Command("curl", append(cs.curl, s.URL+"/x/someMethod"+cs.qs)...)
			var stdOut, stdErr bytes.Buffer
			cmd.Stdout = &stdOut
			cmd.Stderr = &stdErr
			err := cmd.Run()
			require.NoError(t, err)
			assert.Equal(t, "done.", stdOut.String())
			assert.Empty(t, stdErr.String())
		})
	}
}

// ----

func TestProc(t *testing.T) {
	ctx := context.Background()
	command := &xproc.Cmd{
		InterruptDelay: 200 * time.Millisecond,
		KillDelay:      200 * time.Millisecond,
		Command:        "sh",
		Cwd:            ".",
	}
	t.Run("finishNormally", func(t *testing.T) {
		data, err := command.Run(ctx, []string{"-c", "echo OK"}, []string{})
		assert.NoError(t, err)
		assert.Equal(t, []byte("OK\n"), data)
	})
	t.Run("sigint", func(t *testing.T) {
		data, err := command.Run(ctx, []string{"-c", "trap 'echo sigint_catched' SIGINT; sleep 1"}, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sigint_catched") // we have to see our stdout in error text
		assert.Nil(t, data)
	})
	t.Run("sigkill", func(t *testing.T) {
		data, err := command.Run(ctx, []string{"-c", "trap '' SIGINT; sleep 1"}, []string{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "-1") // exit code
		assert.Nil(t, data)
	})
}
