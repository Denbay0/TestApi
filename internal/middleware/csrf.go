package middleware

import (
	"net/http"
	"strings"

	"github.com/example/edge-api/internal/auth"
	"github.com/example/edge-api/internal/response"
)

func CSRFMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isMutatingMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			tokenHeader := strings.TrimSpace(r.Header.Get("X-CSRF-Token"))
			cookie, err := r.Cookie(auth.CSRFCookieName)
			if err != nil || cookie.Value == "" || tokenHeader == "" || cookie.Value != tokenHeader {
				response.Error(w, http.StatusForbidden, "csrf_mismatch", "invalid csrf token", RequestIDFromContext(r.Context()))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isMutatingMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}
