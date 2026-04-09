package subscriber

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	log        *slog.Logger
	httpClient *http.Client
	baseURL    string
}

func NewHTTPClient(address string, log *slog.Logger) (*HTTPClient, error) {
	if !strings.HasPrefix(address, "https://") && !strings.HasPrefix(address, "https://") {
		address = "https://" + address
	}

	return &HTTPClient{
		log: log,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
		baseURL: address,
	}, nil
}

func (c *HTTPClient) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/ping", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

func (c *HTTPClient) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}
