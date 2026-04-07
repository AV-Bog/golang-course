package collector

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	collectorpb "repo-stat/proto/gen/go/collector/v1"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   collectorpb.CollectorServiceClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		pb:   collectorpb.NewCollectorServiceClient(conn),
	}, nil
}

func (c *Client) GetRepository(ctx context.Context, url string) (*collectorpb.GetRepositoryResponse, error) {
	return c.pb.GetRepository(ctx, &collectorpb.GetRepositoryRequest{Url: url})
}

func (c *Client) Close() error {
	return c.conn.Close()
}
