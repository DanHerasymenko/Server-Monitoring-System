package postgres_srvs

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
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
