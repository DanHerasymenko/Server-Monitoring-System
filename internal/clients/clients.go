package clients

import (
	"Server-Monitoring-System/internal/clients/postgres"
	"Server-Monitoring-System/internal/clients/redis"
	"Server-Monitoring-System/internal/config"
	"context"
)

type Clients struct {
	RedisClnt     *redis.Client
	PostgressClnt *postgres.Client
}

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {

	redisClient, err := redis.NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	postgresClient, err := postgres.NewPostgresClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	clients := &Clients{
		RedisClnt:     redisClient,
		PostgressClnt: postgresClient,
	}

	return clients, nil
}
