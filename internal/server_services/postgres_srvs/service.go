package postgres_srvs

import (
	"Server-Monitoring-System/internal/clients"
	"Server-Monitoring-System/internal/logger"
	pb "Server-Monitoring-System/proto"
	"context"
)

type Service struct {
	clnts *clients.Clients
}

func NewService(clnts *clients.Clients) *Service {
	return &Service{
		clnts: clnts,
	}
}

func (srvs *Service) SaveBatchMetricsToPostgres(ctx context.Context, batch []*pb.MetricsRequest) error {

	tx, err := srvs.clnts.PostgressClnt.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range batch {
		var serverId int

		// find server
		err := tx.QueryRow(ctx, `SELECT id FROM servers WHERE ip_address = $1`, m.ServerIp).Scan(&serverId)

		// if there is no row, it will return pgx.ErrNoRows- insert new ip
		if err != nil {
			logger.Info(ctx, "Server not found, inserting new IP")
			err = tx.QueryRow(ctx, `INSERT INTO servers (ip_address, last_active) VALUES ($1, $2)RETURNING id`, m.ServerIp, m.Timestamp).Scan(&serverId)
			if err != nil {
				return err
			}
		} else {
			// update last active
			logger.Info(ctx, "Server found, updating last active")
			_, err := tx.Exec(ctx, `UPDATE servers SET last_active = $1 WHERE id = $2`, m.Timestamp, serverId)
			if err != nil {
				return err
			}
		}

		// insert metrics
		_, err = tx.Exec(ctx, `INSERT INTO metrics (server_id, cpu_usage, ram_usage, disk_usage, timestamp) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
			serverId, m.CpuUsage, m.RamUsage, m.DiskUsage, m.Timestamp)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
