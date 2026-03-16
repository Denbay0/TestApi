package grpcclient

import (
	"context"

	"google.golang.org/grpc"
)

type ReportClient struct {
	conn *grpc.ClientConn
}

func (c *ReportClient) Summary(ctx context.Context, requestID string) (map[string]any, error) {
	_ = c
	_ = withRequestID(ctx, requestID)
	return map[string]any{"todo": "wire Summary RPC"}, nil
}

func (c *ReportClient) ByCategory(ctx context.Context, requestID string) (map[string]any, error) {
	_ = c
	_ = withRequestID(ctx, requestID)
	return map[string]any{"items": []any{}, "todo": "wire ByCategory RPC"}, nil
}
