package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env        string `env:"APP_ENV" envDefault:"local-local"`
	ServerAddr string `env:"SERVER_ADDR" envDefault:"127.0.0.1"`
	ServerPort string `env:"SERVER_ADDR" envDefault:":8082"`
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}
	return cfg, nil
}
