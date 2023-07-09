package httpauthpasswd

import (
	"crypto/hmac"
	"crypto/sha256"

	"github.com/michurin/warehouse/go/basic-auth-middleware/httpauthmw"
)

func StaticAuth(passwd map[string][]byte, hmacKey []byte) httpauthmw.AuthChecker {
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
