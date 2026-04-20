package collector

import (
	"context"
	"time"

	"repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	pb   collectorpb.CollectorServiceClient
}

type Config struct {
	Address string
	Timeout time.Duration
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
		pb:   collectorpb.NewCollectorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, url string) (*domain.Repository, error) {
	resp, err := c.pb.GetRepository(ctx, &collectorpb.GetRepositoryRequest{Url: url})
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
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   createdAt,
		Language:    resp.Language,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
