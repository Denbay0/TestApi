package grpcclient

import (
	"context"

	"google.golang.org/grpc"
)

type IdentityClient struct {
	conn *grpc.ClientConn
}

type RegisterParams struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Session struct {
	Token string `json:"token"`
}

func (c *IdentityClient) Register(ctx context.Context, requestID string, _ RegisterParams) (*Session, error) {
	_ = withRequestID(ctx, requestID)
	return &Session{Token: "todo-register-token"}, nil
}

func (c *IdentityClient) Login(ctx context.Context, requestID string, _ LoginParams) (*Session, error) {
	_ = withRequestID(ctx, requestID)
	return &Session{Token: "todo-login-token"}, nil
}

func (c *IdentityClient) Me(ctx context.Context, requestID string) (map[string]any, error) {
	_ = withRequestID(ctx, requestID)
	return map[string]any{"status": "TODO: wire identity Me rpc"}, nil
}
