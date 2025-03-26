package redis_srvc

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	"context"
	"fmt"
)

type Service struct {
	cfg   *config.Config
	clnts *clients.Clients
}

func NewService(cfg *config.Config, clnts *clients.Clients) *Service {
	return &Service{
		cfg:   cfg,
		clnts: clnts,
	}
}

func (srvs *Service) Ping(ctx context.Context) error {
	if err := srvs.clnts.RedisClnt.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}
