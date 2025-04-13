package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "sms_active_connections",
	Help: "Current number of active connections to the server",
})

var MetricsReceivedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "sms_metrics_received_total",
	Help: "Total number of metrics received from agents",
})

var DBWriteDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "sms_db_write_duration_seconds",
	Help:    "Duration of PostgreSQL write operations in seconds: receive metric from agent -> saving to DB",
	Buckets: prometheus.DefBuckets,
})

var QueueDelaySeconds = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "sms_queue_delay_seconds",
	Help:    "Delay in seconds until metrics are handled from the queue",
	Buckets: prometheus.DefBuckets,
})

var GRPCRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "sms_grpc_requests_total",
	Help: "Total number of gRPC requests received",
})
