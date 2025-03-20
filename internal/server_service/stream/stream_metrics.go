package stream

import (
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"fmt"
	"io"
	"log/slog"
)

type Server struct {
	pb.UnimplementedMonitoringServiceServer
}

func (s *Server) StreamMetrics(stream pb.MonitoringService_StreamMetricsServer) error {

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "Client disconnected or timeout occurred")
			return ctx.Err()
		default:
			req, err := stream.Recv() // Завершення стріму
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err // Помилка під час отримання даних
			}

			// Логуємо отримані метрики
			logger.Info(ctx, "Received metrics",
				slog.String("server_ip", req.ServerIp),
				slog.Any("cpu", req.CpuUsage), // Any підтримує float64
				slog.Any("ram", req.RamUsage),
				slog.Any("disk", req.DiskUsage),
				slog.Int64("timestamp", req.Timestamp),
			)

			// Відповідь назад клієнту
			resp := &pb.MetricsResponse{Status: "Metrics received"}
			if err := stream.Send(resp); err != nil {
				logger.Error(ctx, fmt.Errorf("failed to send response: %w", err))
				return err
			}
		}

	}
}
