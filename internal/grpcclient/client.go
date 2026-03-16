package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ClientSet struct {
	Identity     *IdentityClient
	EventCommand *EventCommandClient
	EventQuery   *EventQueryClient
	Report       *ReportClient
	conns        []*grpc.ClientConn
}

type dialConfig struct {
	Identity     string
	EventCommand string
	EventQuery   string
	Report       string
}

func New(ctx context.Context, cfg dialConfig) (*ClientSet, error) {
	identityConn, err := dial(ctx, cfg.Identity)
	if err != nil {
		return nil, fmt.Errorf("dial identity: %w", err)
	}
	eventCommandConn, err := dial(ctx, cfg.EventCommand)
	if err != nil {
		_ = identityConn.Close()
		return nil, fmt.Errorf("dial event-command: %w", err)
	}
	eventQueryConn, err := dial(ctx, cfg.EventQuery)
	if err != nil {
		_ = identityConn.Close()
		_ = eventCommandConn.Close()
		return nil, fmt.Errorf("dial event-query: %w", err)
	}
	reportConn, err := dial(ctx, cfg.Report)
	if err != nil {
		_ = identityConn.Close()
		_ = eventCommandConn.Close()
		_ = eventQueryConn.Close()
		return nil, fmt.Errorf("dial report: %w", err)
	}

	return &ClientSet{
		Identity:     &IdentityClient{conn: identityConn},
		EventCommand: &EventCommandClient{conn: eventCommandConn},
		EventQuery:   &EventQueryClient{conn: eventQueryConn},
		Report:       &ReportClient{conn: reportConn},
		conns:        []*grpc.ClientConn{identityConn, eventCommandConn, eventQueryConn, reportConn},
	}, nil
}

func (c *ClientSet) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}

func dial(ctx context.Context, rawTarget string) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, NormalizeTarget(rawTarget), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func withRequestID(ctx context.Context, requestID string) context.Context {
	if requestID == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "x-request-id", requestID)
}
