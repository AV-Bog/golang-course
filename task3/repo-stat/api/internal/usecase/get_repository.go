package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type RepositoryClient interface {
	GetRepository(ctx context.Context, url string) (*domain.Repository, error)
}
type GetRepository struct {
	processorClient RepositoryClient
}

func NewGetRepository(processorClient RepositoryClient) *GetRepository {
	return &GetRepository{
		processorClient: processorClient,
	}
}

func (u *GetRepository) Execute(ctx context.Context, url string) (*domain.Repository, error) {
	return u.processorClient.GetRepository(ctx, url)
}
