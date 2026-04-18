package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"sse/handler"
	"sse/loggingmw"
	"sse/room"
)

func noerr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assert[T comparable](t *testing.T, x T, y T) {
	t.Helper()
	// s, ok := any(x).(string) // TODO compare JSON?
	// if ok {
	// 	t.Log("STR", s)
	// } else {
	// 	t.Logf("%T", x)
	// }
	if x != y {
		t.Fatalf("NOT EQUAL:\n%#v\n%#v", x, y)
	}
}

func allignTime() {
	time.Sleep(181482 * time.Hour) // sleep up to September 13, 2020, 18:00:00; time_stamp=1600020000
}

func do(t *testing.T, mx http.Handler, method, path, body, expectedBody string) {
	t.Helper()
	ctx := t.Context()
	b := io.Reader(http.NoBody)
	if method == http.MethodPost {
		b = strings.NewReader(body)
	}
	r, err := http.NewRequestWithContext(ctx, method, "http://x"+path, b)
	noerr(t, err)
	r.RemoteAddr = "127.0.0.1"

	w := httptest.NewRecorder()

	mx.ServeHTTP(w, r)

	assert(t, w.Code, http.StatusOK)
	assert(t, w.Body.String(), expectedBody)
}

/*
A try to /fetch (failed)
A try to /pub (failed)
A try to /lock (failed)
A /enter
A /fetch
A /pub
A /fetch+
B /enter
A /fetch+
B /pub
A /fetch+
sleep
A /pub
sleep
B =die=
A /fetch
A /lock
B /enter (failed)
A /unlock
B /enter (ok)
A /fetch+
B /fetch (all)
*/

func TestHandler_complexFlow(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		allignTime()

		house := room.New()

		// TODO run audit (revision loop)

		mx := loggingmw.MW(handler.Handler(house))

		// enter

		do(t, mx, http.MethodPost, "/bin/enter",
			`{"room":"R","user":"id-A","name":"Alex","color":"#111111"}`,
			`{"users":[{"name":"Alex","color":"#111111"}],"locked":false}`)

		// lock

		do(t, mx, http.MethodPost, "/bin/lock", `{"room":"R","user":"id-A","lock":true}`, ``)

		// fetch (TODO)

		do(t, mx, http.MethodGet, "/bin/fetch?room=R&user=id-A", ``,
			`event: message
retry: 200
id: 1600020000000000002
data: {"message":{"color":"#990099","message":"Alex touched LOCK","name":"#ROBOT","ts":1600020000000},"users":[{"name":"Alex","color":"#111111"}],"locked":true}
data: {"message":{"color":"#990099","message":"Alex HERE!","name":"#ROBOT","ts":1600020000000},"users":[{"name":"Alex","color":"#111111"}],"locked":false}

`)
	})
}
