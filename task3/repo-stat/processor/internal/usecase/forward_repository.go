package usecase

import (
	"context"

	"repo-stat/processor/internal/domain"
)

type CollectorClient interface {
	GetRepository(ctx context.Context, url string) (*domain.Repository, error)
}

type ForwardRepository struct {
	collectorClient CollectorClient
}

func NewForwardRepository(client CollectorClient) *ForwardRepository {
	return &ForwardRepository{
		collectorClient: client,
	}
}

func (uc *ForwardRepository) Execute(ctx context.Context, url string) (*domain.Repository, error) {

	return uc.collectorClient.GetRepository(ctx, url)
}
