package stats

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/net"
)

type NetworkInterface struct {
	Name        string `json:"name"`
	BytesSent   uint64 `json:"bytesSent"`
	BytesRecv   uint64 `json:"bytesRecv"`
	PacketsSent uint64 `json:"packetsSent"`
	PacketsRecv uint64 `json:"packetsRecv"`
	Errors      uint64 `json:"errors"`
	Drops       uint64 `json:"drops"`
	Speed       uint64 `json:"speed"`
	IsUp        bool   `json:"isUp"`
}

type NetworkBandwidth struct {
	Interface    string  `json:"interface"`
	UploadMbps   float64 `json:"uploadMbps"`
	DownloadMbps float64 `json:"downloadMbps"`
}

// GetNetworkInterfaces returns information about all network interfaces
func GetNetworkInterfaces() ([]NetworkInterface, error) {
	interfaces := []NetworkInterface{}

	// Get network interface statistics
	netStats, err := net.IOCounters(true)
	if err != nil {
		return interfaces, err
	}

	for _, stat := range netStats {
		iface := NetworkInterface{
			Name:        stat.Name,
			BytesSent:   stat.BytesSent,
			BytesRecv:   stat.BytesRecv,
			PacketsSent: stat.PacketsSent,
			PacketsRecv: stat.PacketsRecv,
			Errors:      stat.Errin + stat.Errout,
			Drops:       stat.Dropin + stat.Dropout,
			IsUp:        true, // Assume up if stats are available
		}
		interfaces = append(interfaces, iface)
	}

	return interfaces, nil
}

// GetNetworkBandwidth calculates network bandwidth usage over a time interval
func GetNetworkBandwidth(interval time.Duration) ([]NetworkBandwidth, error) {
	// Get initial stats
	initialStats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	// Wait for the specified interval
	time.Sleep(interval)

	// Get final stats
	finalStats, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var bandwidths []NetworkBandwidth

	// Calculate bandwidth for each interface
	for i, initial := range initialStats {
		if i < len(finalStats) {
			final := finalStats[i]
			if initial.Name == final.Name {
				// Calculate bytes transferred during interval
				uploadBytes := final.BytesSent - initial.BytesSent
				downloadBytes := final.BytesRecv - initial.BytesRecv

				// Convert to Mbps
				seconds := interval.Seconds()
				uploadMbps := float64(uploadBytes) * 8 / seconds / 1000000
				downloadMbps := float64(downloadBytes) * 8 / seconds / 1000000

				bandwidth := NetworkBandwidth{
					Interface:    initial.Name,
					UploadMbps:   uploadMbps,
					DownloadMbps: downloadMbps,
				}
				bandwidths = append(bandwidths, bandwidth)
			}
		}
	}

	return bandwidths, nil
}

// FormatNetworkSpeed formats network speed in human-readable format
func FormatNetworkSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%.1f B/s", bytesPerSecond)
	} else if bytesPerSecond < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", bytesPerSecond/1024)
	} else if bytesPerSecond < 1024*1024*1024 {
		return fmt.Sprintf("%.1f MB/s", bytesPerSecond/(1024*1024))
	} else {
		return fmt.Sprintf("%.1f GB/s", bytesPerSecond/(1024*1024*1024))
	}
}
