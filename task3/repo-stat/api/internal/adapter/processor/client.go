package processor

import (
	"context"
	"time"

	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	pb   processorpb.ProcessorServiceClient
}

type Config struct {
	Address string
}

func NewClient(cfg Config) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		pb:   processorpb.NewProcessorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, url string) (*domain.Repository, error) {
	resp, err := c.pb.GetRepository(ctx, &processorpb.GetRepositoryRequest{Url: url})
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		createdAt = time.Time{}
	}

	return &domain.Repository{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   createdAt,
		Language:    resp.Language,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.pb.Ping(ctx, &processorpb.PingRequest{})
	return err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
