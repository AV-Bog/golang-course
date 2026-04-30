package grpc

import (
	"context"
	"errors"
	"log/slog"
	"repo-stat/collector/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	collectorpb "repo-stat/proto/proto/collector"
)

type Server struct {
	collectorpb.UnimplementedCollectorServiceServer
	log             *slog.Logger
	getRepositoryUC *usecase.GetRepository
}

func NewHandler(log *slog.Logger, getRepositoryUC *usecase.GetRepository) *Server {
	return &Server{
		log:             log,
		getRepositoryUC: getRepositoryUC,
	}
}

func (s *Server) GetRepository(ctx context.Context, req *collectorpb.GetRepositoryRequest) (*collectorpb.Repository, error) {
	s.log.Info("getting repository", "url", req.Url)
	repo, err := s.getRepositoryUC.Execute(ctx, req.Url)
	if err != nil {
		s.log.Error("failed to get repository", "error", err)
		return nil, s.mapError(err)
	}
	if repo == nil {
		return nil, status.Error(codes.Internal, "repository is nil")
	}

	return &collectorpb.Repository{
		FullName:    repo.FullName,
		Stars:       int32(repo.Stars),
		Description: repo.Description,
		Forks:       int64(repo.Forks),
		CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Language:    repo.Language,
	}, nil
}

func (s *Server) mapError(err error) error {
	var code codes.Code

	switch {
	case errors.Is(err, usecase.ErrInvalidURL):
		code = codes.InvalidArgument
	case errors.Is(err, usecase.ErrRepoNotFound):
		code = codes.NotFound
	default:
		code = codes.Internal
	}

	return status.Error(code, err.Error())
}
