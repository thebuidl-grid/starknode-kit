package monitoring

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/stats"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

// Format log lines exactly like the JavaScript helperFunctions.js
func formatLogLines(line string) string {
	// Define words which should be highlighted (from JavaScript)
	highlightRules := map[string]string{
		"INFO":                "[green][bold]INFO[white]",
		"WARN":                "[yellow][bold]WARN[white]",
		"ERROR":               "[red][bold]ERROR[white]",
		"updated":             "[yellow][bold]updated[white]",
		"latestProcessedSlot": "[green][bold]latestProcessedSlot[white]",
	}

	// Apply styles to the words
	for word, replacement := range highlightRules {
		line = strings.ReplaceAll(line, word, replacement)
	}

	// Highlight words followed by "=" in green (like number=, hash=, etc.)
	// Simple regex-like replacement for Go
	words := strings.Fields(line)
	for i, word := range words {
		if strings.Contains(word, "=") {
			parts := strings.Split(word, "=")
			if len(parts) == 2 {
				words[i] = fmt.Sprintf("[green][bold]%s[white]=%s", parts[0], parts[1])
			}
		}
		// Highlight words followed by ":" and surrounded by spaces
		if strings.HasSuffix(word, ":") && len(word) > 1 {
			words[i] = fmt.Sprintf("[blue][bold]%s[white]", word)
		}
	}

	line = strings.Join(words, " ")

	// Replace three or more consecutive spaces with two spaces
	for strings.Contains(line, "   ") {
		line = strings.ReplaceAll(line, "   ", "  ")
	}

	return line
}

// Update execution client logs dynamically based on running clients
func (m *MonitorApp) updateExecutionLogs(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Default to Reth logs if no execution client is detected
	elLogs := []string{}

	logIndex := 0
	var logBuffer []string
	var currentClientName string
	var currentLogs []string

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

			// Detect which execution client is running
			runningClients := utils.GetRunningClients()
			var executionClient *types.ClientStatus

			for _, client := range runningClients {
				if client.Name == "Geth" || client.Name == "Reth" {
					executionClient = &client
					break
				}
			}

			// Update panel title and logs based on detected client
			if executionClient != nil {
				if currentClientName != executionClient.Name {
					// Client changed, reset everything
					currentClientName = executionClient.Name
					logIndex = 0
					logBuffer = []string{}

					// Update panel title
					m.App.QueueUpdateDraw(func() {
						if executionClient.Name == "Geth" {
							m.ExecutionLogBox.SetTitle(" Geth ⚙️ ")
						} else {
							m.ExecutionLogBox.SetTitle(" Reth ⚡ ")
						}
					})
					currentLogs = elLogs
				}

				// Try to get real logs first
				realLogs := GetLatestLogs(strings.ToLower(executionClient.Name), 10)
				if len(realLogs) > 0 && realLogs[0] != fmt.Sprintf("No log files found for %s", strings.ToLower(executionClient.Name)) {
					// Use real logs
					content := strings.Join(realLogs, "\n")
					select {
					case m.ExecutionLogChan <- content:
					default:
						// Channel full, skip update
					}
				} else {
					// Fall back to simulated logs
					if len(currentLogs) > 0 {
						currentEntry := currentLogs[logIndex%len(currentLogs)]

						// Update timestamps and block numbers for realism
						if !strings.Contains(currentEntry, "v1.") {
							// Update timestamp for Geth format
							if executionClient.Name == "Geth" {
								timestamp := time.Now().Format("01-02|15:04:05.000")
								if strings.Contains(currentEntry, "06-15|") {
									currentEntry = strings.Replace(currentEntry, "06-15|13:5", timestamp[:len(timestamp)-4]+":5", 1)
								}
							} else {
								// Update timestamp for Reth format
								timestamp := time.Now().Format("2006-01-02T15:04:05.000000Z")
								if strings.Contains(currentEntry, "2025-06-15T") {
									parts := strings.Fields(currentEntry)
									if len(parts) > 0 {
										parts[0] = timestamp
										currentEntry = strings.Join(parts, " ")
									}
								}
							}

						}

						// Format the log line
						formattedLine := formatLogLines(currentEntry)

						// Add to buffer (matching JavaScript ensureBufferFillsWidget logic)
						logBuffer = append(logBuffer, formattedLine)

						// Keep buffer size manageable (matching JavaScript visibleHeight logic)
						if len(logBuffer) > 50 {
							logBuffer = logBuffer[len(logBuffer)-45:] // Keep last 45 lines
						}

						// Update the display
						content := strings.Join(logBuffer, "\n")
						select {
						case m.ExecutionLogChan <- content:
						default:
							// Channel full, skip update
						}

						logIndex++
					}
				}
			} else {
				// No execution client running
				if currentClientName != "None" {
					currentClientName = "None"
					m.App.QueueUpdateDraw(func() {
						m.ExecutionLogBox.SetTitle(" Execution Client (Not Running) ❌ ")
					})

					select {
					case m.ExecutionLogChan <- "[red]No execution client detected.[white]\n[yellow]Start Geth or Reth to see live logs.[white]":
					default:
					}
				}
			}
		}
	}
}

// Update consensus client logs dynamically based on running clients
func (m *MonitorApp) updateConsensusLogs(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	clLogs := []string{}

	logIndex := 0
	var logBuffer []string
	var currentClientName string
	var currentLogs []string

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

			// Detect which consensus client is running
			runningClients := utils.GetRunningClients()
			var consensusClient *types.ClientStatus

			for _, client := range runningClients {
				if client.Name == "Lighthouse" || client.Name == "Prysm" {
					consensusClient = &client
					break
				}
			}

			// Update panel title and logs based on detected client
			if consensusClient != nil {
				if currentClientName != consensusClient.Name {
					// Client changed, reset everything
					currentClientName = consensusClient.Name
					logIndex = 0
					logBuffer = []string{}

					// Update panel title
					m.App.QueueUpdateDraw(func() {
						if consensusClient.Name == "Prysm" {
							m.ConsensusLogBox.SetTitle(" Prysm 🏛️ ")
						} else {
							m.ConsensusLogBox.SetTitle(" Lighthouse 🏛️ ")
						}
					})
					currentLogs = clLogs
				}

				// Try to get real logs first
				realLogs := GetLatestLogs(strings.ToLower(consensusClient.Name), 10)
				if len(realLogs) > 0 && realLogs[0] != fmt.Sprintf("No log files found for %s", strings.ToLower(consensusClient.Name)) {
					// Use real logs
					content := strings.Join(realLogs, "\n")
					select {
					case m.ConsensusLogChan <- content:
					default:
						// Channel full, skip update
					}
				} else {
					// Fall back to simulated logs
					if len(currentLogs) > 0 {
						currentEntry := currentLogs[logIndex%len(currentLogs)]

						// Update timestamps and slot numbers for realism
						if !strings.Contains(currentEntry, "v7.") && !strings.Contains(currentEntry, "v4.") {
							if consensusClient.Name == "Prysm" {
								// Update timestamp for Prysm format
								timestamp := time.Now().Format("2006-01-02T15:04:05Z")
								if strings.Contains(currentEntry, "2024-06-10T") {
									currentEntry = strings.Replace(currentEntry, "2024-06-10T13:5", timestamp[:len(timestamp)-9]+":5", 1)
								}
							} else {
								// Update timestamp for Lighthouse format (maintaining Jun 10 format)
								currentTime := time.Now().Format("Jan 02 15:04:05.000")
								if strings.Contains(currentEntry, "Jun 10 ") {
									// Replace the timestamp part
									parts := strings.Fields(currentEntry)
									if len(parts) >= 3 {
										parts[0] = currentTime[:6] // "Jun 10" part
										parts[1] = currentTime[7:] // Time part
										currentEntry = strings.Join(parts, " ")
									}
								}
							}
						}

						// Format the log line
						formattedLine := formatLogLines(currentEntry)

						// Add to buffer
						logBuffer = append(logBuffer, formattedLine)

						// Keep buffer size manageable
						if len(logBuffer) > 50 {
							logBuffer = logBuffer[len(logBuffer)-45:]
						}

						// Send to consensus log channel
						content := strings.Join(logBuffer, "\n")
						select {
						case m.ConsensusLogChan <- content:
						default:
							// Channel full, skip update
						}

						logIndex++
					}
				}
			} else {
				// No consensus client running
				if currentClientName != "None" {
					currentClientName = "None"
					m.App.QueueUpdateDraw(func() {
						m.ConsensusLogBox.SetTitle(" Consensus Client (Not Running) ❌ ")
					})

					select {
					case m.ConsensusLogChan <- "[red]No consensus client detected.[white]\n[yellow]Start Lighthouse or Prysm to see live logs.[white]":
					default:
					}
				}
			}
		}
	}
}

// Update status box (matching statusBox.js and updateLogic.js synchronizeAndUpdateWidgets)
func (m *MonitorApp) updateStatusBox(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

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

			// Add some dynamic information
			config, _ := utils.LoadConfig()
			currentTime := time.Now()
			ethStatus := utils.GetGethSyncStatus()
			currentStrkBlock := ethStatus.CurrentBlock
			peers := ethStatus.PeersCount
			netowrk := config.Network
			isSyncing := ethStatus.IsSyncing

			l1statusContent := fmt.Sprintf("Block: [yellow]%d[white]\n", currentStrkBlock)
			l1statusContent += fmt.Sprintf("Peers: [green]%d[white]\n", peers)
			l1statusContent += fmt.Sprintf("Syncing: [green]%t[white]\n", isSyncing)

			l2Status := GetJunoMetrics()
			l2statusContent := fmt.Sprintf("Current Block: [yellow]%d[white]\n", l2Status.CurrentBlock)
			l2statusContent += fmt.Sprintf("Syncing: [green]%t[white]\n", l2Status.IsSyncing)

			networkChanContent := fmt.Sprintf("Network: %s\ntime: %s", netowrk, currentTime.Format("15:04:05"))

			// Send to status channel
			select {
			case m.StatusChan <- l1statusContent:
			case m.JunoStatusChan <- l2statusContent:
			case m.NetworkChan <- networkChanContent:
			default:
				// Channel full, skip update
			}

		}
	}
}

// Update chain info box (matching chainInfoBox.js populateChainInfoBox)
func (m *MonitorApp) updateChainInfoBox(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

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

			content := ""
			separator := strings.Repeat("-", 25)
			content += separator + "\n"

			metrics := GetEthereumMetrics()
			gasPrice := metrics.GasPrice
			content += fmt.Sprintf("[blue][bold]GAS:[white]   %s\n", gasPrice)
			content += separator

			// Send to chain info channel instead of direct update
			select {
			case m.ChainInfoChan <- content:
			default:
				// Channel full, skip update
			}
		}
	}
}

// Update system stats gauge (matching systemStatsGauge.js)
func (m *MonitorApp) updateSystemStatsGauge(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

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

			systemStats, err := stats.GetSystemStats()
			if err != nil {
				continue
			}

			// Create gauge display like JavaScript version
			gaugeNames := []string{"MEMORY", "STORAGE", "CPU TEMP"}
			gaugeColors := []string{"[magenta]", "[green]", "[blue]"}
			units := []string{"%", "%", "C"}

			values := []float64{
				systemStats.Memory.UsedPercent,
				systemStats.Disk.UsedPercent,
				float64(45 + rand.Intn(20)), // Simulated CPU temp
			}

			content := ""
			boxWidth := 20 // Fixed width for gauges

			for i, value := range values {
				percentComplete := value / 100.0
				if percentComplete > 1 {
					percentComplete = 1
				}

				percentString := fmt.Sprintf("%.0f%s", value, units[i])
				filledBars := int(float64(boxWidth) * percentComplete)

				bar := strings.Repeat("█", filledBars) + strings.Repeat(" ", boxWidth-filledBars)

				content += fmt.Sprintf("%s%s\n[%s] %s\n",
					gaugeColors[i], gaugeNames[i], bar, percentString)
			}

			// Send to system stats channel
			select {
			case m.SystemStatsChan <- content:
			default:
				// Channel full, skip update
			}
		}
	}
}

// Update RPC info box to show connection status and endpoints
func (m *MonitorApp) updateRPCInfo(ctx context.Context) {
	ticker := time.NewTicker(8 * time.Second)
	defer ticker.Stop()

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

			content := "[yellow][bold]RPC STATUS[white]\n"
			content += strings.Repeat("-", 15) + "\n"

			// Execution RPC
			execRPCURL := "http://localhost:8545"
			execStatus, _ := utils.CheckRPCStatus(execRPCURL, "web3_clientVersion")
			content += fmt.Sprintf("[blue]Execution:[white]\n%s\n[dim]%s[white]\n\n", execStatus, execRPCURL)

			// Consensus RPC
			consRPCURL := "http://localhost:5054"
			consStatus, _ := utils.CheckRPCStatus(consRPCURL, "eth/v1/node/health") // Use a valid consensus-layer method
			content += fmt.Sprintf("[blue]Consensus:[white]\n%s\n[dim]%s[white]\n", consStatus, consRPCURL)

			// Send to RPC info channel
			select {
			case m.RPCInfoChan <- content:
			default:
				// Channel full, skip update
			}
		}
	}
}

// Update Juno client logs dynamically based on running client and real logs
func (m *MonitorApp) updateJunoLogs(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Realistic Juno logs matching typical Starknet node output (fallback)
	junoLogs := []string{}

	var logBuffer []string
	logIndex := 0
	var currentClientName string

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

			// Detect if Juno client is running
			runningClients := utils.GetRunningClients()
			var junoClient *types.ClientStatus

			for _, client := range runningClients {
				if client.Name == "Juno" {
					junoClient = &client
					break
				}
			}

			// Update panel title and logs based on detected client
			if junoClient != nil {
				if currentClientName != junoClient.Name {
					// Client changed, reset everything
					currentClientName = junoClient.Name
					logIndex = 0
					logBuffer = []string{}

					// Update panel title to show it's running
					m.App.QueueUpdateDraw(func() {
						m.JunoLogBox.SetTitle(" Juno 🌟 (Running) ")
					})
				}

				// Try to get real logs first from Juno log directory
				realLogs := GetLatestLogs("juno", 10)
				if len(realLogs) > 0 && realLogs[0] != "No log files found for juno" {
					// Use real logs from Juno client
					var formattedRealLogs []string
					for _, logLine := range realLogs {
						if strings.TrimSpace(logLine) != "" {
							formattedLine := formatLogLines(logLine)
							formattedRealLogs = append(formattedRealLogs, formattedLine)
						}
					}

					content := strings.Join(formattedRealLogs, "\n")
					select {
					case m.JunoLogChan <- content:
					default:
						// Channel full, skip update
					}
				} else {
					// Fall back to simulated logs if real logs aren't available
					if logIndex < len(junoLogs) {
						currentEntry := junoLogs[logIndex]

						// Dynamic updates for realistic logs
						if strings.Contains(currentEntry, "block_number=") {
							// Update block numbers progressively
							baseBlock := 650328
							currentBlock := baseBlock + int(time.Now().Unix()%100)
							currentEntry = strings.ReplaceAll(currentEntry, "650328", fmt.Sprintf("%d", currentBlock))
							currentEntry = strings.ReplaceAll(currentEntry, "650329", fmt.Sprintf("%d", currentBlock+1))
							currentEntry = strings.ReplaceAll(currentEntry, "650330", fmt.Sprintf("%d", currentBlock+2))
							currentEntry = strings.ReplaceAll(currentEntry, "650331", fmt.Sprintf("%d", currentBlock+3))
							currentEntry = strings.ReplaceAll(currentEntry, "650335", fmt.Sprintf("%d", currentBlock+7))
						}

						// Update timestamps to current time
						if strings.Contains(currentEntry, "15:32:") {
							now := time.Now()
							timeStr := now.Format("15:04:05")
							// Replace the timestamp part
							parts := strings.Split(currentEntry, "] ")
							if len(parts) >= 2 {
								parts[0] = fmt.Sprintf("INFO [12-07|%s.%03d", timeStr, now.Nanosecond()/1000000)
								currentEntry = strings.Join(parts, "] ")
							}
						}

						// Format the log line
						formattedLine := formatLogLines(currentEntry)

						// Add to buffer
						logBuffer = append(logBuffer, formattedLine)

						// Keep buffer size manageable
						if len(logBuffer) > 50 {
							logBuffer = logBuffer[len(logBuffer)-45:]
						}

						// Send to Juno log channel
						content := strings.Join(logBuffer, "\n")
						select {
						case m.JunoLogChan <- content:
						default:
							// Channel full, skip update
						}

						logIndex++
					} else {
						// Reset to beginning for continuous simulation
						logIndex = 0
					}
				}
			} else {
				// No Juno client running
				if currentClientName != "None" {
					currentClientName = "None"
					m.App.QueueUpdateDraw(func() {
						m.JunoLogBox.SetTitle(" Juno (Not Running) ❌ ")
					})

					select {
					case m.JunoLogChan <- "[red]No Juno client detected.[white]\n[yellow]Start Juno to see live logs.[white]":
					default:
					}
				}
			}
		}
	}
}

// Legacy update methods for backward compatibility

func (m *MonitorApp) updateSystemStats(ctx context.Context) {
	// This is now handled by updateSystemStatsGauge - just forward
	go m.updateSystemStatsGauge(ctx)
}

func (m *MonitorApp) updateClientStatus(ctx context.Context) {
	// Legacy method - no longer needed with new layout
}

func (m *MonitorApp) updateChainInfo(ctx context.Context) {
	// This is now handled by updateChainInfoBox - just forward
	go m.updateChainInfoBox(ctx)
}

func (m *MonitorApp) updatePeerInfo(ctx context.Context) {
	// This is now handled by updatePeerCountGauge - just forward
	// TODO: Implement updatePeerCountGauge method or remove this call
}

func (m *MonitorApp) updateLogInfo(ctx context.Context) {
	// This is now handled by updateExecutionLogs and updateConsensusLogs
	// No need to do anything here
}

func (m *MonitorApp) updateGraphs(ctx context.Context) {
	// Legacy method - graphs are not part of the JavaScript component layout
}

// Helper functions

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatNumberWithCommas(n uint64) string {
	str := fmt.Sprintf("%d", n)
	if len(str) <= 3 {
		return str
	}

	result := ""
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(char)
	}
	return result
}
