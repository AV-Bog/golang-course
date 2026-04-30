package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/controller/http"
	"repo-stat/api/internal/usecase"
	"repo-stat/platform/httpserver"
	"repo-stat/platform/logger"
	"syscall"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting api gateway...")

	// gRPC клиент к Processor
	processorClient, err := processor.NewClient(processor.Config{
		Address: cfg.Services.Processor,
	})
	if err != nil {
		return fmt.Errorf("create processor client: %w", err)
	}
	defer func() {
		if err := processorClient.Close(); err != nil {
			log.Error("failed to close processor client", "error", err)
		}
	}()

	// usecase
	pingUC := usecase.NewPing(processorClient, nil) // если subscriber пока не нужен
	getRepoUC := usecase.NewGetRepository(processorClient)

	// HTTP хендлер с usecase
	handler := http.NewHandler(pingUC, getRepoUC) // ← NewHandler (без 's')

	// HTTP сервер
	srv := httpserver.New(cfg.HTTP, handler)
	if err := srv.Run(ctx); err != nil {
		return fmt.Errorf("run http server: %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var exitCode int
	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
