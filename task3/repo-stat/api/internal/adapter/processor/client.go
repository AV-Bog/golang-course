package processor

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	processorpb "repo-stat/proto/proto/processor"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorpb.ProcessorServiceClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   processorpb.NewProcessorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, url string) (*processorpb.Repository, error) {
	return c.pb.GetRepository(ctx, &processorpb.GetRepositoryRequest{Url: url})
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
