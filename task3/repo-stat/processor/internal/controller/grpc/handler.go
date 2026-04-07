package grpc

import (
	"context"
	"log/slog"

	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/gen/go/processor/v1"
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
		return &processorpb.GetRepositoryResponse{
			ErrorCode:    int32(500),
			ErrorMessage: err.Error(),
		}, err
	}

	return &processorpb.GetRepositoryResponse{
		Repository:   resp.Repository,
		ErrorCode:    resp.ErrorCode,
		ErrorMessage: resp.ErrorMessage,
	}, nil
}

func (s *Server) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	s.log.Debug("Ping received")
	return &processorpb.PingResponse{Status: "pong"}, nil
}
