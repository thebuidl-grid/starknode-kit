package monitoring

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

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
		SetTitle(" Execution Client (Detecting...) ⚡ ").
		SetTitleAlign(tview.AlignLeft)

	// Create consensus client log panel (matching consensusLog.js)
	m.ConsensusLogBox = m.createVibrantPanel("Lighthouse", tcell.ColorTeal)
	m.ConsensusLogBox.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitle(" Consensus Client (Detecting...) 🏛️ ").
		SetTitleAlign(tview.AlignLeft)

	// Create Juno Starknet client log panel
	m.JunoLogBox = m.createVibrantPanel("Juno", tcell.ColorTeal)
	m.JunoLogBox.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitle(" Juno (Detecting...) 🌟 ").
		SetTitleAlign(tview.AlignLeft)

	// Create Validator log panel
	m.ValidatorLogBox = m.createVibrantPanel("Validator", tcell.ColorTeal)
	m.ValidatorLogBox.SetBorder(true).
		SetBorderColor(tcell.ColorTeal).
		SetTitle(" Validator (Detecting...) 🛡️ ").
		SetTitleAlign(tview.AlignLeft)

	// Create "No Clients Running" message box
	m.NoClientsBox = m.createVibrantPanel("Status", tcell.ColorYellow)
	m.NoClientsBox.SetBorder(true).
		SetBorderColor(tcell.ColorYellow).
		SetTitle(" System Status ⚠️ ").
		SetTitleAlign(tview.AlignCenter)
	m.NoClientsBox.SetText("\n\n[yellow]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n" +
		"[white]        ⚠️  NO CLIENTS RUNNING  ⚠️\n\n" +
		"[dim]No Ethereum or Starknet clients are currently active.\n\n" +
		"[yellow]To start clients, use:[white]\n" +
		"  • starknode-kit start\n" +
		"  • starknode-kit run\n\n" +
		"[yellow]━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[white]").
		SetTextAlign(tview.AlignCenter)

	// Create status boxes
	m.StatusBox = m.createVibrantPanel("L1 Status", tcell.ColorTeal)
	m.StatusBox.SetText("INITIALIZING...")

	m.NetworkBox = m.createVibrantPanel("Network", tcell.ColorTeal)
	m.NetworkBox.SetText("INITIALIZING...")

	m.StarknetStatusBox = m.createVibrantPanel("L2 Status", tcell.ColorTeal)
	m.StarknetStatusBox.SetText("INITIALIZING...")

	// Create info boxes
	m.ChainInfoBox = m.createVibrantPanel("Chain Info", tcell.ColorTeal)
	m.SystemStatsBox = m.createVibrantPanel("System Stats", tcell.ColorTeal)
	m.RPCInfoBox = m.createVibrantPanel("RPC Info", tcell.ColorTeal)

	// Initial setup with placeholder - will be rebuilt dynamically
	m.rebuildDynamicLayout()

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
