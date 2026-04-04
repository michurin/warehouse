package handler_test

import (
	"context"
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

func httpDo(
	t *testing.T,
	ctx context.Context,
	mx http.Handler,
	method,
	path,
	body,
	expetedBody string) {
	t.Helper()
	w := httptest.NewRecorder()
	b := io.Reader(http.NoBody)
	if body != "" {
		b = strings.NewReader(body)
	}
	r, err := http.NewRequestWithContext(ctx, method, "http://x"+path, b)
	noerr(t, err)
	mx.ServeHTTP(w, r)
	assert(t, w.Code, http.StatusOK)
	assert(t, expetedBody, w.Body.String())
	t.Log("header:", w.Header()) // TODO assert
}

func TestMux(t *testing.T) {
	t.Run("static", func(t *testing.T) {
		house := room.New()
		mx := handler.Handler(house)
		t.Run("favicon", func(t *testing.T) {
			ctx := t.Context()
			w := httptest.NewRecorder()
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://x/favicon.ico", nil)
			noerr(t, err)
			mx.ServeHTTP(w, r)
			assert(t, w.Code, http.StatusOK)
		})
		// TODO css, js, index, /XX
	})
	t.Run("messaging", func(t *testing.T) {
		t.Run("just_empty_clinet_timeout", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				house := room.New()
				mx := handler.Handler(house)
				ctx := t.Context()
				ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
				defer cancel()
				w := httptest.NewRecorder()
				r, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://x/fetch", nil)
				noerr(t, err)
				mx.ServeHTTP(w, r)
				assert(t, w.Code, http.StatusOK)
				t.Log("body:", w.Body.String()) // TODO assert
				t.Log("header:", w.Header())    // TODO assert
			})
		})
		t.Run("just_empty_server_timeout", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				house := room.New()
				mx := handler.Handler(house)
				ctx := t.Context()
				ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
				defer cancel()
				w := httptest.NewRecorder()
				r, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://x/fetch", nil)
				noerr(t, err)
				mx.ServeHTTP(w, r)
				assert(t, w.Code, http.StatusOK)
				t.Log("body:", w.Body.String()) // TODO assert
				t.Log("header:", w.Header())    // TODO assert
			})
		})
		t.Run("pub_sub", func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				house := room.New()
				mx := loggingmw.MW(handler.Handler(house))
				ctx := t.Context()
				httpDo(t, ctx, mx, http.MethodGet, "/fetch?room=CHAN&user=xx", "", "") // it will take 28s
				httpDo(t, ctx, mx, http.MethodPost, "/pub",
					`{"room":"CHAN","user":"xx","color":"#fff","name":"M","message":"text"}`,
					``)
				httpDo(t, ctx, mx, http.MethodGet, "/fetch?room=CHAN&user=xx", "",
					`event: message
retry: 200
id: 946684800000000001
data: {"message":{"color":"#fff","message":"text","name":"M","ts":946684828000}}

`)
			})
		})
	})
}
