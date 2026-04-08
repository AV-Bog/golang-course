package grpc

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/proto/processor"
)

type Server struct {
	processorpb.UnimplementedProcessorServiceServer
	log               *slog.Logger
	forwardRepository *usecase.ForwardRepository
}

func NewHandler(log *slog.Logger, forwardRepository *usecase.ForwardRepository) *Server {
	return &Server{
		log:               log,
		forwardRepository: forwardRepository,
	}
}

func (s *Server) GetRepository(ctx context.Context, req *processorpb.GetRepositoryRequest) (*processorpb.GetRepositoryResponse, error) {
	resp, err := s.forwardRepository.Execute(ctx, req.Url)
	if err != nil {
		s.log.Error("failed to forward request", "error", err, "url", req.Url)
		return nil, err
	}

	return &processorpb.GetRepositoryResponse{
		Repository: &processorpb.Repository{
			FullName:    resp.Repository.FullName,
			Description: resp.Repository.Description,
			Stars:       resp.Repository.Stars,
			Forks:       resp.Repository.Forks,
			CreatedAt:   resp.Repository.CreatedAt,
			Language:    resp.Repository.Language,
		},
	}, nil
}

func (s *Server) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	s.log.Debug("Ping received")
	return &processorpb.PingResponse{Status: "pong"}, nil
}
