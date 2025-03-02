package main

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/services"
	"context"
	"fmt"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfigFromEnv(ctx)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err))
	}

	svs := services.NewServices(cfg, ctx)
	svs.Agent.RunAgentService()
}
