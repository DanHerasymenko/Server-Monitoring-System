package utils

import (
	"Server-Monitoring-System/internal/logger"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetPublicIP(ctx context.Context) string {
	resp, err := http.Get("https://api64.ipify.org")
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("failed to get public IP: %v", err))
		return "unknown"
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("failed to read public IP: %v", err))
		return "unknown"
	}

	return strings.TrimSpace(string(ip))
}
