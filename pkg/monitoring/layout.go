package monitoring

import (
	"context"
	"time"

	"github.com/rivo/tview"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

// rebuildDynamicLayout rebuilds the entire grid layout based on running clients
func (m *MonitorApp) rebuildDynamicLayout() {
	runningClients := utils.GetRunningClients()

	// Determine which clients are running
	hasExecution := false
	hasConsensus := false
	hasJuno := false
	hasValidator := false

	for _, client := range runningClients {
		switch client.Name {
		case "Geth", "Reth":
			hasExecution = true
		case "Lighthouse", "Prysm":
			hasConsensus = true
		case "Juno":
			hasJuno = true
		case "Validator":
			hasValidator = true
		}
	}

	// Count active log panels
	activeLogPanels := 0
	if hasExecution {
		activeLogPanels++
	}
	if hasConsensus {
		activeLogPanels++
	}
	if hasJuno {
		activeLogPanels++
	}
	if hasValidator {
		activeLogPanels++
	}

	// Clear the grid
	m.Grid.Clear()

	// If no clients are running, show "No Clients Running" message
	if activeLogPanels == 0 {
		m.Grid.SetRows(-1).
			SetColumns(-1).
			SetBorders(false)
		m.Grid.AddItem(m.NoClientsBox, 0, 0, 1, 1, 0, 0, false)
		return
	}

	// Create dynamic row configuration based on number of active clients
	rows := make([]int, activeLogPanels)
	for i := range rows {
		rows[i] = -1 // Equal height for all rows
	}

	m.Grid.SetRows(rows...).
		SetColumns(-3, -2). // LEFT(60%), RIGHT(40%)
		SetBorders(false).
		SetGap(0, 0)

	// Add active log panels to the left side
	currentRow := 0
	if hasExecution {
		m.Grid.AddItem(m.ExecutionLogBox, currentRow, 0, 1, 1, 0, 0, false)
		currentRow++
	}
	if hasConsensus {
		m.Grid.AddItem(m.ConsensusLogBox, currentRow, 0, 1, 1, 0, 0, false)
		currentRow++
	}
	if hasJuno {
		m.Grid.AddItem(m.JunoLogBox, currentRow, 0, 1, 1, 0, 0, false)
		currentRow++
	}
	if hasValidator {
		m.Grid.AddItem(m.ValidatorLogBox, currentRow, 0, 1, 1, 0, 0, false)
		currentRow++
	}

	// RIGHT SIDE - Create sub-grid for info panels (5 rows total)
	rightGrid := tview.NewGrid().
		SetRows(-1, -1, -1, -1, -1). // 5 equal rows
		SetColumns(-1).              // Single column
		SetBorders(false).
		SetGap(0, 0)

	// Create status grid for ETH and Starknet status side by side
	statusGrid := tview.NewGrid().
		SetRows(-1).        // Single row
		SetColumns(-1, -1). // 2 equal columns for ETH and Starknet
		SetBorders(false).
		SetGap(1, 0) // Small gap between status panels

	// Add ETH and Starknet status to the status grid
	statusGrid.AddItem(m.StatusBox, 0, 0, 1, 1, 0, 0, false)         // ETH Status (left)
	statusGrid.AddItem(m.StarknetStatusBox, 0, 1, 1, 1, 0, 0, false) // Starknet Status (right)

	// Add all panels to the right side sub-grid
	rightGrid.AddItem(m.NetworkBox, 0, 0, 1, 1, 0, 0, false)     // Row 0: Network
	rightGrid.AddItem(statusGrid, 1, 0, 1, 1, 0, 0, false)       // Row 1: Status grid (ETH + Starknet)
	rightGrid.AddItem(m.ChainInfoBox, 2, 0, 1, 1, 0, 0, false)   // Row 2: Chain Info
	rightGrid.AddItem(m.RPCInfoBox, 3, 0, 1, 1, 0, 0, false)     // Row 3: RPC Info
	rightGrid.AddItem(m.SystemStatsBox, 4, 0, 1, 1, 0, 0, false) // Row 4: System Stats

	// Add the right side sub-grid to main grid (spans all rows on right)
	m.Grid.AddItem(rightGrid, 0, 1, activeLogPanels, 1, 0, 0, false)
}

// updateLayoutDynamically periodically checks for running clients and updates the layout
func (m *MonitorApp) updateLayoutDynamically(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds
	defer ticker.Stop()

	previousState := ""

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case <-ticker.C:
			if m.paused {
				continue
			}

			// Get current running clients
			runningClients := utils.GetRunningClients()

			// Create a state signature based on running clients
			currentState := ""
			for _, client := range runningClients {
				currentState += client.Name + ","
			}

			// Only rebuild layout if the state has changed
			if currentState != previousState {
				m.App.QueueUpdateDraw(func() {
					m.rebuildDynamicLayout()
				})
				previousState = currentState
			}
		}
	}
}
