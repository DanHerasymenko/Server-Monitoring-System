package main

import (
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type server struct {
	pb.MonitoringServiceServer
}

//func (s *server) SayHello(ctx context.Context, in *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
//	return &pb.HelloWorldResponse{Message: "Hello, World! "}, nil
//}

func main() {

	ctx, _ := context.WithCancel(context.Background())

	defer logger.Close()

	// listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to listen: %w", err))
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// register the server with the gRPC server
	pb.RegisterMonitoringServiceServer(grpcServer, &server{})

	// start the server
	lisAddrStr := lis.Addr().String()
	logger.Info(ctx, "Starting server...", slog.String("address", lisAddrStr))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to serve: %w", err))
	}
}
