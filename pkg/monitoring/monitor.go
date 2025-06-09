package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/rivo/tview"
)

func NewMonitorApp() *MonitorApp {
	app := &MonitorApp{
		App:         tview.NewApplication(),
		Grid:        tview.NewGrid(),
		UpdateRate:  2 * time.Second,
		SystemChan:  make(chan string, 10),
		ClientsChan: make(chan string, 10),
		LogsChan:    make(chan string, 100),
		PeersChan:   make(chan string, 10),
		ChainChan:   make(chan string, 10),
		StopChan:    make(chan bool, 1),
	}

	app.setupUI()
	return app
}

func (m *MonitorApp) Start(ctx context.Context) error {
	// Start update goroutines
	go m.updateSystemStats(ctx)
	go m.updateClientStatus(ctx)
	go m.updateChainInfo(ctx)
	go m.updatePeerInfo(ctx)
	go m.updateLogInfo(ctx)

	// Start UI update handler
	go m.handleUpdates(ctx)

	// Initialize status message
	go func() {
		m.App.QueueUpdateDraw(func() {
			timestamp := time.Now().Format("15:04:05")
			statusMsg := fmt.Sprintf("[dim]%s[white] [green]Monitor started[white] - Press ? for help", timestamp)
			m.StatusBox.SetText(statusMsg)
		})
	}()

	// Run the application
	return m.App.Run()
}

func (m *MonitorApp) Stop() {
	m.StopChan <- true
	m.App.Stop()
}
