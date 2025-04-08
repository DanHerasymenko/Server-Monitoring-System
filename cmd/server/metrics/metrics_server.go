package metrics

import (
	"Server-Monitoring-System/internal/logger"
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func StartMetricsServer(ctx context.Context) {

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info(ctx, "Starting metrics server on :2112")
		if err := http.ListenAndServe(":2112", nil); err != nil {
			logger.Fatal(ctx, err)
		}
	}()
}
