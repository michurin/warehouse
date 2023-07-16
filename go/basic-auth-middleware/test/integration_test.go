package test_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michurin/warehouse/go/basic-auth-middleware/httpauthmw"
	"github.com/michurin/warehouse/go/basic-auth-middleware/httpauthpasswd"
)

func TestAuthBasic(t *testing.T) { //nolint:funlen,gocognit,cyclop
	nakedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK: user=" + httpauthmw.UserName(r.Context())))
	})

	hmacKey := []byte{1, 2, 3, 4}
	user := "one"
	password := []byte("secret")
	passwd := map[string][]byte{
		user: hmac.New(sha256.New, hmacKey).Sum(password),
	}

	checker := httpauthpasswd.StaticAuth(passwd, hmacKey)

	mux := http.NewServeMux()
	mux.Handle("/test", httpauthmw.AuthBasic(nakedHandler, "test", checker))
	mux.Handle("/norealm", httpauthmw.AuthBasic(nakedHandler, "", checker))

	server := httptest.NewServer(mux)
	defer server.Close()

	for _, cs := range []struct {
		name         string
		path         string
		username     string
		password     string
		expectedCode int
	}{
		{"noauth", "/test", "", "", http.StatusUnauthorized},
		{"wrongpw", "/test", "one", "xxxxxx", http.StatusUnauthorized},
		{"wronguser", "/test", "xxx", "secret", http.StatusUnauthorized},
		{"authok", "/test", "one", "secret", http.StatusOK},
		{"noauthoknorealm", "/norealm", "", "", http.StatusUnauthorized},
		{"authoknorealm", "/norealm", "one", "secret", http.StatusOK},
	} {
		cs := cs
		t.Run(cs.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, server.URL+cs.path, nil) //nolint:noctx
			errPanic(err)
			req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(cs.username+":"+cs.password)))
			resp, err := http.DefaultClient.Do(req)
			errPanic(err)
			defer func() { errPanic(resp.Body.Close()) }()
			body, err := io.ReadAll(resp.Body)
			errPanic(err)
			wwwah := resp.Header.Get("Www-Authenticate")
			if cs.expectedCode == http.StatusOK { //nolint:nestif
				if resp.StatusCode != cs.expectedCode {
					t.Fail()
				}
				if wwwah != "" {
					t.Fail()
				}
				if string(body) != "OK: user=one" {
					t.Fail()
				}
			} else {
				if resp.StatusCode != cs.expectedCode {
					t.Fail()
				}
				expHdr := "Basic"
				if cs.path != "/norealm" {
					expHdr += ` realm="test", charset="UTF-8"`
				}
				if wwwah != expHdr {
					t.Fail()
				}
				if string(body) != "Unauthorized\n" {
					t.Fail()
				}
			}
			for k, v := range resp.Header {
				t.Log(k, v)
			}
			t.Log(resp.StatusCode)
			t.Log(resp.Header.Get("Www-Authenticate"))
			t.Log(string(body))
		})
	}
}

func TestRealmValidation(t *testing.T) {
	validCahrs := []byte(nil)
	nakedHandler := http.HandlerFunc(nil)
	for i := 0; i < 256; i++ {
		func() {
			defer func() {
				if err := recover(); err == nil {
					validCahrs = append(validCahrs, byte(i))
				}
			}()
			_ = httpauthmw.AuthBasic(nakedHandler, fmt.Sprintf("test_%c", i), nil)
		}()
	}
	if string(validCahrs) != "\x20!#$%&'()*+,-./0123456789:;=?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]_abcdefghijklmnopqrstuvwxyz~" {
		t.Fatalf("Invalid set of valid chars: %q", string(validCahrs))
	}
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
