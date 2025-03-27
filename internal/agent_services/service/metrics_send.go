package service

import (
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"log/slog"
)

func (s *Service) SendMetrics(collectedMetrics *Metrics, client pb.MonitoringService_StreamMetricsClient) error {

	metricsReq := &pb.MetricsRequest{
		ServerIp:  s.cfg.AgentIP,
		CpuUsage:  collectedMetrics.CpuUsage,
		RamUsage:  collectedMetrics.RamUsage,
		DiskUsage: collectedMetrics.DiskUsage,
		Timestamp: collectedMetrics.Timestamp,
	}

	if err := client.Send(metricsReq); err != nil {
		return fmt.Errorf("error sending metrics: %w", err)
	}

	logger.Info(s.context, "Metrics sent",
		slog.Float64("CPU", metricsReq.CpuUsage),
		slog.Float64("RAM", metricsReq.RamUsage),
		slog.Float64("Disk", metricsReq.DiskUsage),
		slog.Int64("Timestamp", metricsReq.Timestamp),
	)

	resp, err := client.Recv()
	if err != nil {
		logger.Error(s.context, fmt.Errorf("failed to receive response: %w", err))
		return err
	}

	logger.Info(s.context, "Response received", slog.Any("Response", resp))

	return nil
}
