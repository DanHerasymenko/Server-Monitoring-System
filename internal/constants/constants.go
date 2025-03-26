package constants

const (
	ServiceName        = "MonitoringAgent"
	ServiceDisplayName = "Monitoring Agent"
	ServiceDescription = "Collects metrics and sends them to the server_services via gRPC streaming"
	DependencieNetwork = "Requires=network.target"
	DependencieAfter   = "After=network-online.target syslog.target"
)
