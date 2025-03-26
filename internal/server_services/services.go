package server_services

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/server_services/redis_srvc"
)

type Services struct {
	//Metrics *metrics.Service
	RedisS *redis_srvc.Service
}

func NewServices(cfg *config.Config, clnts *clients.Clients) *Services {
	return &Services{
		//Metrics: metrics.NewService(),
		RedisS: redis_srvc.NewService(cfg, clnts),
	}
}
