package redis_srvc

import (
	"Server-Monitoring-System/internal/utils"
	pb "Server-Monitoring-System/proto"
	"fmt"
)

// ConvertRedisDataToMetrics Redis hash -> pb.MetricsRequest
func ConvertRedisDataToMetrics(ip string, data map[string]string) (*pb.MetricsRequest, error) {
	cpu, err := utils.ParseFloat(data["cpu"])
	if err != nil {
		return nil, fmt.Errorf("invalid CPU value: %w", err)
	}

	ram, err := utils.ParseFloat(data["ram"])
	if err != nil {
		return nil, fmt.Errorf("invalid RAM value: %w", err)
	}

	disk, err := utils.ParseFloat(data["disk"])
	if err != nil {
		return nil, fmt.Errorf("invalid Disk value: %w", err)
	}

	timestamp, err := utils.ParseInt64(data["timestamp"])
	if err != nil {
		return nil, fmt.Errorf("invalid Timestamp value: %w", err)
	}

	return &pb.MetricsRequest{
		ServerIp:  ip,
		CpuUsage:  cpu,
		RamUsage:  ram,
		DiskUsage: disk,
		Timestamp: timestamp,
	}, nil
}
