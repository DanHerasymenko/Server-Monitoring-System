package stream

import (
	"Server-Monitoring-System/internal/clients/redis_client"
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"io"
	"log/slog"
)

type Server struct {
	pb.UnimplementedMonitoringServiceServer
	Redis *redis_client.Client
}

func (s *Server) StreamMetrics(stream pb.MonitoringService_StreamMetricsServer) error {

	ctx := stream.Context()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			logger.Info(ctx, "Client closed the stream")
			return nil
		}
		if err != nil {
			logger.Error(ctx, fmt.Errorf("failed to receive metrics: %w", err))
			return err
		}

		logger.Info(ctx, "Server received metrics",
			slog.String("server_ip", req.ServerIp),
			slog.Any("cpu", req.CpuUsage),
			slog.Any("ram", req.RamUsage),
			slog.Any("disk", req.DiskUsage),
			slog.Int64("timestamp", req.Timestamp),
		)

		// Send response back to client
		resp := &pb.MetricsResponse{Status: "Metrics received"}
		if err := stream.Send(resp); err != nil {
			logger.Error(ctx, fmt.Errorf("failed to send response: %w", err))
			return err
		}
	}
}
