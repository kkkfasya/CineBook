package middleware

import (
	"net/http"
)

func StripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if len(p) > 1 && p[len(p)-1] == '/' {
			http.Redirect(w, r, p[:len(p)-1], http.StatusPermanentRedirect)
			return
		}

		next.ServeHTTP(w, r)
	})
}
