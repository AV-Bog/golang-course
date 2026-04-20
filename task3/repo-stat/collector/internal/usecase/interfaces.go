package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type GitHubClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}
