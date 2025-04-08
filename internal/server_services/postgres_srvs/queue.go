package postgres_srvs

import (
	"Server-Monitoring-System/internal/constants"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"time"
)

type ServerWorker struct {
	MetricQueue chan *pb.MetricsRequest
}

func (srvs *Service) NewServerMetricsQueue() *ServerWorker {
	return &ServerWorker{
		MetricQueue: make(chan *pb.MetricsRequest, constants.MetricQueueSize),
	}
}

func (srvs *Service) AddMetricsToQueueWithTimeout(metricQueue chan *pb.MetricsRequest, req *pb.MetricsRequest) error {
	select {
	case metricQueue <- req:
		return nil
	case <-time.After(constants.MetricQueueTimeout):
		return fmt.Errorf("timeout after %v", constants.MetricQueueTimeout)
	}
}
