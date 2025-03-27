package stream

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/server_services"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"io"
	"log/slog"
)

type Server struct {
	pb.UnimplementedMonitoringServiceServer
	Clients  *clients.Clients
	Services *server_services.Services
}

func (s *Server) StreamMetrics(stream pb.MonitoringService_StreamMetricsServer) error {

	ctx := stream.Context()

	if err := s.Services.RedisS.Ping(ctx); err != nil {
		logger.Error(ctx, err)
	}
	logger.Info(ctx, "Redis ping successful")

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
			slog.Float64("cpu", req.CpuUsage),
			slog.Float64("ram", req.RamUsage),
			slog.Float64("disk", req.DiskUsage),
			slog.Int64("timestamp", req.Timestamp),
		)

		err = s.Services.RedisS.SaveMetrics(ctx, req.ServerIp, req)
		if err != nil {
			logger.Error(ctx, fmt.Errorf("failed to save metrics after receiving: %w", err))
		}
		logger.Info(ctx, "Metrics saved successfully")

		metricsbyIP, err := s.Services.RedisS.GetMetricsByIp(ctx, req.ServerIp)
		fmt.Println(metricsbyIP)

		allMetrics, err := s.Services.RedisS.GetAllMetrics(ctx)
		fmt.Println(allMetrics)

		// Send response back to client
		resp := &pb.MetricsResponse{Status: "Metrics received"}
		if err := stream.Send(resp); err != nil {
			logger.Error(ctx, fmt.Errorf("failed to send response: %w", err))
			return err
		}
	}
}
