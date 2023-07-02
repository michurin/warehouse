package basicauthmiddleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
)

type userNameKeyType int

const userNameKey userNameKeyType = iota

type AuthChecker func(username, password string) bool

func StaticAuth(passwd map[string][]byte, hmacKey []byte) AuthChecker {
	return func(username, password string) bool {
		pwHmac, ok := passwd[username]
		if !ok {
			return false
		}
		if hmac.Equal(hmac.New(sha256.New, hmacKey).Sum([]byte(password)), pwHmac) {
			return true
		}
		return false
	}
}

// valid chars
// !#$%&'()*+,-./0123456789:;=?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]_abcdefghijklmnopqrstuvwxyz~ .
var realmCharMap = [4]uint64{ //nolint:gochecknoglobals
	0b_101011111_1111111_11111111_11111010_00000000_00000000_00000000_00000000,
	0b_01000111_11111111_11111111_11111110_10101111_11111111_11111111_11111111,
	0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
	0b_00000000_00000000_00000000_00000000_00000000_00000000_00000000_00000000,
}

func AuthBasic(next http.HandlerFunc, realm string, checker AuthChecker) http.HandlerFunc {
	ah := "Basic"
	if realm != "" {
		for _, c := range realm {
			if (realmCharMap[c/64] & (uint64(1) << (c % 64))) == 0 { //nolint:gomnd
				panic(fmt.Sprintf("Invalid char in realm: %q", c))
			}
		}
		ah += ` realm="` + realm + `", charset="UTF-8"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok && checker(username, password) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userNameKey, username)))
			return
		}
		w.Header().Set("WWW-Authenticate", ah)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func UserName(ctx context.Context) string {
	s, _ := ctx.Value(userNameKey).(string)
	return s
}
