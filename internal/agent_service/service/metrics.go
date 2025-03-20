package service

import (
	"Server-Monitoring-System/internal/logger"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"log/slog"
	"time"
)

type Metrics struct {
	CpuUsage  float64
	RamUsage  float64
	DiskUsage float64
	Timestamp int64
}

func (s *Service) CollectMetrics() {

	m, err := getMetrics()
	if err != nil {
		logger.Error(s.context, fmt.Errorf("error getting metrics %v", err))
		return
	}

	logger.Info(s.context, "Metrics collected",
		slog.Float64("CPU", m.CpuUsage),
		slog.Float64("RAM", m.RamUsage),
		slog.Float64("Disk", m.DiskUsage),
	)

	time.Sleep(time.Duration(s.cfg.CollectMetricsInterval) * time.Second)
}

func getMetrics() (*Metrics, error) {

	now := time.Now().Unix()

	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	memV, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	diskU, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &Metrics{
		CpuUsage:  cpuUsage[0],
		RamUsage:  memV.UsedPercent,
		DiskUsage: diskU.UsedPercent,
		Timestamp: now,
	}, nil

}
