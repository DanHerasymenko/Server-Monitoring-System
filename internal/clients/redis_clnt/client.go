package redis_clnt

import (
	"Server-Monitoring-System/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	Redis *redis.Client
}

func NewRedisClient(cfg *config.Config) *Client {

	//client := redis.NewClient(&redis.Options{
	//	Addr:     cfg.RedisAddr,
	//	Username: cfg.RedisUser,
	//	Password: cfg.RedisUserPassword,
	//	DB:       cfg.RedisDB,
	//})

	redisURL := fmt.Sprintf("redis://%s:%s@%s/%d",
		cfg.RedisUser, cfg.RedisUserPassword, cfg.RedisAddr, cfg.RedisDB)

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	return &Client{Redis: client}

}
