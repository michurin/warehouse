package basicauthmiddleware_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	basicauthmiddleware "github.com/michurin/warehouse/go/basic-auth-middleware"
)

func ExampleAuthBasic() {
	nakedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK: user=" + basicauthmiddleware.UserName(r.Context())))
	})

	hmacKey := []byte{1, 2, 3, 4}
	user := "one"
	password := []byte("secret")
	passwd := map[string][]byte{
		user: hmac.New(sha256.New, hmacKey).Sum(password),
	}

	checker := basicauthmiddleware.StaticAuth(passwd, hmacKey)

	wrappedHandler := basicauthmiddleware.AuthBasic(nakedHandler, "test", checker)

	server := httptest.NewServer(wrappedHandler)
	defer server.Close()

	url := server.URL
	curl := func(args ...string) {
		fmt.Println("curl " + strings.Join(append(args, "http://testserver/"), " "))
		cmd := exec.Command("curl", append([]string{"-qs"}, append(args, url)...)...) //nolint:gosec // -q must be first arg
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	curl()
	curl("-u", "xxx:secret")
	curl("-u", "one:xxxxxx")
	curl("-u", "one:secret")

	// output:
	// curl http://testserver/
	// Unauthorized
	// curl -u xxx:secret http://testserver/
	// Unauthorized
	// curl -u one:xxxxxx http://testserver/
	// Unauthorized
	// curl -u one:secret http://testserver/
	// OK: user=one
}

func TestAuthBasic(t *testing.T) { //nolint:funlen,gocognit,cyclop
	nakedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK: user=" + basicauthmiddleware.UserName(r.Context())))
	})

	hmacKey := []byte{1, 2, 3, 4}
	user := "one"
	password := []byte("secret")
	passwd := map[string][]byte{
		user: hmac.New(sha256.New, hmacKey).Sum(password),
	}

	checker := basicauthmiddleware.StaticAuth(passwd, hmacKey)

	mux := http.NewServeMux()
	mux.Handle("/test", basicauthmiddleware.AuthBasic(nakedHandler, "test", checker))
	mux.Handle("/norealm", basicauthmiddleware.AuthBasic(nakedHandler, "", checker))

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
			_ = basicauthmiddleware.AuthBasic(nakedHandler, fmt.Sprintf("test_%c", i), nil)
		}()
	}
	if string(validCahrs) != "!#$%&'()*+,-./0123456789:;=?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]_abcdefghijklmnopqrstuvwxyz~" {
		t.Fail()
	}
}

func errPanic(err error) {
	if err != nil {
		panic(err)
	}
}
