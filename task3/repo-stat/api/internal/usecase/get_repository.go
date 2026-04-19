package usecase

import (
	"context"
	processorpb "repo-stat/proto/proto/processor"
)

type RepositoryClient interface {
	GetRepository(ctx context.Context, url string) (*processorpb.Repository, error)
}
type GetRepository struct {
	processorClient RepositoryClient
}

func NewGetRepository(processorClient RepositoryClient) *GetRepository {
	return &GetRepository{
		processorClient: processorClient,
	}
}

func (u *GetRepository) Execute(ctx context.Context, url string) (*processorpb.Repository, error) {
	return u.processorClient.GetRepository(ctx, url)
}
