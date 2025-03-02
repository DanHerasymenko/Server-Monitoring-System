package main

import (
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/logger"
	"Server-Monitoring-System/internal/metrics"
	"context"
	"fmt"
	"log/slog"
)

func main() {

	ctx := context.Background()

	cfg, err := config.NewConfigFromEnv(ctx)
	if err != nil {
		logger.Fatal(ctx, fmt.Errorf("failed to load config: %w", err))
	}

	// add server and agent IP to context
	ctx = logger.SetServerIP(ctx, cfg.ServerIp)
	ctx = logger.SetAgentIP(ctx, cfg.AgentIP)

	// get metrics from the agent system
	m, err := metrics.GetMetrics()
	if err != nil {
		logger.Error(ctx, fmt.Errorf("agent error in getting metrics: %v", err))
		return
	}

	// print metrics
	logger.Info(ctx, "Agent metrics",
		slog.Float64("CPU", m.CpuUsage),
		slog.Float64("RAM", m.RamUsage),
		slog.Float64("Disk", m.DiskUsage),
	)

}
