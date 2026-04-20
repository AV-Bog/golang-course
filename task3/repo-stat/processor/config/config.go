package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type Config struct {
	App      App               `yaml:"app"`
	GRPC     grpcserver.Config `yaml:"grpc"`
	Logger   logger.Config     `yaml:"logger"`
	Services Services          `yaml:"services"`
}

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-processor"`
}

type Services struct {
	Collector string `yaml:"collector" env:"COLLECTOR_ADDRESS" env-default:"localhost:50052"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
