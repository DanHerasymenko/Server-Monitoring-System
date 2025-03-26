package clients

import (
	"Server-Monitoring-System/internal/clients/redis_clnt"
	"Server-Monitoring-System/internal/config"
	"context"
	"fmt"
)

type Clients struct {
	RedisClnt *redis_clnt.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {

	redisClient, err := redis_clnt.NewRedisClient("", "", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis_clnt client: %w", err)
	}

	clients := &Clients{
		RedisClnt: redisClient,
	}

	return clients, nil
}
