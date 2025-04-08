package stream

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/server_services"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"io"
	"log/slog"
)

type Server struct {
	pb.UnimplementedMonitoringServiceServer
	Clients     *clients.Clients
	Services    *server_services.Services
	MetricQueue chan *pb.MetricsRequest
	Ctx         context.Context
}

func (s *Server) StreamMetrics(stream pb.MonitoringService_StreamMetricsServer) error {

	for {
		// Receive metrics from the client
		req, err := stream.Recv()
		if err == io.EOF {
			logger.Info(s.Ctx, "Client closed the stream")
			return nil
		}
		if err != nil {
			logger.Error(s.Ctx, fmt.Errorf("failed to receive metrics: %w", err))
			return err
		}

		logger.Info(s.Ctx, "Server received metrics",
			slog.String("server_ip", req.ServerIp),
			slog.Float64("cpu", req.CpuUsage),
			slog.Float64("ram", req.RamUsage),
			slog.Float64("disk", req.DiskUsage),
			slog.Int64("timestamp", req.Timestamp),
		)

		// Save metrics to Postgres
		err = s.Services.RedisS.SaveMetrics(s.Ctx, req)
		if err != nil {
			logger.Error(s.Ctx, fmt.Errorf("failed to save metrics after receiving: %w", err))
			return err
		}
		logger.Info(s.Ctx, "Metrics saved to Redis successfully")

		// Add metrics to the queue with timeout
		err = s.Services.Postgres.AddMetricsToQueueWithTimeout(s.MetricQueue, req)
		if err != nil {
			logger.Error(s.Ctx, fmt.Errorf("failed to add metrics to queue: %w", err))
			return err
		}

		// Send response back to client
		resp := &pb.MetricsResponse{Status: "Metrics received on the server"}
		if err := stream.Send(resp); err != nil {
			logger.Error(s.Ctx, fmt.Errorf("failed to send response: %w", err))
			return err
		}
	}
}
