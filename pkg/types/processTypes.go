package types

import "time"

type ProcessInfo struct {
	PID      int           `json:"pid"`
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Uptime   time.Duration `json:"uptime"`
	CPUUsage float64       `json:"cpu_usage"`
	MemUsage uint64        `json:"mem_usage"`
}

// EthereumMetrics holds blockchain metrics
type EthereumMetrics struct {
	CurrentBlock uint64  `json:"current_block"`
	HighestBlock uint64  `json:"highest_block"`
	SyncPercent  float64 `json:"sync_percent"`
	PeerCount    int     `json:"peer_count"`
	IsSyncing    bool    `json:"is_syncing"`
	GasPrice     string  `json:"gas_price"`
	NetworkName  string  `json:"network"`
}
