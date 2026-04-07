package usecase

import (
	"context"
)

type RepositoryClient interface {
	GetRepository(ctx context.Context, url string) (interface{}, error) // замените на конкретный тип
}

type GetRepository struct {
	processorClient RepositoryClient
}

func NewGetRepository(processorClient RepositoryClient) *GetRepository {
	return &GetRepository{
		processorClient: processorClient,
	}
}

func (u *GetRepository) Execute(ctx context.Context, url string) (interface{}, error) {
	return u.processorClient.GetRepository(ctx, url)
}
