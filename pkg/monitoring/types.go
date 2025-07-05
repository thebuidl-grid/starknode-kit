package monitoring

import (
	"time"

	"github.com/rivo/tview"
)

type MonitorApp struct {
	App  *tview.Application
	Grid *tview.Grid

	// Main dashboard panels matching JavaScript components exactly
	ExecutionLogBox *tview.TextView // executionLog.js - Reth logs
	ConsensusLogBox *tview.TextView // consensusLog.js - Lighthouse logs
	StatusBox       *tview.TextView // statusBox.js - Chain status
	ChainInfoBox    *tview.TextView // chainInfoBox.js - Block info with ETH prices
	SystemStatsBox  *tview.TextView // systemStatsGauge.js - Memory/Storage/CPU temp gauges
	RPCInfoBox      *tview.TextView // rpcInfoBox.js - RPC connection info
	StatusBar       *tview.TextView // Status bar for displaying current status

	// Legacy panels (for backward compatibility during transition)
	SystemBox       *tview.TextView
	ClientsBox      *tview.TextView
	LogsBox         *tview.TextView
	PeersBox        *tview.TextView
	ChainBox        *tview.TextView
	NetworkStatsBox *tview.TextView

	// Graph panels
	CPUGraphBox     *tview.TextView
	NetworkGraphBox *tview.TextView
	DiskGraphBox    *tview.TextView

	// Update channels matching the new component structure
	ExecutionLogChan chan string
	ConsensusLogChan chan string
	StatusChan       chan string
	ChainInfoChan    chan string
	SystemStatsChan  chan string
	RPCInfoChan      chan string

	// Legacy channels (for backward compatibility)
	SystemChan  chan string
	ClientsChan chan string
	LogsChan    chan string
	PeersChan   chan string
	ChainChan   chan string
	GraphsChan  chan string
	StatsChan   chan string

	// Control
	StopChan   chan bool
	UpdateRate time.Duration

	// UI state
	darkTheme bool
	paused    bool

	// Graph data storage
	CPUHistory     []float64
	NetworkHistory []NetworkPoint
	DiskHistory    []float64
}

type NetworkPoint struct {
	Upload   float64
	Download float64
	Time     time.Time
}
