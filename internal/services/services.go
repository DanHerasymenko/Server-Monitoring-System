package services

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/services/agent"
	"context"
)

type Services struct {
	Agent *agent.Service
}

func NewServices(cfg *config.Config, ctx context.Context, cancel context.CancelFunc) *Services {
	return &Services{
		Agent: agent.NewService(cfg, ctx, cancel),
	}
}
