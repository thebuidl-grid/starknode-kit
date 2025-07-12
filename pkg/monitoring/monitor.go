package monitoring

import (
	"context"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func NewMonitorApp() *MonitorApp {
	app := &MonitorApp{
		App:        tview.NewApplication(),
		Grid:       tview.NewGrid(),
		StatusBar:  tview.NewTextView(),
		UpdateRate: 2 * time.Second,

		// New channels matching JavaScript components
		ExecutionLogChan: make(chan string, 100),
		ConsensusLogChan: make(chan string, 100),
		JunoLogChan:      make(chan string, 100),
		StatusChan:       make(chan string, 10),
		ChainInfoChan:    make(chan string, 10),
		SystemStatsChan:  make(chan string, 10),
		RPCInfoChan:      make(chan string, 10),

		// Legacy channels for backward compatibility
		SystemChan:  make(chan string, 10),
		ClientsChan: make(chan string, 10),
		LogsChan:    make(chan string, 100),
		PeersChan:   make(chan string, 10),
		ChainChan:   make(chan string, 10),
		GraphsChan:  make(chan string, 10),
		StatsChan:   make(chan string, 10),
		StopChan:    make(chan bool, 1),

		// UI state
		darkTheme: true, // Start with dark theme
		paused:    false,

		// Initialize graph data
		CPUHistory:     make([]float64, 0, 60), // Store 60 data points
		NetworkHistory: make([]NetworkPoint, 0, 60),
		DiskHistory:    make([]float64, 0, 60),
	}

	app.setupUI()
	return app
}

func (m *MonitorApp) Start(ctx context.Context) error {
	// Start new update goroutines matching JavaScript components exactly
	go m.updateExecutionLogs(ctx)    // executionLog.js equivalent
	go m.updateConsensusLogs(ctx)    // consensusLog.js equivalent
	go m.updateJunoLogs(ctx)         // junoLog.js equivalent (Starknet client)
	go m.updateStatusBox(ctx)        // statusBox.js equivalent
	go m.updateChainInfoBox(ctx)     // chainInfoBox.js equivalent
	go m.updateSystemStatsGauge(ctx) // systemStatsGauge.js equivalent
	go m.updateRPCInfo(ctx)          // RPC info component
	// Removed: go m.updateBandwidthGauge(ctx)   // Bandwidth component removed
	// Removed: go m.updatePeerCountGauge(ctx)   // Peer count component removed

	// Start legacy update goroutines for backward compatibility
	go m.updateSystemStats(ctx)
	go m.updateClientStatus(ctx)
	go m.updateChainInfo(ctx)
	go m.updatePeerInfo(ctx)
	go m.updateLogInfo(ctx)
	go m.updateGraphs(ctx)

	// Start UI update handler
	go m.handleUpdates(ctx)

	// Run the application
	return m.App.Run()
}

func (m *MonitorApp) handleUpdates(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case text := <-m.ExecutionLogChan:
			m.App.QueueUpdateDraw(func() {
				m.ExecutionLogBox.SetText(text)
				m.ExecutionLogBox.ScrollToEnd()
			})
		case text := <-m.ConsensusLogChan:
			m.App.QueueUpdateDraw(func() {
				m.ConsensusLogBox.SetText(text)
				m.ConsensusLogBox.ScrollToEnd()
			})
		case text := <-m.JunoLogChan:
			m.App.QueueUpdateDraw(func() {
				m.JunoLogBox.SetText(text)
				m.JunoLogBox.ScrollToEnd()
			})
		case text := <-m.StatusChan:
			m.App.QueueUpdateDraw(func() {
				m.StatusBox.SetText(text)
			})
		case text := <-m.ChainInfoChan:
			m.App.QueueUpdateDraw(func() {
				m.ChainInfoBox.SetText(text)
			})
		case text := <-m.SystemStatsChan:
			m.App.QueueUpdateDraw(func() {
				m.SystemStatsBox.SetText(text)
			})
		case text := <-m.RPCInfoChan:
			m.App.QueueUpdateDraw(func() {
				m.RPCInfoBox.SetText(text)
			})
		// Removed bandwidth and peer count channel handlers
		// Legacy channel handlers for backward compatibility
		case text := <-m.SystemChan:
			m.App.QueueUpdateDraw(func() {
				if m.SystemBox != nil {
					m.SystemBox.SetText(text)
				}
			})
		case text := <-m.ClientsChan:
			m.App.QueueUpdateDraw(func() {
				if m.ClientsBox != nil {
					m.ClientsBox.SetText(text)
				}
			})
		case text := <-m.ChainChan:
			m.App.QueueUpdateDraw(func() {
				if m.ChainBox != nil {
					m.ChainBox.SetText(text)
				}
			})
		case text := <-m.PeersChan:
			m.App.QueueUpdateDraw(func() {
				if m.PeersBox != nil {
					m.PeersBox.SetText(text)
				}
			})
		case text := <-m.LogsChan:
			m.App.QueueUpdateDraw(func() {
				if m.LogsBox != nil {
					m.LogsBox.SetText(text)
				}
			})
		case text := <-m.GraphsChan:
			m.App.QueueUpdateDraw(func() {
				// Parse the prefixed messages for consensus client info
				if strings.HasPrefix(text, "cpu:") && m.CPUGraphBox != nil {
					m.CPUGraphBox.SetText(strings.TrimPrefix(text, "cpu:"))
				} else if strings.HasPrefix(text, "network:") && m.NetworkGraphBox != nil {
					m.NetworkGraphBox.SetText(strings.TrimPrefix(text, "network:"))
				}
			})
		case text := <-m.StatsChan:
			m.App.QueueUpdateDraw(func() {
				if m.NetworkStatsBox != nil {
					m.NetworkStatsBox.SetText(text) // System stats go to NetworkStatsBox (System Stats)
				}
			})
		}
	}
}

func (m *MonitorApp) Stop() {
	m.StopChan <- true
	m.App.Stop()
}
