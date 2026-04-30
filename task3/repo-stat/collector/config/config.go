package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/grpcserver"
	"repo-stat/platform/logger"
)

type GitHub struct {
	Token string `yaml:"token" env:"GITHUB_TOKEN"`
}

type Config struct {
	App    App               `yaml:"app"`
	GRPC   grpcserver.Config `yaml:"grpc"`
	Logger logger.Config     `yaml:"logger"`
	GitHub GitHub            `yaml:"github"`
}

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-collector"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
