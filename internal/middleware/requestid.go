package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func InjectRequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := chimiddleware.GetReqID(r.Context())
			r = r.WithContext(WithRequestID(r.Context(), reqID))
			next.ServeHTTP(w, r)
		})
	}
}
