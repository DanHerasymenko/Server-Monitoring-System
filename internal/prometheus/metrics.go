package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "sms_active_connections",
	Help: "Current number of active connections to the server",
})

var WorkerPoolSize = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "sms_worker_pool_size",
	Help: "Current size of the worker pool",
})

var QueueLength = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "sms_queue_length",
	Help: "Current length of the queue",
})

var MetricsReceivedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "sms_metrics_received_total",
	Help: "Total number of metrics received from agents",
})

var DBWriteDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "sms_db_write_duration_seconds",
	Help:    "Duration of PostgreSQL write operations in seconds",
	Buckets: prometheus.DefBuckets,
})
