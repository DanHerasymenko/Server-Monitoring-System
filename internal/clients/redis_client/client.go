package redis_client

import "github.com/redis/go-redis/v9"

type Client struct {
	Redis *redis.Client
}

func NewRedisClient(addr, password string, db int) (*Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       db,
	})

	return &Client{Redis: client}, nil

}
