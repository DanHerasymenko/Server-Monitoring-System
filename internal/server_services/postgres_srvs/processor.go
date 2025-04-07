package postgres_srvs

import (
	"Server-Monitoring-System/internal/constants"
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"context"
	"fmt"
	"time"
)

func (srvs *Service) StartWorkerPool(ctx context.Context, queue *ServerWorker) {
	for i := 0; i < constants.WorkerCount; i++ {
		go worker(ctx, srvs, queue)
	}
}

func worker(ctx context.Context, srvs *Service, queue *ServerWorker) {

	batch := make([]*pb.MetricsRequest, 0, constants.WorkerBatchSize)
	// ticker to flush the queue if batch size is not reached (to not freeze and not to lose data)
	ticker := time.NewTicker(constants.WorkerFlushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "Worker stopped(ctx.Done())")
			return
		case metric := <-queue.MetricQueue:
			batch = append(batch, metric)
			if len(batch) >= constants.WorkerBatchSize {
				logger.Info(ctx, "Batch is full: flushing batch to DB")
				err := srvs.SaveBatchMetricsToPostgres(ctx, batch)
				if err != nil {
					logger.Error(ctx, fmt.Errorf("failed to save batch to DB: %w", err))
				} else {
					logger.Info(ctx, "Batch saved to DB successfully")
				}
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				logger.Info(ctx, "Ticker: flushing batch to DB")
				err := srvs.SaveBatchMetricsToPostgres(ctx, batch)
				if err != nil {
					logger.Error(ctx, fmt.Errorf("failed to save batch to DB: %w", err))
				} else {
					logger.Info(ctx, "Batch saved to DB successfully")
				}
				batch = batch[:0]
			}
		}
	}

}
