package middleware

import (
	"crypto/subtle"
	"errors"
	"net/http"

	"github.com/kkkfasya/CineBook/internal/utils"
)

func BasicAuth(username, password string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rusrn, rpwd, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}

			userMatch := subtle.ConstantTimeCompare([]byte(rusrn), []byte(username)) == 1
			passwordMatch := subtle.ConstantTimeCompare([]byte(rpwd), []byte(password)) == 1

			if !userMatch || !passwordMatch {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				utils.WriteError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}

}
