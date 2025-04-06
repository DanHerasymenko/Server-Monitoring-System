package postgres_srvs

import (
	"Server-Monitoring-System/internal/constants"
	pb "Server-Monitoring-System/proto"
)

type ServerWorker struct {
	MetricQueue chan *pb.MetricsRequest
}

func (srvs *Service) NewServerMetricsQueue() *ServerWorker {
	return &ServerWorker{
		MetricQueue: make(chan *pb.MetricsRequest, constants.WorkerQueueSize),
	}
}
