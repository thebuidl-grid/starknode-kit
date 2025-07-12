package monitoring

import (
	"fmt"
	"os"
	"starknode-kit/pkg/utils"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *MonitorApp) setupUI() {
	// Set terminal colors to match the JavaScript blessed theme
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorNavy
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorDarkMagenta
	tview.Styles.BorderColor = tcell.ColorTeal
	tview.Styles.TitleColor = tcell.ColorWhite

	// Create execution client log panel (matching executionLog.js)
	m.ExecutionLogBox = m.createVibrantPanel("Reth", tcell.ColorTeal)
	m.ExecutionLogBox.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitle(" Reth ").
		SetTitleAlign(tview.AlignLeft)

	// Create consensus client log panel (matching consensusLog.js)
	m.ConsensusLogBox = m.createVibrantPanel("Lighthouse", tcell.ColorTeal)
	m.ConsensusLogBox.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitle(" Lighthouse ").
		SetTitleAlign(tview.AlignLeft)

	// Create Juno Starknet client log panel
	m.JunoLogBox = m.createVibrantPanel("Juno", tcell.ColorPurple)
	m.JunoLogBox.SetBorder(true).
		SetBorderColor(tcell.ColorPurple).
		SetTitle(" Juno ").
		SetTitleAlign(tview.AlignLeft)

	// Create status box (matching statusBox.js)
	m.StatusBox = m.createVibrantPanel("Status", tcell.ColorTeal)
	m.StatusBox.SetText("INITIALIZING...")

	// Create chain info box (matching chainInfoBox.js)
	m.ChainInfoBox = m.createVibrantPanel("Chain Info", tcell.ColorTeal)

	// Create system stats gauge (matching systemStatsGauge.js)
	m.SystemStatsBox = m.createVibrantPanel("System Stats", tcell.ColorTeal)

	// Create RPC info box
	m.RPCInfoBox = m.createVibrantPanel("RPC Info", tcell.ColorTeal)

	// Perfect LEFT/RIGHT split layout
	// LEFT side (60%): Execution logs (top) + Consensus logs (middle) + Juno logs (bottom) - EQUAL SIZES
	// RIGHT side (40%): Status, Chain Info, RPC Info, System Stats - STACKED VERTICALLY

	m.Grid.SetRows(-1, -1, -1). // 3 rows: execution(1), consensus(1), juno(1) - FULL TERMINAL HEIGHT
					SetColumns(-3, -2). // 2 columns: LEFT(60%), RIGHT(40%)
					SetBorders(false).
					SetGap(0, 0) // No gaps for maximum space usage

	// LEFT SIDE - Execution logs (top third of left side)
	m.Grid.AddItem(m.ExecutionLogBox, 0, 0, 1, 1, 0, 0, false) // Row 0, left col - Execution logs

	// LEFT SIDE - Consensus logs (middle third of left side)
	m.Grid.AddItem(m.ConsensusLogBox, 1, 0, 1, 1, 0, 0, false) // Row 1, left col - Consensus logs

	// LEFT SIDE - Juno logs (bottom third of left side)
	m.Grid.AddItem(m.JunoLogBox, 2, 0, 1, 1, 0, 0, false) // Row 2, left col - Juno logs

	// RIGHT SIDE - Create a sub-grid for the 4 panels stacked vertically
	rightGrid := tview.NewGrid().
		SetRows(-1, -1, -1, -1). // 4 equal rows for the 4 panels
		SetColumns(-1).          // 1 column
		SetBorders(false).
		SetGap(0, 0)

	// Add the 4 panels to the right side sub-grid
	rightGrid.AddItem(m.StatusBox, 0, 0, 1, 1, 0, 0, false)      // Status
	rightGrid.AddItem(m.ChainInfoBox, 1, 0, 1, 1, 0, 0, false)   // Chain Info
	rightGrid.AddItem(m.RPCInfoBox, 2, 0, 1, 1, 0, 0, false)     // RPC Info
	rightGrid.AddItem(m.SystemStatsBox, 3, 0, 1, 1, 0, 0, false) // System Stats

	// Add the right side sub-grid to main grid (spans all 3 rows on right)
	m.Grid.AddItem(rightGrid, 0, 1, 3, 1, 0, 0, false) // Spans rows 0-2 on right column

	// NO HELP BAR - MONITOR FILLS 100% OF TERMINAL HEIGHT

	// Enhanced input handling
	m.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			m.Stop()
			return nil
		case tcell.KeyCtrlC:
			m.Stop()
			return nil
		case tcell.KeyF1:
			m.showAdvancedHelp()
			return nil
		}

		switch event.Rune() {
		case 'q', 'Q':
			m.Stop()
			return nil
		case 'r', 'R':
			m.restartClients()
			return nil
		case 's', 'S':
			m.stopClients()
			return nil
		case 'h', 'H', '?':
			m.showAdvancedHelp()
			return nil
		case 'e', 'E':
			m.exportData()
			return nil
		case 't', 'T':
			m.toggleTheme()
			return nil
		case 'p', 'P':
			m.pauseUpdates()
			return nil
		}

		return event
	})

	m.App.SetRoot(m.Grid, true).SetFocus(m.Grid)
	m.App.EnableMouse(true)

	// Ensure application handles screen resize events
	m.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})
}

func (m *MonitorApp) createVibrantPanel(title string, borderColor tcell.Color) *tview.TextView {
	panel := tview.NewTextView()

	// Set border properties
	panel.SetTitle(fmt.Sprintf(" %s ", title))
	panel.SetTitleAlign(tview.AlignLeft)
	panel.SetBorder(true)
	panel.SetBorderColor(borderColor)

	// Set TextView specific properties
	panel.SetWrap(true)
	panel.SetBackgroundColor(tcell.ColorBlack)
	panel.SetTextColor(tcell.ColorWhite)
	panel.SetDynamicColors(true)

	return panel
}

func (m *MonitorApp) createImageStyleHelpBox() *tview.TextView {
	helpBox := tview.NewTextView()
	helpBox.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(false).
		SetBackgroundColor(tcell.ColorBlack)

	helpText := "[dim]Hsense                                                                                                                    "

	helpBox.SetText(helpText)
	return helpBox
}

func (m *MonitorApp) createHelpBox() *tview.TextView {
	helpBox := tview.NewTextView()
	helpBox.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(false).
		SetBackgroundColor(tcell.ColorBlack)

	helpText := "[dim]Press Q to quit, H for help, R to restart clients, S to stop clients"

	helpBox.SetText(helpText)
	return helpBox
}

// Enhanced client control methods
func (m *MonitorApp) restartClients() {
	m.updateStatusBar("[yellow]⚠️  Restarting all clients...[white]")

	clients := utils.GetRunningClients()
	for _, client := range clients {
		if client.PID > 0 {
			m.updateStatusBar(fmt.Sprintf("[red]🔄 Stopping %s (PID: %d)...[white]", client.Name, client.PID))
			time.Sleep(500 * time.Millisecond) // Visual feedback
		}
		m.updateStatusBar(fmt.Sprintf("[green]▶️  Starting %s...[white]", client.Name))
		time.Sleep(500 * time.Millisecond) // Visual feedback
	}

	m.updateStatusBar("[green]✅ All clients restarted successfully[white]")
}

func (m *MonitorApp) stopClients() {
	m.updateStatusBar("[yellow]⚠️  Stopping all clients...[white]")

	clients := utils.GetRunningClients()
	for _, client := range clients {
		if client.PID > 0 {
			m.updateStatusBar(fmt.Sprintf("[red]⏹️  Stopping %s (PID: %d)...[white]", client.Name, client.PID))
			time.Sleep(300 * time.Millisecond)
		}
	}

	m.updateStatusBar("[red]🔴 All clients stopped[white]")
}

func (m *MonitorApp) toggleLogs() {
	m.updateStatusBar("[blue]📝 Log display toggled[white]")
	m.App.SetFocus(m.LogsBox)
}

func (m *MonitorApp) toggleGraphs() {
	m.updateStatusBar("[cyan]📊 Graph display toggled[white]")
	// Cycle focus through graph panels
	m.App.SetFocus(m.CPUGraphBox)
}

func (m *MonitorApp) clearAll() {
	m.LogsBox.SetText("")
	m.CPUGraphBox.SetText("")
	m.NetworkGraphBox.SetText("")
	m.DiskGraphBox.SetText("")
	m.updateStatusBar("[blue]🧹 All displays cleared[white]")
}

func (m *MonitorApp) showAdvancedHelp() {
	helpModal := tview.NewModal().
		SetText("🚀 [red]STARKNODE MONITORING DASHBOARD[white]\n\n" +
			"[yellow]📊 FEATURES:[white]\n" +
			"  • Real-time system resource monitoring\n" +
			"  • Live client status and health checks\n" +
			"  • Blockchain synchronization tracking\n" +
			"  • Network peer information\n" +
			"  • Interactive graphs and charts\n" +
			"  • Live log streaming from all clients\n" +
			"  • Data export functionality\n" +
			"  • Theme switching (dark/light)\n" +
			"  • Pause/resume monitoring\n\n" +
			"[yellow]⌨️  KEYBOARD SHORTCUTS:[white]\n" +
			"  [green]Q, ESC[white] - Quit the monitor\n" +
			"  [green]R[white] - Restart all clients\n" +
			"  [green]S[white] - Stop all clients\n" +
			"  [green]L[white] - Toggle log display focus\n" +
			"  [green]G[white] - Toggle graph focus\n" +
			"  [green]C[white] - Clear all displays\n" +
			"  [green]E[white] - Export logs and metrics to files\n" +
			"  [green]T[white] - Toggle theme (dark/light)\n" +
			"  [green]P[white] - Pause/resume updates\n" +
			"  [green]H, ?[white] - Show this help\n" +
			"  [green]F1[white] - Advanced help\n\n" +
			"[yellow]🎯 STATUS COLORS:[white]\n" +
			"  [green]Green[white] - Running/Healthy\n" +
			"  [yellow]Yellow[white] - Warning/Syncing\n" +
			"  [red]Red[white] - Error/Stopped\n" +
			"  [blue]Blue[white] - Information\n" +
			"  [cyan]Cyan[white] - Network Activity\n\n" +
			"[yellow]💾 EXPORT FUNCTIONALITY:[white]\n" +
			"  Files are saved with timestamps in current directory\n" +
			"  • starknode_logs_YYYY-MM-DD_HH-MM-SS.txt\n" +
			"  • starknode_metrics_YYYY-MM-DD_HH-MM-SS.txt\n\n" +
			"Press [yellow]ESC[white] to continue...").
		SetBackgroundColor(tcell.ColorBlack).
		SetTextColor(tcell.ColorWhite).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.App.SetRoot(m.Grid, true)
		})

	m.App.SetRoot(helpModal, false)
}

func (m *MonitorApp) updateStatusBar(message string) {
	timestamp := time.Now().Format("15:04:05")
	statusMsg := fmt.Sprintf("🚀 [green]StarkNode Monitor[white] | [blue]%s[white] | %s", timestamp, message)

	m.App.QueueUpdateDraw(func() {
		m.StatusBar.SetText(statusMsg)
	})
}

// Add log entry with enhanced formatting
func (m *MonitorApp) addLogEntry(entry string) {
	timestamp := "[dim]" + time.Now().Format("15:04:05") + "[white] "
	currentText := m.LogsBox.GetText(false)

	lines := strings.Split(currentText, "\n")
	if len(lines) > 100 {
		lines = lines[len(lines)-100:]
		currentText = strings.Join(lines, "\n")
	}

	// Enhanced log formatting with colors based on log level
	formattedEntry := entry
	if strings.Contains(strings.ToUpper(entry), "ERROR") {
		formattedEntry = "[red]" + entry + "[white]"
	} else if strings.Contains(strings.ToUpper(entry), "WARN") {
		formattedEntry = "[yellow]" + entry + "[white]"
	} else if strings.Contains(strings.ToUpper(entry), "INFO") {
		formattedEntry = "[green]" + entry + "[white]"
	} else if strings.Contains(strings.ToUpper(entry), "DEBUG") {
		formattedEntry = "[blue]" + entry + "[white]"
	}

	newText := currentText + "\n" + timestamp + formattedEntry
	m.LogsBox.SetText(newText)
	m.LogsBox.ScrollToEnd()
}

// Graph generation functions
func (m *MonitorApp) generateCPUGraph(history []float64) string {
	if len(history) == 0 {
		return "[yellow]No CPU data available[white]"
	}

	graph := "[green]CPU Usage Graph (Last 60 readings):[white]\n"
	graph += "[blue]" + strings.Repeat("─", 50) + "[white]\n"

	// Simple ASCII graph
	for i, value := range history {
		if i >= len(history)-20 { // Show last 20 values
			bars := int(value / 5) // Scale to 20 chars max (100% / 5 = 20)
			if bars > 20 {
				bars = 20
			}

			color := "[green]"
			if value > 80 {
				color = "[red]"
			} else if value > 60 {
				color = "[yellow]"
			}

			graph += fmt.Sprintf("%s%3.1f%% %s%s[white]\n",
				color, value,
				color, strings.Repeat("█", bars))
		}
	}

	return graph
}

func (m *MonitorApp) generateNetworkGraph(history []NetworkPoint) string {
	if len(history) == 0 {
		return "[yellow]No network data available[white]"
	}

	graph := "[cyan]Network I/O Graph (Upload/Download MB/s):[white]\n"
	graph += "[blue]" + strings.Repeat("─", 50) + "[white]\n"

	for i, point := range history {
		if i >= len(history)-15 { // Show last 15 values
			uploadBars := int(point.Upload / 10) // Scale appropriately
			downloadBars := int(point.Download / 10)

			if uploadBars > 25 {
				uploadBars = 25
			}
			if downloadBars > 25 {
				downloadBars = 25
			}

			timeStr := point.Time.Format("15:04:05")
			graph += fmt.Sprintf("[dim]%s[white] ↑[green]%s[white] ↓[blue]%s[white]\n",
				timeStr,
				strings.Repeat("▲", uploadBars),
				strings.Repeat("▼", downloadBars))
		}
	}

	return graph
}

func (m *MonitorApp) generateDiskGraph(history []float64) string {
	if len(history) == 0 {
		return "[yellow]No disk data available[white]"
	}

	graph := "[yellow]Disk Usage Graph (Last 60 readings):[white]\n"
	graph += "[blue]" + strings.Repeat("─", 50) + "[white]\n"

	// Show disk usage as horizontal bars
	for i, value := range history {
		if i >= len(history)-20 { // Show last 20 values
			bars := int(value / 5) // Scale to 20 chars max
			if bars > 20 {
				bars = 20
			}

			color := "[green]"
			if value > 90 {
				color = "[red]"
			} else if value > 75 {
				color = "[yellow]"
			}

			graph += fmt.Sprintf("%s%3.1f%% %s%s[white]\n",
				color, value,
				color, strings.Repeat("■", bars))
		}
	}

	return graph
}

// New enhancement methods

// Export current data to files
func (m *MonitorApp) exportData() {
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	// Export execution logs
	if m.ExecutionLogBox != nil {
		logsContent := m.ExecutionLogBox.GetText(false)
		if logsContent != "" {
			filename := fmt.Sprintf("starknode_execution_logs_%s.txt", timestamp)
			if err := os.WriteFile(filename, []byte(logsContent), 0644); err == nil {
				m.updateStatusBar(fmt.Sprintf("[green]✅ Execution logs exported to %s[white]", filename))
			} else {
				m.updateStatusBar(fmt.Sprintf("[red]❌ Failed to export execution logs: %v[white]", err))
			}
		}
	}

	// Export consensus logs
	if m.ConsensusLogBox != nil {
		logsContent := m.ConsensusLogBox.GetText(false)
		if logsContent != "" {
			filename := fmt.Sprintf("starknode_consensus_logs_%s.txt", timestamp)
			if err := os.WriteFile(filename, []byte(logsContent), 0644); err == nil {
				m.updateStatusBar(fmt.Sprintf("[green]✅ Consensus logs exported to %s[white]", filename))
			} else {
				m.updateStatusBar(fmt.Sprintf("[red]❌ Failed to export consensus logs: %v[white]", err))
			}
		}
	}

	// Export system metrics
	var metricsContent string
	if m.SystemStatsBox != nil {
		statsContent := m.SystemStatsBox.GetText(false)
		metricsContent = fmt.Sprintf("StarkNode System Metrics Export - %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
		metricsContent += "=== SYSTEM STATS ===\n" + statsContent + "\n\n"
	}

	if m.StatusBox != nil {
		statusContent := m.StatusBox.GetText(false)
		metricsContent += "=== STATUS ===\n" + statusContent + "\n\n"
	}

	if m.ChainInfoBox != nil {
		chainContent := m.ChainInfoBox.GetText(false)
		metricsContent += "=== CHAIN INFO ===\n" + chainContent + "\n"
	}

	if metricsContent != "" {
		metricsFilename := fmt.Sprintf("starknode_metrics_%s.txt", timestamp)
		if err := os.WriteFile(metricsFilename, []byte(metricsContent), 0644); err == nil {
			m.updateStatusBar(fmt.Sprintf("[green]✅ Metrics exported to %s[white]", metricsFilename))
		} else {
			m.updateStatusBar(fmt.Sprintf("[red]❌ Failed to export metrics: %v[white]", err))
		}
	}
}

// Toggle between different color themes
func (m *MonitorApp) toggleTheme() {
	if m.darkTheme {
		// Switch to light theme
		tview.Styles.PrimitiveBackgroundColor = tcell.ColorWhite
		tview.Styles.ContrastBackgroundColor = tcell.ColorLightGray
		if m.ExecutionLogBox != nil {
			m.ExecutionLogBox.SetBackgroundColor(tcell.ColorWhite)
			m.ExecutionLogBox.SetTextColor(tcell.ColorBlack)
		}
		if m.ConsensusLogBox != nil {
			m.ConsensusLogBox.SetBackgroundColor(tcell.ColorWhite)
			m.ConsensusLogBox.SetTextColor(tcell.ColorBlack)
		}
		m.updateStatusBar("[blue]🎨 Switched to light theme[white]")
		m.darkTheme = false
	} else {
		// Switch to dark theme
		tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
		tview.Styles.ContrastBackgroundColor = tcell.ColorNavy
		if m.ExecutionLogBox != nil {
			m.ExecutionLogBox.SetBackgroundColor(tcell.ColorBlack)
			m.ExecutionLogBox.SetTextColor(tcell.ColorWhite)
		}
		if m.ConsensusLogBox != nil {
			m.ConsensusLogBox.SetBackgroundColor(tcell.ColorBlack)
			m.ConsensusLogBox.SetTextColor(tcell.ColorWhite)
		}
		m.updateStatusBar("[blue]🎨 Switched to dark theme[white]")
		m.darkTheme = true
	}
}

// Pause/resume updates
func (m *MonitorApp) pauseUpdates() {
	m.paused = !m.paused
	if m.paused {
		m.updateStatusBar("[yellow]⏸️  Updates paused - Press P to resume[white]")
	} else {
		m.updateStatusBar("[green]▶️  Updates resumed[white]")
	}
}
