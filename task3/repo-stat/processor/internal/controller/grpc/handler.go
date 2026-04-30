package grpc

import (
	"context"
	"errors"
	"log/slog"

	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	processorpb.UnimplementedProcessorServiceServer
	log       *slog.Logger
	forwardUC *usecase.ForwardRepository
}

func NewHandler(log *slog.Logger, forwardUC *usecase.ForwardRepository) *Server {
	return &Server{
		log:       log,
		forwardUC: forwardUC,
	}
}

func (s *Server) GetRepository(ctx context.Context, req *processorpb.GetRepositoryRequest) (*processorpb.Repository, error) {
	s.log.Info("processing repository request", "url", req.Url)

	repo, err := s.forwardUC.Execute(ctx, req.Url)
	if err != nil {
		s.log.Error("failed to get repository", "error", err)
		return nil, s.mapError(err)
	}

	if repo == nil {
		return nil, status.Error(codes.NotFound, "repository not found")
	}

	return s.toProto(repo), nil
}

func (s *Server) mapError(err error) error {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (s *Server) toProto(repo *domain.Repository) *processorpb.Repository {
	return &processorpb.Repository{
		FullName:    repo.FullName,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Language:    repo.Language,
	}
}
