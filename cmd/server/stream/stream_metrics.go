package stream

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/prometheus"
	"Server-Monitoring-System/internal/server_services"
	"Server-Monitoring-System/internal/server_services/postgres_srvs"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

type Server struct {
	pb.UnimplementedMonitoringServiceServer
	Clients     *clients.Clients
	Services    *server_services.Services
	MetricQueue chan postgres_srvs.MetricsItem
	Ctx         context.Context
}

func (s *Server) StreamMetrics(stream pb.MonitoringService_StreamMetricsServer) error {

	// Increment Prometheus metric: GRPCRequestsTotal
	prometheus.GRPCRequestsTotal.Inc()

	// Increment Prometheus metric: ActiveConnections
	prometheus.ActiveConnections.Inc()
	defer prometheus.ActiveConnections.Dec()

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

		// Increment Prometheus metric: MetricsReceivedTotal
		prometheus.MetricsReceivedTotal.Inc()

		logger.Info(s.Ctx, "Server received metrics",
			slog.String("server_ip", req.ServerIp),
			slog.Float64("cpu", req.CpuUsage),
			slog.Float64("ram", req.RamUsage),
			slog.Float64("disk", req.DiskUsage),
			slog.Int64("timestamp", req.Timestamp),
		)

		// Save metrics to Redis
		err = s.Services.RedisS.SaveMetrics(s.Ctx, req)
		if err != nil {
			logger.Error(s.Ctx, fmt.Errorf("failed to save metrics after receiving: %w", err))
			return err
		}
		logger.Info(s.Ctx, "Metrics saved to Redis successfully")

		// Add metrics to the queue with timeout
		// Save to PostgreSQL

		item := postgres_srvs.MetricsItem{
			Metric:     req,
			EnqueuedAt: time.Now(),
		}

		err = s.Services.Postgres.AddMetricsToQueueWithTimeout(s.MetricQueue, item)
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
