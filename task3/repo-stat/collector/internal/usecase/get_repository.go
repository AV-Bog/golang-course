package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/domain"
)

var (
	ErrInvalidURL   = errors.New("invalid repository URL")
	ErrRepoNotFound = errors.New("repository not found")
)

type GetRepository struct {
	githubClient *github.Client
}

func NewGetRepository(githubToken string) *GetRepository {
	config := github.Config{
		AuthToken: githubToken,
		Timeout:   10 * time.Second,
	}
	return &GetRepository{
		githubClient: github.NewClient(config),
	}
}

func (uc *GetRepository) Execute(ctx context.Context, repoURL string) (*domain.Repository, error) {
	owner, name, err := parseGitHubURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	repo, err := uc.githubClient.GetRepository(owner, name)
	if err != nil {
		if errors.Is(err, github.ErrRepoNotFound) {
			return nil, ErrRepoNotFound
		}
		return nil, err
	}

	return repo, nil
}

func parseGitHubURL(repoURL string) (owner, repo string, err error) {
	repoURL = strings.TrimSpace(repoURL)
	repoURL = strings.TrimSuffix(repoURL, ".git")

	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "https://" + repoURL
	}

	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	if parsedURL.Host != "github.com" {
		return "", "", fmt.Errorf("not a GitHub URL")
	}

	path := strings.Trim(parsedURL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid format: expected /owner/repo")
	}

	return parts[0], parts[1], nil
}
