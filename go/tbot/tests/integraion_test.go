package tests_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michurin/warehouse/go/tbot/tests/apiserver"
	"github.com/michurin/warehouse/go/tbot/xbot"
	"github.com/michurin/warehouse/go/tbot/xctrl"
	"github.com/michurin/warehouse/go/tbot/xloop"
	"github.com/michurin/warehouse/go/tbot/xproc"
)

// TODO setup xlog

func TestAPI_justCall(t *testing.T) {
	/* case
	tg        bot
	|         |
	|<--req---|
	|---resp->|
	*/
	tgURL, tgClose := apiserver.APIServer(t, nil, map[string][]apiserver.APIAct{
		"/botMORN/xMorn": {{
			IsJSON:   true,
			Request:  `{"ok":1}`,
			Response: []byte(`{"response":1}`),
		}},
	})
	defer tgClose()

	ctx := context.Background()

	bot := buildBot(tgURL)

	body, err := bot.API(ctx, &xbot.Request{
		Method:      "xMorn",
		ContentType: "application/json",
		Body:        []byte(`{"ok":1}`),
	})

	require.NoError(t, err)
	assert.JSONEq(t, `{"response":1}`, string(body))
}

func TestLoop(t *testing.T) {
	/* cases
	tg        bots loop
	|         |
	|<--req---| (call for update)
	|---resp->|
	|         |
	|         |--exec-->| script
	|         |<-stdout-|
	|         |
	|<--req---| (call for update)
	|---resp->|
	|<--req---| (and send response from script)
	|---resp->| (the order of update and send doesn't metter)
	*/
	simpleUpdates := []apiserver.APIAct{
		{
			IsJSON:  true,
			Request: `{"offset":0,"timeout":30}`,
			Response: []byte(`{"ok": true, "result": [{"update_id": 500, "message": {
"message_id": 100,
"from": {"id": 1500, "is_bot": false, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "language_code": "en"},
"chat": {"id": 1500, "first_name": "Alex", "last_name": "Morn", "username": "AlexMorn", "type": "private"},
"date": 1682222222,
"text": "word"}}]}`),
		},
		{
			IsJSON:   true,
			Request:  `{"offset":501,"timeout":30}`,
			Response: nil,
		},
	}
	sendMessageResponseJSON := []byte(`{"ok": true, "result": {}}`)
	for _, cs := range []struct {
		name   string
		script string
		api    map[string][]apiserver.APIAct
	}{
		{
			name:   "simple_text",
			script: "scripts/just_ok.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "text": "ok\n"}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "simple_text",
			script: "scripts/preformatted_ok.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendMessage": {
					{
						IsJSON:   true,
						Request:  `{"chat_id": 1500, "text": "ok", "entities": [{"type": "pre", "offset": 0, "length": 2}]}`,
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "media_jpeg",
			script: "scripts/media_jpeg.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendPhoto": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"photo\"; filename=\"image.jpeg\"\r\nContent-Type: image/jpeg\r\n\r\n\xff\xd8\xff\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "media_png",
			script: "scripts/media_png.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendPhoto": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"photo\"; filename=\"image.png\"\r\nContent-Type: image/png\r\n\r\n\x89PNG\r\n\x1a\n\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "media_mp3",
			script: "scripts/media_mp3.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendAudio": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"audio\"; filename=\"audio.mpeg\"\r\nContent-Type: audio/mpeg\r\n\r\nID3\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "media_ogg",
			script: "scripts/media_ogg.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": { // consider ogg as document, it seems it's not fully supported
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document\"\r\nContent-Type: application/ogg\r\n\r\nOggS\x00\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "media_mp4",
			script: "scripts/media_mp4.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendVideo": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"video\"; filename=\"video.mp4\"\r\nContent-Type: video/mp4\r\n\r\n\x00\x00\x00\fftypmp4_\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
		{
			name:   "media_pdf",
			script: "scripts/media_pdf.sh",
			api: map[string][]apiserver.APIAct{
				"/botMORN/getUpdates": simpleUpdates,
				"/botMORN/sendDocument": {
					{
						IsJSON:   false,
						Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"chat_id\"\r\n\r\n1500\r\n--BOUND\r\nContent-Disposition: form-data; name=\"document\"; filename=\"document\"\r\nContent-Type: application/pdf\r\n\r\n%PDF-\r\n--BOUND--\r\n",
						Response: sendMessageResponseJSON,
					},
				},
			},
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tgURL, tgClose := apiserver.APIServer(t, cancel, cs.api)
			defer tgClose()

			bot := buildBot(tgURL)

			command := buildCommand(cs.script)

			err := xloop.Loop(ctx, bot, command)
			require.Error(t, err)
			require.Contains(t, err.Error(), "context canceled") // like "api: client: Post \"http://127.0.0.1:34241/botMORN/getUpdates\": context canceled"
		})
	}
}

func TestHttp(t *testing.T) {
	/* cases
	tg        bots loop
	|         |
	|         |<-- someone external calls bot over http
	|<--req---| (request to send)
	|---resp->|
	|         |--> reply to external client
	*/
	for _, cs := range []struct {
		name string
		curl []string
		qs   string
		api  map[string][]apiserver.APIAct
	}{
		{
			name: "curl_F", // curl -F works transparently as is
			curl: []string{"-q", "-s", "-F", "user_id=10", "-F", "text=ok"},
			qs:   "",
			api: map[string][]apiserver.APIAct{
				"/botMORN/someMethod": {{
					IsJSON:   false,
					Request:  "--BOUND\r\nContent-Disposition: form-data; name=\"user_id\"\r\n\r\n10\r\n--BOUND\r\nContent-Disposition: form-data; name=\"text\"\r\n\r\nok\r\n--BOUND--\r\n",
					Response: []byte("done."),
				}},
			},
		},
		{
			name: "curl_d",
			curl: []string{"-q", "-s", "-d", "ok"},
			qs:   "?to=111",
			api: map[string][]apiserver.APIAct{
				"/botMORN/sendMessage": {{
					IsJSON:   true,
					Request:  `{"chat_id":111, "text":"ok"}`,
					Response: []byte("done."),
				}},
			},
		},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			tgURL, tgClose := apiserver.APIServer(t, nil, cs.api)
			defer tgClose()

			bot := buildBot(tgURL)

			h := xctrl.Handler(bot, nil) // we won't use second argument in this test

			s := httptest.NewServer(h)

			ou, er := runCurl(t, append(cs.curl, s.URL+"/x/someMethod"+cs.qs)...)
			assert.Equal(t, "done.", ou)
			assert.Empty(t, er)
		})
	}
}

func TestHttp_long(t *testing.T) { // CAUTION: test has sleep
	/* cases
	tg        bots loop
	|         |
	|         |<-- someone external calls bot over http (method=RUN)
	|         |
	|         |--exec-->| long-running external script
	|         |<-stdout-|
	|         |
	|<--req---| (request to send)
	|---resp->| (response will be skipped; and test tries cover it by making small sleep)
	*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tgURL, tgClose := apiserver.APIServer(t, cancel, map[string][]apiserver.APIAct{
		"/botMORN/sendMessage": {{
			IsJSON:   true,
			Request:  `{"chat_id":222, "text":"args: a1 a2\n"}`,
			Response: nil, // response will be skipped, but in fact, we do not test this fact
		}},
	})
	defer tgClose()

	bot := buildBot(tgURL)
	command := buildCommand("scripts/longrunning.sh")

	h := xctrl.Handler(bot, command)

	s := httptest.NewServer(h)

	ou, er := runCurl(t, "-q", "-s", "-X", "RUN", s.URL+"/?to=222&a=a1&a=a2")
	assert.Empty(t, ou)
	assert.Empty(t, er)
	<-ctx.Done()
	time.Sleep(time.Millisecond * 100) // we give small amount of time to let Bot.API method finishing after receiving response; it is not necessary
}

func TestProc(t *testing.T) { // CAUTION: test has sleep indirectly
	ctx := context.Background()
	t.Run("argsEnvs", func(t *testing.T) {
		data, err := buildCommand("scripts/run_show_args.sh").Run(ctx, []string{"ARG1", "ARG2"}, []string{"test1=TEST1", "test2=TEST2"})
		require.NoError(t, err, "data="+string(data))
		assert.Equal(t, "arg1=ARG1 arg2=ARG2 test1=TEST1 test2=TEST2 TEST=test\n", string(data))
	})
	t.Run("exit", func(t *testing.T) {
		data, err := buildCommand("scripts/run_exit.sh").Run(ctx, nil, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "wait: exit status 28")
		assert.Nil(t, data)
	})
	t.Run("sigint", func(t *testing.T) {
		data, err := buildCommand("scripts/run_slow.sh").Run(ctx, nil, nil)
		require.NoError(t, err)
		assert.Equal(t,
			`start
trap SIGINT
trap ERR
end
trap EXIT
`, string(data))
	})
	t.Run("sigkill", func(t *testing.T) {
		data, err := buildCommand("scripts/run_immortal.sh").Run(ctx, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wait: signal: killed")
		assert.Nil(t, data)
	})
	t.Run("notfound", func(t *testing.T) {
		data, err := buildCommand("scripts/NOTFOUND").Run(ctx, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start: fork/exec scripts/NOTFOUND: no such file or directory")
		assert.Nil(t, data)
	})
}

func buildCommand(cmd string) *xproc.Cmd {
	return &xproc.Cmd{
		InterruptDelay: 200 * time.Millisecond, // timeouts important for TestProc
		KillDelay:      200 * time.Millisecond,
		Env:            []string{"TEST=test"},
		Command:        cmd,
		Cwd:            ".",
	}
}

func buildBot(origin string) *xbot.Bot {
	return &xbot.Bot{
		APIOrigin: origin,
		Token:     "MORN",
		Client:    http.DefaultClient,
	}
}

func runCurl(t *testing.T, args ...string) (string, string) {
	t.Helper()
	t.Logf("Run curl %s", strings.Join(args, " "))
	cmd := exec.Command("curl", args...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	err := cmd.Run()
	require.NoError(t, err)
	return stdOut.String(), stdErr.String()
}
