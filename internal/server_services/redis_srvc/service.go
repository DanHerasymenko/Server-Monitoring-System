package redis_srvc

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/config"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"strings"
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

func (srvs *Service) SaveMetrics(ctx context.Context, agentIP string, metrics *pb.MetricsRequest) error {

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

	metrics, err := ConvertRedisDataToMetrics(agentIP, values)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics by IP: %w", err)
	}

	return metrics, nil
}

func (srvs *Service) GetAllMetrics(ctx context.Context) (map[string]*pb.MetricsRequest, error) {

	var (
		cursor        uint64
		allKeys       []string
		allMetricsMap = make(map[string]*pb.MetricsRequest)
	)

	// scan all keys with pattern "metrics:*" and 100 keys for 1 iteration - use scan to not block
	for {

		keys, cursor, err := srvs.clnts.RedisClnt.Redis.Scan(ctx, cursor, "metrics:*", 100).Result()

		if err != nil {
			fmt.Errorf("failed to scan keys: %w", err)
		}

		if cursor == 0 {
			break
		}

		allKeys = append(allKeys, keys...)
	}

	// get all metrics for each key
	for _, key := range allKeys {

		ip := strings.TrimPrefix(key, "metrics:")
		data, err := srvs.clnts.RedisClnt.Redis.HGetAll(ctx, ip).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get all metrics: %w", err)
		} else if len(data) == 0 {
			return nil, fmt.Errorf("no metrics found: %s", ip)
		}

		allMetricsMap[ip], err = ConvertRedisDataToMetrics(ip, data)
	}

	return allMetricsMap, nil

}
