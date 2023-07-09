package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"net/http"

	"github.com/michurin/warehouse/go/basic-auth-middleware/httpauthmw"
	"github.com/michurin/warehouse/go/basic-auth-middleware/httpauthpasswd"
)

func main() {
	const addr = "localhost:9999"
	const user = "u"
	const passwd = "1234"
	const realm = "x"
	const key = "KeY"

	log.Printf("Starting server at http://%s, try to login with username `%s` and password `%s`", addr, user, passwd)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK: " + httpauthmw.UserName(r.Context())))
	})

	hkey := []byte(key)
	checker := httpauthpasswd.StaticAuth(map[string][]byte{
		user: hmac.New(sha256.New, hkey).Sum([]byte(passwd)),
	}, hkey)

	handler = httpauthmw.AuthBasic(handler, realm, checker)

	_ = http.ListenAndServe(addr, handler) //nolint:gosec // no timeouts
}
