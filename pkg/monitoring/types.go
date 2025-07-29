package monitoring

import (
	"time"

	"github.com/rivo/tview"
)

type MonitorApp struct {
	App  *tview.Application
	Grid *tview.Grid

	// Main dashboard panels matching JavaScript components exactly
	ExecutionLogBox   *tview.TextView
	ConsensusLogBox   *tview.TextView
	JunoLogBox        *tview.TextView
	StatusBox         *tview.TextView
	NetworkBox        *tview.TextView
	StarknetStatusBox *tview.TextView
	ChainInfoBox      *tview.TextView
	SystemStatsBox    *tview.TextView
	RPCInfoBox        *tview.TextView
	StatusBar         *tview.TextView

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
	JunoLogChan      chan string
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
