package agent_service

import (
	"Server-Monitoring-System/internal/agent_service/service"
	"Server-Monitoring-System/internal/config"
	"context"
)

type Services struct {
	Agent *service.Service
}

func NewServices(cfg *config.Config, ctx context.Context, cancel context.CancelFunc) *Services {
	return &Services{
		Agent: service.NewService(cfg, ctx, cancel),
	}
}
