package config

import (
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/utils"
	"context"
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

type Config struct {
	Env                    string `env:"APP_ENV" envDefault:"local"`
	ServerIP               string `env:"SERVER_IP" envDefault:"localhost"`
	ServerPort             string `env:"SERVER_PORT" envDefault:"50051"`
	AgentIP                string `env:"AGENT_IP" envDefault:"localhost"`
	AgentPort              string `env:"AGENT_PORT" envDefault:"50052"`
	CollectMetricsInterval int    `env:"COLLECT_METRICS_INTERVAL" envDefault:"30"`
	RedisAddr              string `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	RedisPassword          string `env:"REDIS_PASSWORD"`
	RedisUser              string `env:"REDIS_USER"`
	RedisUserPassword      string `env:"REDIS_USER_PASSWORD"`
	RedisDB                int    `env:"REDIS_DB"`
	PostgresHost           string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort           int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser           string `env:"POSTGRES_USER"`
	PostgresPassword       string `env:"POSTGRES_PASSWORD"`
	PostgresDB             string `env:"POSTGRES_DB"`
}

func NewConfigFromEnv(ctx context.Context) (*Config, error) {

	//// read explicitly from .env file
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envPath := filepath.Join(dir, ".env")

	if err := godotenv.Load(envPath); err != nil {
		logger.Warn(ctx, "Failed to load .env file, using OS env",
			slog.String("path", envPath),
			slog.Any("err", err))
	}

	cfg := &Config{}

	err = env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	if cfg.AgentIP == "" {
		cfg.AgentIP = utils.GetPublicIP(ctx)
		logger.Info(ctx, "AgentIP is not set, using public IP",
			slog.String("AgentIP", cfg.AgentIP))
	}

	if cfg.ServerIP == "" {
		logger.Warn(ctx, "ServerIP is not set, can be issues with the server_services connection")
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
