package stats

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type SystemStats struct {
	CPU     CPUStats     `json:"cpu"`
	Memory  MemoryStats  `json:"memory"`
	Disk    DiskStats    `json:"disk"`
	Network NetworkStats `json:"network"`
	Uptime  int64        `json:"uptime"`
}

type CPUStats struct {
	Usage   float64   `json:"usage"`
	Cores   int       `json:"cores"`
	Model   string    `json:"model"`
	LoadAvg []float64 `json:"loadAvg"`
}

type MemoryStats struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

type DiskStats struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type NetworkStats struct {
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
}

var startTime = time.Now()

func GetSystemStats() (*SystemStats, error) {
	stats := &SystemStats{}

	// CPU Stats
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err == nil && len(cpuPercent) > 0 {
		stats.CPU.Usage = cpuPercent[0]
	}
	stats.CPU.Cores = runtime.NumCPU()

	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		stats.CPU.Model = cpuInfo[0].ModelName
	}

	// Memory Stats
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		stats.Memory.Total = memInfo.Total
		stats.Memory.Available = memInfo.Available
		stats.Memory.Used = memInfo.Used
		stats.Memory.UsedPercent = memInfo.UsedPercent
		stats.Memory.Free = memInfo.Free
	}

	// Disk Stats (root partition)
	diskInfo, err := disk.Usage("/")
	if err == nil {
		stats.Disk.Total = diskInfo.Total
		stats.Disk.Free = diskInfo.Free
		stats.Disk.Used = diskInfo.Used
		stats.Disk.UsedPercent = diskInfo.UsedPercent
	}

	// Network Stats
	netStats, err := net.IOCounters(false)
	if err == nil && len(netStats) > 0 {
		stats.Network.BytesSent = netStats[0].BytesSent
		stats.Network.BytesRecv = netStats[0].BytesRecv
		stats.Network.PacketsSent = netStats[0].PacketsSent
		stats.Network.PacketsRecv = netStats[0].PacketsRecv
	}

	// Uptime
	stats.Uptime = int64(time.Since(startTime).Seconds())

	return stats, nil
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatUptime(seconds int64) string {
	duration := time.Duration(seconds) * time.Second
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

