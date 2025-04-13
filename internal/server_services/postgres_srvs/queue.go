package postgres_srvs

import (
	"Server-Monitoring-System/internal/constants"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"time"
)

type MetricsItem struct {
	Metric     *pb.MetricsRequest
	EnqueuedAt time.Time
}

type ServerWorker struct {
	MetricQueue chan MetricsItem
}

func (srvs *Service) NewServerWorker() *ServerWorker {
	return &ServerWorker{
		MetricQueue: make(chan MetricsItem, constants.MetricQueueSize),
	}
}

func (srvs *Service) AddMetricsToQueueWithTimeout(metricQueue chan MetricsItem, metricsItem MetricsItem) error {
	select {
	case metricQueue <- metricsItem:
		return nil
	case <-time.After(constants.MetricQueueTimeout):
		return fmt.Errorf("timeout after %v", constants.MetricQueueTimeout)
	}
}
