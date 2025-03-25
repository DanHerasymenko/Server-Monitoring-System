package main

import (
	"Server-Monitoring-System/internal/agent_service"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfigFromEnv(ctx)

	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err))
	}

	// new gRPC connection to server without TLS
	serverAddr := cfg.ServerIP + ":" + cfg.ServerPort
	fmt.Println(serverAddr)
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to connect to server: %w", err))
	}
	defer conn.Close()

	//create a gRPC client
	client := pb.NewMonitoringServiceClient(conn)

	stream, err := client.StreamMetrics(ctx)
	if err != nil {
		log.Fatalf("Error opening grpc stream: %v", err)
	}

	svs := agent_service.NewServices(cfg, ctx, cancel, stream)
	defer logger.Close()
	svs.Agent.RunAgentService()
}
