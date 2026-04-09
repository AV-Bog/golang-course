package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/collector"
	"repo-stat/processor/internal/controller/grpc"
	"repo-stat/processor/internal/usecase"
	processorpb "repo-stat/proto/proto/processor"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting processor server...")

	collectorClient, err := collector.NewClient(cfg.Services.Collector, log)
	if err != nil {
		return fmt.Errorf("create collector client: %w", err)
	}
	defer func(collectorClient *collector.Client) {
		err := collectorClient.Close()
		if err != nil {

		}
	}(collectorClient)

	forwardUC := usecase.NewForwardRepository(collectorClient)
	handler := grpc.NewHandler(log, forwardUC)

	srv, err := grpcserver.New(cfg.GRPC.Address)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}

	processorpb.RegisterProcessorServiceServer(srv.GRPC(), handler)

	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run grpc server: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	var _, cancel = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := run(ctx); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
