package agent_services

import (
	"Server-Monitoring-System/internal/agent_services/service"
	"Server-Monitoring-System/internal/config"
	pb "Server-Monitoring-System/proto"
	"context"
)

type Services struct {
	Agent *service.Service
}

func NewServices(cfg *config.Config, ctx context.Context, cancel context.CancelFunc, client pb.MonitoringService_StreamMetricsClient) *Services {
	return &Services{
		Agent: service.NewService(cfg, ctx, cancel, client),
	}
}
