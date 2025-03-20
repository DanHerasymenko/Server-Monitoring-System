package service

import (
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"log/slog"
)

func (s *Service) SendMetrics(collectedMetrics *Metrics, client pb.MonitoringService_StreamMetricsClient) error {

	metricsReq := &pb.MetricsRequest{
		CpuUsage:  float32(collectedMetrics.CpuUsage),
		RamUsage:  float32(collectedMetrics.RamUsage),
		DiskUsage: float32(collectedMetrics.DiskUsage),
		Timestamp: collectedMetrics.Timestamp,
	}

	if err := client.Send(metricsReq); err != nil {
		logger.Error(s.context, fmt.Errorf("error sending metrics: %v", err))
		return err
	}

	logger.Info(s.context, "Metrics collected",
		slog.Any("CPU", metricsReq.CpuUsage),
		slog.Any("RAM", metricsReq.RamUsage),
		slog.Any("Disk", metricsReq.DiskUsage),
		slog.Any("Timestamp", metricsReq.Timestamp),
	)

	resp, err := client.Recv()
	if err != nil {
		logger.Error(s.context, fmt.Errorf("failed to receive response: %w", err))
		return err
	}

	logger.Info(s.context, "Response received", slog.Any("Response", resp))

	return nil
}
