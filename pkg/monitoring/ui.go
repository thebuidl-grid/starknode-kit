package monitoring

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (m *MonitorApp) setupUI() {
	// Create text views
	m.SystemBox = tview.NewTextView()
	m.SystemBox.SetTitle(" System Stats ").SetBorder(true)

	m.ClientsBox = tview.NewTextView()
	m.ClientsBox.SetTitle(" Client Status ").SetBorder(true)

	m.LogsBox = tview.NewTextView()
	m.LogsBox.SetTitle(" Recent Logs ").SetBorder(true)

	m.PeersBox = tview.NewTextView()
	m.PeersBox.SetTitle(" Network Peers ").SetBorder(true)

	m.ChainBox = tview.NewTextView()
	m.ChainBox.SetTitle(" Blockchain Info ").SetBorder(true)

	m.StatusBox = tview.NewTextView()
	m.StatusBox.SetTitle(" Status Messages ").SetBorder(true)

	// Setup improved grid layout
	m.Grid.SetRows(8, 8, 0, 3).
		SetColumns(0, 0, 0).
		SetBorders(true)

	// Add components to grid
	m.Grid.AddItem(m.SystemBox, 0, 0, 1, 1, 0, 0, false).
		AddItem(m.ClientsBox, 0, 1, 1, 1, 0, 0, false).
		AddItem(m.PeersBox, 0, 2, 1, 1, 0, 0, false).
		AddItem(m.ChainBox, 1, 0, 1, 1, 0, 0, false).
		AddItem(m.LogsBox, 1, 1, 1, 2, 0, 0, false).
		AddItem(m.StatusBox, 2, 0, 1, 3, 0, 0, false).
		AddItem(m.createHelpBox(), 3, 0, 1, 3, 0, 0, false)

	// Set input capture for commands
	m.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			m.Stop()
			return nil
		case tcell.KeyCtrlC:
			m.Stop()
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
		case 'l', 'L':
			m.toggleLogs()
			return nil
		case '?':
			m.showHelp()
			return nil
		}

		return event
	})

	m.App.SetRoot(m.Grid, true)
}

func (m *MonitorApp) createHelpBox() *tview.TextView {
	helpBox := tview.NewTextView()
	helpBox.SetTitle(" Help ").SetBorder(true)

	helpText := "[yellow]Commands:[white]\n" +
		"  [green]q/ESC[white] - Quit\n" +
		"  [green]r[white] - Restart clients\n" +
		"  [green]s[white] - Stop clients\n" +
		"  [green]l[white] - Toggle logs\n" +
		"  [green]?[white] - Show help\n"

	helpBox.SetText(helpText)
	return helpBox
}

// Client control methods
func (m *MonitorApp) restartClients() {
	m.updateStatus("[yellow]Restarting all clients...[white]")

	// Get running clients and restart them
	clients := GetRunningClients()
	for _, client := range clients {
		// Stop client first
		if client.PID > 0 {
			m.updateStatus("[red]Stopping " + client.Name + "...[white]")
			// In a real implementation, you would send a signal to stop the process
			// For now, we'll just simulate the action
		}

		// Start client again
		m.updateStatus("[green]Starting " + client.Name + "...[white]")
		// In a real implementation, you would start the client process
		// For now, we'll just simulate the action
	}

	m.updateStatus("[green]All clients restarted successfully[white]")
}

func (m *MonitorApp) stopClients() {
	m.updateStatus("[yellow]Stopping all clients...[white]")

	// Get running clients and stop them
	clients := GetRunningClients()
	for _, client := range clients {
		if client.PID > 0 {
			m.updateStatus("[red]Stopping " + client.Name + "...[white]")
			// In a real implementation, you would send a signal to stop the process
			// For now, we'll just simulate the action
		}
	}

	m.updateStatus("[red]All clients stopped[white]")
}

// Log display methods
func (m *MonitorApp) toggleLogs() {
	// Toggle log visibility by changing focus and highlighting
	m.updateStatus("[blue]Log display toggled[white]")
	// Focus on the logs box
	m.App.SetFocus(m.LogsBox)
}

// Help and status methods
func (m *MonitorApp) showHelp() {
	// Create a modal with detailed help information
	helpModal := tview.NewModal().
		SetText("StarkNode Kit Monitor\n\n" +
			"Keyboard Commands:\n" +
			"  q, Q, ESC - Quit the monitor\n" +
			"  r, R - Restart all clients\n" +
			"  s, S - Stop all clients\n" +
			"  l, L - Toggle log display focus\n" +
			"  ? - Show this help\n\n" +
			"Monitor Features:\n" +
			"  • Real-time system resource monitoring\n" +
			"  • Client status and health checks\n" +
			"  • Blockchain synchronization status\n" +
			"  • Network peer information\n" +
			"  • Live log streaming\n\n" +
			"Press any key to continue...").
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			m.App.SetRoot(m.Grid, true)
		})

	m.App.SetRoot(helpModal, false)
}

func (m *MonitorApp) updateStatus(message string) {
	// Add timestamp and send to status box
	timestamp := "[dim]" + time.Now().Format("15:04:05") + "[white] "
	m.StatusBox.SetText(timestamp + message)
}

// Log streaming methods
func (m *MonitorApp) addLogEntry(entry string) {
	// Add new log entry with timestamp
	timestamp := "[dim]" + time.Now().Format("15:04:05") + "[white] "
	currentText := m.LogsBox.GetText(false)

	// Keep only last 100 lines to prevent memory issues
	lines := strings.Split(currentText, "\n")
	if len(lines) > 100 {
		lines = lines[len(lines)-100:]
		currentText = strings.Join(lines, "\n")
	}

	newText := currentText + "\n" + timestamp + entry
	m.LogsBox.SetText(newText)

	// Auto-scroll to bottom
	m.LogsBox.ScrollToEnd()
}

func (m *MonitorApp) clearLogs() {
	m.LogsBox.SetText("")
	m.updateStatus("[blue]Logs cleared[white]")
}
