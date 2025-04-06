package server_services

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/server_services/postgres_srvs"
	"Server-Monitoring-System/internal/server_services/redis_srvc"
)

type Services struct {
	Postgres *postgres_srvs.Service
	RedisS   *redis_srvc.Service
}

func NewServices(cfg *config.Config, clnts *clients.Clients) *Services {
	return &Services{
		Postgres: postgres_srvs.NewService(clnts),
		RedisS:   redis_srvc.NewService(cfg, clnts),
	}
}
