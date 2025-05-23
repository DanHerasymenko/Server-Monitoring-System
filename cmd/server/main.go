package main

import (
	"Server-Monitoring-System/cmd/server/metrics"
	"Server-Monitoring-System/cmd/server/stream"
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/server_services"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// simlulate agent with `grpcui.exe -plaintext localhost:50051`
// configure workers settings in constants

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer logger.Close()

	// load config
	cfg, err := config.NewConfigFromEnv(ctx)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err))
	}

	// initialize clients
	clnts, err := clients.NewClients(ctx, cfg)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to initialize clients: %w", err))
	}

	// initialize services
	srvc := server_services.NewServices(cfg, clnts)

	// start workers
	queue := srvc.Postgres.NewServerWorker()
	srvc.Postgres.StartWorkerPool(ctx, queue)

	// listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to listen: %w", err))
	}

	// create a gRPC server_services object
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// register the server_services with the gRPC server_services
	pb.RegisterMonitoringServiceServer(grpcServer, &stream.Server{
		Clients:     clnts,
		Services:    srvc,
		MetricQueue: queue.MetricQueue,
		Ctx:         ctx,
	})

	lisAddrStr := lis.Addr().String()

	// start metrics server
	metrics.StartMetricsServer(ctx)
	logger.Info(ctx, "Metrics server started")

	// start the server_services
	logger.Info(ctx, "Starting server_services...", slog.String("address", lisAddrStr))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Panic(ctx, fmt.Errorf("failed to start server_services: %w", err))
		}
	}()
	logger.Info(ctx, "Server started")

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info(ctx, "Shutting down server_services...")
	grpcServer.GracefulStop()
	lis.Close()
	logger.Info(ctx, "Server stopped")

}
