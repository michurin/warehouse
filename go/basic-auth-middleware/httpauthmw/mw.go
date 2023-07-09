package httpauthmw

import (
	"context"
	"net/http"
)

func AuthBasic(next http.HandlerFunc, realm string, checker AuthChecker) http.HandlerFunc {
	ah := "Basic"
	if realm != "" {
		validateRialm(realm)
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
