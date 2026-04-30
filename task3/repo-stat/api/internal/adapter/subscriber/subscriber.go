package subscriber

import (
	"context"
	"log/slog"

	subscirberpb "repo-stat/proto/subscriber"

	"google.golang.org/grpc"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   subscirberpb.SubscriberClient
}

func (c *Client) Ping(ctx context.Context) error {
	c.log.Info("Calling subscriber Ping")
	_, err := c.pb.Ping(ctx, &subscirberpb.PingRequest{})
	if err != nil {
		c.log.Error("subscriber ping failed", "error", err)
		return err
	}
	c.log.Info("subscriber ping succeeded")
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
