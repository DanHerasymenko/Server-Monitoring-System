package main

import (
	"Server-Monitoring-System/cmd/server/stream"
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/logger"
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
	_, err = clients.NewClients(ctx, cfg)

	// listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to listen: %w", err))
	}

	// create a gRPC server_service object
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	// register the server_service with the gRPC server_service
	pb.RegisterMonitoringServiceServer(grpcServer, &stream.Server{})

	lisAddrStr := lis.Addr().String()

	// start the server_service
	logger.Info(ctx, "Starting server_service...", slog.String("address", lisAddrStr))
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Panic(ctx, fmt.Errorf("failed to start server_service: %w", err))
		}
	}()
	logger.Info(ctx, "Server started")

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info(ctx, "Shutting down server_service...")
	grpcServer.GracefulStop()
	lis.Close()
	logger.Info(ctx, "Server stopped")

}
