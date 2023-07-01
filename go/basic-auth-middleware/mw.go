package basicauthmiddleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
)

type userNameKeyType int

const userNameKey userNameKeyType = iota

func BasicAuth(next http.HandlerFunc, passwd map[string][]byte, hmacKey []byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			if pwHmac, ok := passwd[username]; ok {
				if hmac.Equal(hmac.New(sha256.New, hmacKey).Sum([]byte(password)), pwHmac) {
					next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userNameKey, username)))
					return
				}
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func UserName(ctx context.Context) string {
	s, _ := ctx.Value(userNameKey).(string)
	return s
}
