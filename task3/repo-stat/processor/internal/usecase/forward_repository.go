package usecase

import (
	"context"

	collectorpb "repo-stat/proto/proto/collector"
)

type CollectorClient interface {
	GetRepository(ctx context.Context, url string) (*collectorpb.GetRepositoryResponse, error)
}

type ForwardRepository struct {
	collectorClient CollectorClient
}

func NewForwardRepository(collectorClient CollectorClient) *ForwardRepository {
	return &ForwardRepository{
		collectorClient: collectorClient,
	}
}

func (uc *ForwardRepository) Execute(ctx context.Context, url string) (*collectorpb.GetRepositoryResponse, error) {
	return uc.collectorClient.GetRepository(ctx, url)
}
