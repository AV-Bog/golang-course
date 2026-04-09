package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/controller/grpc"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	collectorpb "repo-stat/proto/proto/collector"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting collector server...")

	getRepoUC := usecase.NewGetRepository(cfg.GitHub.Token)
	handler := grpc.NewHandler(log, getRepoUC)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	collectorpb.RegisterCollectorServiceServer(srv.GRPC(), handler)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
}
