package config

import (
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/utils"
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Env                    string `env:"APP_ENV" envDefault:"local"`
	ServerIP               string `env:"SERVER_IP" envDefault:"localhost"`
	ServerPort             string `env:"SERVER_PORT" envDefault:"8082"`
	AgentIP                string `env:"AGENT_IP" envDefault:"50051"`
	AgentPort              string `env:"AGENT_PORT" envDefault:"8080"`
	CollectMetricsInterval int    `env:"COLLECT_METRICS_INTERVAL" envDefault:"15"`
}

func NewConfigFromEnv(ctx context.Context) (*Config, error) {

	// read explicitly from .env file
	envPath := filepath.Join(filepath.Dir(os.Args[0]), ".env")
	if err := godotenv.Load(envPath); err != nil {
		logger.Warn(ctx, "Failed to load .env file, using OS env",
			slog.String("path", envPath))
	}

	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	if cfg.AgentIP == "" {
		cfg.AgentIP = utils.GetPublicIP(ctx)
		logger.Info(ctx, "AgentIP is not set, using public IP",
			slog.String("AgentIP", cfg.AgentIP))
	}

	if cfg.ServerIP == "" {
		logger.Warn(ctx, "ServerIP is not set, can be issues with the server_service connection")
	}

	logger.Info(ctx, "config loaded",
		slog.String("ServerIP", cfg.ServerIP),
		slog.String("ServerPort", cfg.ServerPort),
		slog.String("AgentIP", cfg.AgentIP),
		slog.String("AgentPort", cfg.AgentPort),
		slog.Int("CollectMetricsInterval", cfg.CollectMetricsInterval),
	)

	return cfg, nil
}
