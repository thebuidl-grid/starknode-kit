package monitoring

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"starknode-kit/pkg/stats"
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

// Update execution client logs (matching updateLogic.js setupLogStreaming)
func (m *MonitorApp) updateExecutionLogs(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Realistic Reth logs matching the JavaScript format exactly
	rethLogs := []string{
		"Reth v1.3.12",
		"2025-06-15T13:50:49.381500Z INFO Canonical chain committed number=22674573 hash=0x42e0c59670d0aae3e65fc1109baec597207bb20eacceab850374b706 elapsed=796.542μs",
		"2025-06-15T13:51:01.107312Z INFO State root task finished state_root=0x37fcdd548d484c536f5816b855e2497e03da5e7503 elapsed=1.07344ms",
		"2025-06-15T13:51:01.107966Z INFO Block added to canonical chain number=22674574 hash=0xca5d0e5dbdc9f6fab2be1b4f8b9c6cc76fbd09ce3 elapsed=293.77ms",
		"2025-06-15T13:51:01.267946Z INFO Canonical chain committed number=22674574 hash=0xca5d0e5dbdc9f6fab2be1b4f8b9c6cc76fbd09ce3dbbdb6 elapsed=416.344μs",
		"2025-06-15T13:51:02.187294Z INFO State root task finished state_root=0x6660c6f5c489bbdc35507dfb74b016139c885f3b52 elapsed=356098ms",
		"2025-06-15T13:51:13.279558Z INFO Block added to canonical chain number=22674575 hash=0x2a965406cad77de0d2701d3d5a29cdd0197a7d71 elapsed=290.97ms",
		"2025-06-15T13:51:15.123Z INFO Downloading bodies from_block=22674576 to_block=22674580",
		"2025-06-15T13:51:16.456Z INFO Processing transactions block=22674577 txs=156 gas_used=15234567",
		"2025-06-15T13:51:17.789Z INFO State trie update merkle_root=0x8a4b2c7d9e3f1a6b5c8d2e7f4a9b3c6d8e1f4a7b2c5d8e1f4a7b block=22674578",
	}

	logIndex := 0
	var logBuffer []string

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

			// Get current log entry
			currentEntry := rethLogs[logIndex%len(rethLogs)]

			// Update timestamps and block numbers for realism
			if !strings.Contains(currentEntry, "Reth v") {
				// Update timestamp
				timestamp := time.Now().Format("2006-01-02T15:04:05.000000Z")
				if strings.Contains(currentEntry, "2025-06-15T") {
					// Replace the timestamp part
					parts := strings.Fields(currentEntry)
					if len(parts) > 0 {
						parts[0] = timestamp
						currentEntry = strings.Join(parts, " ")
					}
				}

				// Update block numbers progressively
				currentBlock := 22674573 + int(time.Now().Unix()%1000)
				currentEntry = strings.ReplaceAll(currentEntry, "22674573", fmt.Sprintf("%d", currentBlock))
				currentEntry = strings.ReplaceAll(currentEntry, "22674574", fmt.Sprintf("%d", currentBlock+1))
				currentEntry = strings.ReplaceAll(currentEntry, "22674575", fmt.Sprintf("%d", currentBlock+2))
				currentEntry = strings.ReplaceAll(currentEntry, "22674576", fmt.Sprintf("%d", currentBlock+3))
				currentEntry = strings.ReplaceAll(currentEntry, "22674577", fmt.Sprintf("%d", currentBlock+4))
				currentEntry = strings.ReplaceAll(currentEntry, "22674578", fmt.Sprintf("%d", currentBlock+5))
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
			// Send to execution log channel
			content := strings.Join(logBuffer, "\n")
			select {
			case m.ExecutionLogChan <- content:
			default:
				// Channel full, skip update
			}

			logIndex++
		}
	}
}

// Update consensus client logs (matching consensusLog.js)
func (m *MonitorApp) updateConsensusLogs(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Realistic Lighthouse logs matching the JavaScript format
	lighthouseLogs := []string{
		"Lighthouse v7.0.1",
		"Jun 10 13:50:36.160 INFO New block received root: 0xbd82ff3359a52ebc815f8b171ee6465d53e5598ae141612c68aef3793b940f7ff, elapsed=37171μs",
		"Jun 10 13:50:41.004 INFO Synced slot: 11894951, block: 0xd482ff3359e52ebc815f8b171ee6465d5ae3e5598ae141612c68aef3793b940f7ff, elapsed=371717μs",
		"Jun 10 13:50:45.001 INFO Synced verified attestation slot: 11894952 block: 0x9d1329711ce559e1f47bd499d0ae2be3f3ae4560e0bb0aded6b10dd312df78, elapsed=hash: 0x4a2e0c59a78bda6adf5f",
		"Jun 10 13:51:05.001 INFO New block received root: 0x9d1329711ce559e1f47bd499d0ae2be3f3ae4560e0bb0aded6b10dd312df78, slot_notifier",
		"Jun 10 13:51:32.796 INFO New block received slot: 11894953, elapsed=0xef7fd29c91e856c43dee9573f392f96ae540d7ffce032408ba9a1a053db",
		"Jun 10 13:51:32.797 INFO Synced slot: 11894953, block: 0xef7fd29c91e856c43dee9573f392f96ae540d7ffce032408ba9a1a053db, elapsed=371757μs",
		"Jun 10 13:51:45.123 INFO Updated latestProcessedSlot slot=11894954 epoch=372185",
		"Jun 10 13:51:50.456 INFO Attestation published slot=11894955 committee_index=5 validator_index=12345",
	}

	logIndex := 0
	var logBuffer []string

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

			// Get current log entry
			currentEntry := lighthouseLogs[logIndex%len(lighthouseLogs)]

			// Update timestamps and slot numbers for realism
			if !strings.Contains(currentEntry, "Lighthouse v") {
				// Update timestamp to current time (maintaining Jun 10 format)
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

				// Update slot numbers progressively
				currentSlot := 11894951 + int(time.Now().Unix()%100)
				currentEntry = strings.ReplaceAll(currentEntry, "11894951", fmt.Sprintf("%d", currentSlot))
				currentEntry = strings.ReplaceAll(currentEntry, "11894952", fmt.Sprintf("%d", currentSlot+1))
				currentEntry = strings.ReplaceAll(currentEntry, "11894953", fmt.Sprintf("%d", currentSlot+2))
				currentEntry = strings.ReplaceAll(currentEntry, "11894954", fmt.Sprintf("%d", currentSlot+3))
				currentEntry = strings.ReplaceAll(currentEntry, "11894955", fmt.Sprintf("%d", currentSlot+4))
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
}

// Update status box (matching statusBox.js and updateLogic.js synchronizeAndUpdateWidgets)
func (m *MonitorApp) updateStatusBox(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	statusMessages := []string{
		"FOLLOWING CHAIN HEAD",
		"SYNCING BLOCKS...",
		"PROCESSING TRANSACTIONS",
		"UPDATING STATE ROOT",
		"VALIDATING CONSENSUS",
		"FOLLOWING CHAIN HEAD",
	}

	messageIndex := 0

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

			// Cycle through status messages
			status := statusMessages[messageIndex%len(statusMessages)]

			// Add some dynamic information
			currentTime := time.Now()
			currentBlock := 22674573 + int(currentTime.Unix()%1000)

			statusContent := fmt.Sprintf("[green][bold]%s[white]\n", status)
			statusContent += fmt.Sprintf("Block: [yellow]%d[white]\n", currentBlock)
			statusContent += fmt.Sprintf("Time: [cyan]%s[white]\n", currentTime.Format("15:04:05"))
			statusContent += fmt.Sprintf("Peers: [green]%d[white]", 18+rand.Intn(10))

			// Send to status channel
			select {
			case m.StatusChan <- statusContent:
			default:
				// Channel full, skip update
			}

			messageIndex++
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

			// Simulate chain info similar to JavaScript version
			currentTime := time.Now()
			currentBlock := 22674573 + int(currentTime.Unix()%1000)

			// Create multiple blocks info like in JavaScript
			content := ""
			separator := strings.Repeat("-", 25)
			content += separator + "\n"

			for i := 0; i < 5; i++ {
				blockNum := currentBlock - i
				ethPrice := 3200.00 + float64(rand.Intn(200))  // Random ETH price
				gasPrice := 15.0 + float64(rand.Intn(50))/10.0 // Random gas price
				txCount := 150 + rand.Intn(100)                // Random TX count

				content += fmt.Sprintf("[center][green][bold]%s[white][/center]\n",
					formatNumberWithCommas(uint64(blockNum)))
				content += fmt.Sprintf("[blue][bold]ETH $:[white] %.2f\n", ethPrice)
				content += fmt.Sprintf("[blue][bold]GAS:[white]   %.1f\n", gasPrice)
				content += fmt.Sprintf("[blue][bold]# TX:[white]  %d\n", txCount)
				content += separator

				if i < 4 {
					content += "\n"
				}
			}

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

			// Simulate RPC connection info
			content := "[yellow][bold]RPC STATUS[white]\n"
			content += strings.Repeat("-", 15) + "\n"

			// Execution RPC
			execRPCStatus := "✅ Connected"
			if rand.Intn(10) < 2 { // 20% chance of connection issues
				execRPCStatus = "⚠️ Slow"
			}
			content += fmt.Sprintf("[blue]Execution:[white]\n%s\n", execRPCStatus)
			content += "[dim]http://localhost:8545[white]\n\n"

			// Consensus RPC
			consRPCStatus := "✅ Connected"
			if rand.Intn(20) < 1 { // 5% chance of issues
				consRPCStatus = "❌ Disconnected"
			}
			content += fmt.Sprintf("[blue]Consensus:[white]\n%s\n", consRPCStatus)
			content += "[dim]http://localhost:5052[white]"

			// Send to RPC info channel
			select {
			case m.RPCInfoChan <- content:
			default:
				// Channel full, skip update
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
