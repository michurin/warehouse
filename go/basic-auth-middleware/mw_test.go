package basicauthmiddleware_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"

	basicauthmiddleware "github.com/michurin/warehouse/go/basic-auth-middleware"
)

func ExampleBasicAuth() {
	hmacKey := []byte{1, 2, 3, 4}
	user := "one"
	password := []byte("secret")
	passwd := map[string][]byte{
		user: hmac.New(sha256.New, hmacKey).Sum(password),
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK: user=" + basicauthmiddleware.UserName(r.Context())))
	})

	h = basicauthmiddleware.BasicAuth(h, passwd, hmacKey)

	s := httptest.NewServer(h)
	defer s.Close()

	curl("-qs", s.URL)
	curl("-qs", s.URL, "-u", "on1:secret")
	curl("-qs", s.URL, "-u", "one:secre1")
	curl("-qs", s.URL, "-u", "one:secret")

	// output:
	// Unauthorized
	// Unauthorized
	// Unauthorized
	// OK: user=one
}

func curl(args ...string) {
	cmd := exec.Command("curl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
