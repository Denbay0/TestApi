package middleware

import (
	"net/http"

	"github.com/example/edge-api/internal/auth"
)

func ParseAuthCookie() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(auth.AuthCookieName)
			if err == nil && cookie.Value != "" {
				r = r.WithContext(WithAuthToken(r.Context(), cookie.Value))
			}
			next.ServeHTTP(w, r)
		})
	}
}
