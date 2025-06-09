package monitoring

import (
	"context"
	"fmt"
	"strings"
	"time"

	"starknode-kit/pkg/stats"
)

func (m *MonitorApp) updateSystemStats(ctx context.Context) {
	ticker := time.NewTicker(m.UpdateRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case <-ticker.C:
			systemStats, err := stats.GetSystemStats()
			if err != nil {
				continue
			}

			statsText := fmt.Sprintf(
				"[green]CPU Usage:[white] %.1f%%\n"+
					"[green]CPU Cores:[white] %d\n"+
					"[green]Memory:[white] %s / %s (%.1f%%)\n"+
					"[green]Disk:[white] %s / %s (%.1f%%)\n"+
					"[green]Network:[white] ↑%s ↓%s\n"+
					"[green]Uptime:[white] %s",
				systemStats.CPU.Usage,
				systemStats.CPU.Cores,
				stats.FormatBytes(systemStats.Memory.Used),
				stats.FormatBytes(systemStats.Memory.Total),
				systemStats.Memory.UsedPercent,
				stats.FormatBytes(systemStats.Disk.Used),
				stats.FormatBytes(systemStats.Disk.Total),
				systemStats.Disk.UsedPercent,
				stats.FormatBytes(systemStats.Network.BytesSent),
				stats.FormatBytes(systemStats.Network.BytesRecv),
				stats.FormatUptime(systemStats.Uptime),
			)

			m.SystemChan <- statsText
		}
	}
}

func (m *MonitorApp) updateClientStatus(ctx context.Context) {
	ticker := time.NewTicker(m.UpdateRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case <-ticker.C:
			// Get real client status data
			clients := GetRunningClients()

			var clientText string
			if len(clients) == 0 {
				clientText = "[red]No clients running[white]\n"
			} else {
				for _, client := range clients {
					status := client.Status
					statusColor := "red"
					switch status {
					case "running":
						statusColor = "green"
					case "syncing":
						statusColor = "yellow"
					case "stopped":
						statusColor = "red"
					}

					syncText := ""
					if client.SyncStatus.IsSyncing {
						syncText = fmt.Sprintf(" (%.1f%%)", client.SyncStatus.SyncPercent)
					}

					clientText += fmt.Sprintf("[green]%s:[%s] %s[white]%s\n",
						client.Name, statusColor, status, syncText)

					if client.PID > 0 {
						clientText += fmt.Sprintf("  [blue]PID:[white] %d [blue]Uptime:[white] %s\n",
							client.PID, client.Uptime.Truncate(time.Second))
					}

					if client.SyncStatus.PeersCount > 0 {
						clientText += fmt.Sprintf("  [blue]Peers:[white] %d\n", client.SyncStatus.PeersCount)
					}

					clientText += "\n"
				}
			}

			m.ClientsChan <- clientText
		}
	}
}

func (m *MonitorApp) updateChainInfo(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case <-ticker.C:
			// Get real Ethereum metrics
			metrics := GetEthereumMetrics()

			chainText := fmt.Sprintf("[green]Network:[white] %s\n", metrics.NetworkName)

			if metrics.CurrentBlock > 0 {
				chainText += fmt.Sprintf("[green]Block:[white] %s\n",
					formatNumber(metrics.CurrentBlock))
			} else {
				chainText += "[red]Block:[white] Unavailable\n"
			}

			if metrics.IsSyncing {
				chainText += fmt.Sprintf("[yellow]Sync:[white] %.1f%%\n", metrics.SyncPercent)
			} else {
				chainText += "[green]Sync:[white] 100%\n"
			}

			if metrics.GasPrice != "" {
				chainText += fmt.Sprintf("[green]Gas Price:[white] %s\n", metrics.GasPrice)
			}

			if metrics.PeerCount > 0 {
				chainText += fmt.Sprintf("[green]Peers:[white] %d\n", metrics.PeerCount)
			}

			m.ChainChan <- chainText
		}
	}
}

// formatNumber formats a large number with commas
func formatNumber(n uint64) string {
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

func (m *MonitorApp) updatePeerInfo(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case <-ticker.C:
			// Get actual peer information from running clients
			clients := GetRunningClients()

			var executionPeers, consensusPeers int
			var executionClient, consensusClient string

			for _, client := range clients {
				switch client.Name {
				case "Geth", "Reth":
					executionPeers = client.SyncStatus.PeersCount
					executionClient = client.Name
				case "Lighthouse", "Prysm":
					consensusPeers = client.SyncStatus.PeersCount
					consensusClient = client.Name
				}
			}

			peerText := ""
			if executionClient != "" {
				peerText += fmt.Sprintf("[green]%s Peers:[white] %d\n", executionClient, executionPeers)
			} else {
				peerText += "[red]Execution Client:[white] Not running\n"
			}

			if consensusClient != "" {
				peerText += fmt.Sprintf("[green]%s Peers:[white] %d\n", consensusClient, consensusPeers)
			} else {
				peerText += "[red]Consensus Client:[white] Not running\n"
			}

			// Add network health indicator
			totalPeers := executionPeers + consensusPeers
			if totalPeers > 20 {
				peerText += "[green]Network Health:[white] Good\n"
			} else if totalPeers > 5 {
				peerText += "[yellow]Network Health:[white] Fair\n"
			} else {
				peerText += "[red]Network Health:[white] Poor\n"
			}

			m.PeersChan <- peerText
		}
	}
}

func (m *MonitorApp) updateLogInfo(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case <-ticker.C:
			// Get logs from running clients
			var allLogs []string

			clients := GetRunningClients()
			for _, client := range clients {
				logs := GetLatestLogs(strings.ToLower(client.Name), 5)
				for _, log := range logs {
					if strings.TrimSpace(log) != "" {
						timestamp := time.Now().Format("15:04:05")
						allLogs = append(allLogs, fmt.Sprintf("[dim]%s[white] [green]%s[white]: %s",
							timestamp, client.Name, log))
					}
				}
			}

			// Keep only the latest 50 log entries
			if len(allLogs) > 50 {
				allLogs = allLogs[len(allLogs)-50:]
			}

			logText := strings.Join(allLogs, "\n")
			if logText == "" {
				logText = "[yellow]No recent log entries...[white]"
			}

			m.LogsChan <- logText
		}
	}
}

func (m *MonitorApp) handleUpdates(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.StopChan:
			return
		case text := <-m.SystemChan:
			m.App.QueueUpdateDraw(func() {
				m.SystemBox.SetText(text)
			})
		case text := <-m.ClientsChan:
			m.App.QueueUpdateDraw(func() {
				m.ClientsBox.SetText(text)
			})
		case text := <-m.ChainChan:
			m.App.QueueUpdateDraw(func() {
				m.ChainBox.SetText(text)
			})
		case text := <-m.PeersChan:
			m.App.QueueUpdateDraw(func() {
				m.PeersBox.SetText(text)
			})
		case text := <-m.LogsChan:
			m.App.QueueUpdateDraw(func() {
				m.LogsBox.SetText(text)
			})
		}
	}
}
