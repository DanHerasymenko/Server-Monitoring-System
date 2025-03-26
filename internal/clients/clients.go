package clients

import (
	"Server-Monitoring-System/internal/clients/redis_client"
	"Server-Monitoring-System/internal/config"
	"context"
	"fmt"
)

type Clients struct {
	Redis *redis_client.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {

	redisClient, err := redis_client.NewRedisClient("", "", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create redis_client client: %w", err)
	}

	clients := &Clients{
		Redis: redisClient,
	}

	return clients, nil
}
