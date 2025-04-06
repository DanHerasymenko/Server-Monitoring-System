package postgres_srvs

import (
	"Server-Monitoring-System/internal/clients"
)

type Service struct {
	clnts *clients.Clients
}

func NewService(clnts *clients.Clients) *Service {
	return &Service{
		clnts: clnts,
	}
}
