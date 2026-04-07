package grpc

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"repo-stat/collector/internal/usecase"
	collectorpb "repo-stat/proto/gen/go/collector/v1"
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

func (s *Server) GetRepository(ctx context.Context, req *collectorpb.GetRepositoryRequest) (*collectorpb.GetRepositoryResponse, error) {
	// Ваша реализация из Task2
	repo, err := s.getRepositoryUC.Execute(ctx, req.Url)
	if err != nil {
		return s.mapError(err)
	}

	return &collectorpb.GetRepositoryResponse{
		Repository: &collectorpb.Repository{
			FullName:    repo.FullName,
			Description: repo.Description,
			Stars:       int32(repo.Stars),
			Forks:       int64(repo.Forks),
			CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Language:    repo.Language,
		},
		ErrorCode:    0,
		ErrorMessage: "",
	}, nil
}

func (s *Server) mapError(err error) (*collectorpb.GetRepositoryResponse, error) {
	// Ваша реализация маппинга ошибок из Task2
	var code codes.Code
	var message string

	switch {
	case err == usecase.ErrInvalidURL:
		code = codes.InvalidArgument
		message = "Invalid repository URL"
	case err == usecase.ErrRepoNotFound:
		code = codes.NotFound
		message = "Repository not found"
	default:
		code = codes.Internal
		message = "Internal server error"
	}

	return &collectorpb.GetRepositoryResponse{
		ErrorCode:    int32(code),
		ErrorMessage: message,
	}, status.Error(code, message)
}
