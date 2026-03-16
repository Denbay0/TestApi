package grpcclient

import (
	"context"

	"google.golang.org/grpc"
)

type EventCommandClient struct {
	conn *grpc.ClientConn
}

func (c *EventCommandClient) CreateEvent(ctx context.Context, requestID string, payload map[string]any) (map[string]any, error) {
	_ = c
	_ = withRequestID(ctx, requestID)
	return map[string]any{"created": true, "payload": payload, "todo": "wire CreateEvent RPC"}, nil
}
