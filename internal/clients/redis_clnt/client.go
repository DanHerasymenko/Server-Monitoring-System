package redis_clnt

import (
	"Server-Monitoring-System/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis *redis.Client
}

func NewRedisClient(cfg *config.Config) *Client {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &Client{Redis: client}

}
