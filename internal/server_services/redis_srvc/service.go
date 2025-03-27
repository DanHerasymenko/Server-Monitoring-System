package redis_srvc

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	"Server-Monitoring-System/internal/utils"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
)

type Service struct {
	cfg   *config.Config
	clnts *clients.Clients
}

func NewService(cfg *config.Config, clnts *clients.Clients) *Service {
	return &Service{
		cfg:   cfg,
		clnts: clnts,
	}
}

func (srvs *Service) Ping(ctx context.Context) error {
	if err := srvs.clnts.RedisClnt.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

func (srvs Service) SaveMetrics(ctx context.Context, agentIP string, metrics *pb.MetricsRequest) error {

	key := fmt.Sprintf("metrics:%s", agentIP)
	value := map[string]interface{}{
		"cpu":       metrics.CpuUsage,
		"ram":       metrics.RamUsage,
		"disk":      metrics.DiskUsage,
		"timestamp": metrics.Timestamp,
	}

	return srvs.clnts.RedisClnt.Redis.HSet(ctx, key, value).Err()
}

func (srvs *Service) GetMetricsByIp(ctx context.Context, agentIP string) (*pb.MetricsRequest, error) {

	key := fmt.Sprintf("metrics:%s", agentIP)
	values, err := srvs.clnts.RedisClnt.Redis.HGetAll(ctx, key).Result() // returns string: map[string]string

	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no metrics found for agent: %s", agentIP)
	}

	return &pb.MetricsRequest{
		ServerIp:  agentIP,
		CpuUsage:  utils.ParseFloat(values["cpu"]),
		RamUsage:  utils.ParseFloat(values["ram"]),
		DiskUsage: utils.ParseFloat(values["disk"]),
		Timestamp: utils.ParseInt(values["timestamp"]),
	}, nil
}
