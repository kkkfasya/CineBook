package middleware

import (
	"crypto/subtle"
	"errors"
	"net/http"

	"github.com/kkkfasya/CineBook/internal/utils"
)

var errUnauthorized = errors.New("unauthorized access")

func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rusn, rpwd, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				utils.WriteError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}
			usnMatch := subtle.ConstantTimeCompare([]byte(rusn), []byte(username)) == 1
			pwdMatch := subtle.ConstantTimeCompare([]byte(rpwd), []byte(password)) == 1

			if !usnMatch || !pwdMatch {
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				utils.WriteError(w, http.StatusUnauthorized, errUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

}
