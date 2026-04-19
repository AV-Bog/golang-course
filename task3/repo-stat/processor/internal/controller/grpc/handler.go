package grpc

import (
	"context"
	"log/slog"
	"repo-stat/processor/internal/adapter/collector"
	processorpb "repo-stat/proto/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	processorpb.UnimplementedProcessorServiceServer
	log    *slog.Logger
	client *collector.Client
}

func NewHandler(log *slog.Logger, client *collector.Client) *Server {
	return &Server{
		log:    log,
		client: client,
	}
}

func (s *Server) GetRepository(ctx context.Context, req *processorpb.GetRepositoryRequest) (*processorpb.Repository, error) {
	repo, err := s.client.GetRepository(ctx, req.Url)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if repo == nil {
		return nil, status.Error(codes.NotFound, "repository not found")
	}

	return &processorpb.Repository{
		FullName:    repo.FullName,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
		Language:    repo.Language,
	}, nil
}
