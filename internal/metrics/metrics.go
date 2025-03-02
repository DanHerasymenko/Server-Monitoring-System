package metrics

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

type Metrics struct {
	CpuUsage  float64
	RamUsage  float64
	DiskUsage float64
	Timestamp int64
}

func GetMetrics() (*Metrics, error) {

	now := time.Now().Unix()

	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	mem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	disk, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &Metrics{
		CpuUsage:  cpuUsage[0],
		RamUsage:  mem.UsedPercent,
		DiskUsage: disk.UsedPercent,
		Timestamp: now,
	}, nil
}
