package constants

import "time"

const (
	ServiceName        = "MonitoringAgent"
	ServiceDisplayName = "Monitoring Agent"
	ServiceDescription = "Collects metrics and sends them to the server_services via gRPC streaming"
	DependencyNetwork  = "Requires=network.target"
	DependencyAfter    = "After=network-online.target syslog.target"
	WorkerCount        = 5
	WorkerBatchSize    = 100
	WorkerFlushTimeout = 5 * time.Second
	MetricQueueSize    = 10000
	MetricQueueTimeout = 10 * time.Second
)
