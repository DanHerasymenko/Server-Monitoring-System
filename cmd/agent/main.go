package main

import (
	"Server-Monitoring-System/internal/agent_service"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/logger"
	"context"
	"fmt"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfigFromEnv(ctx)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err))
	}

	svs := agent_service.NewServices(cfg, ctx, cancel)
	defer logger.Close()
	svs.Agent.RunAgentService()
}
