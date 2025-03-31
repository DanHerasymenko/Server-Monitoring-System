package clients

import (
	"Server-Monitoring-System/internal/clients/redis_clnt"
	"Server-Monitoring-System/internal/config"
)

type Clients struct {
	RedisClnt *redis_clnt.Client
}

func NewClients(cfg *config.Config) (*Clients, error) {

	redisClient := redis_clnt.NewRedisClient(cfg)
	//postgressClient := postgress_clnt.NewPostgressClient(cfg)

	clients := &Clients{
		RedisClnt: redisClient,
		//PostgressClnt: postgressClient,
	}

	return clients, nil
}
