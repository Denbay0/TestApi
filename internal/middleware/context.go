package middleware

import "context"

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	authTokenKey contextKey = "auth_token"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

func WithAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenKey, token)
}

func AuthTokenFromContext(ctx context.Context) string {
	v, _ := ctx.Value(authTokenKey).(string)
	return v
}
