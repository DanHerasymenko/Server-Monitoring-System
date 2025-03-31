package redis_srvc

import (
	"Server-Monitoring-System/internal/logger"
	"context"
)

func (s *Service) PingCheck(ctx context.Context) error {

	if err := s.clnts.RedisClnt.Redis.Ping(ctx).Err(); err != nil {
		return err
	} else if err == nil {
		logger.Info(ctx, "Redis ping successful")
		return nil
	}
	return nil
}
