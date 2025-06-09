package monitoring

import (
	// "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	// "regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"starknode-kit/pkg"
	"starknode-kit/pkg/clients"
)

// ProcessInfo holds information about a running process
type ProcessInfo struct {
	PID      int           `json:"pid"`
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Uptime   time.Duration `json:"uptime"`
	CPUUsage float64       `json:"cpu_usage"`
	MemUsage uint64        `json:"mem_usage"`
}

// EthereumMetrics holds blockchain metrics
type EthereumMetrics struct {
	CurrentBlock uint64  `json:"current_block"`
	HighestBlock uint64  `json:"highest_block"`
	SyncPercent  float64 `json:"sync_percent"`
	PeerCount    int     `json:"peer_count"`
	IsSyncing    bool    `json:"is_syncing"`
	GasPrice     string  `json:"gas_price"`
	NetworkName  string  `json:"network"`
}

// GetRunningClients returns information about running Ethereum clients
func GetRunningClients() []ClientStatus {
	var clients []ClientStatus

	// Check for Geth
	if gethInfo := getProcessInfo("geth"); gethInfo != nil {
		status := ClientStatus{
			Name:       "Geth",
			Status:     gethInfo.Status,
			PID:        gethInfo.PID,
			Uptime:     gethInfo.Uptime,
			Version:    getClientVersion("geth"),
			SyncStatus: getGethSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Reth
	if rethInfo := getProcessInfo("reth"); rethInfo != nil {
		status := ClientStatus{
			Name:       "Reth",
			Status:     rethInfo.Status,
			PID:        rethInfo.PID,
			Uptime:     rethInfo.Uptime,
			Version:    getClientVersion("reth"),
			SyncStatus: getRethSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Lighthouse
	if lighthouseInfo := getProcessInfo("lighthouse"); lighthouseInfo != nil {
		status := ClientStatus{
			Name:       "Lighthouse",
			Status:     lighthouseInfo.Status,
			PID:        lighthouseInfo.PID,
			Uptime:     lighthouseInfo.Uptime,
			Version:    getClientVersion("lighthouse"),
			SyncStatus: getLighthouseSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Prysm
	if prysmInfo := getProcessInfo("prysm"); prysmInfo != nil {
		status := ClientStatus{
			Name:       "Prysm",
			Status:     prysmInfo.Status,
			PID:        prysmInfo.PID,
			Uptime:     prysmInfo.Uptime,
			Version:    getClientVersion("prysm"),
			SyncStatus: getPrysmSyncStatus(),
		}
		clients = append(clients, status)
	}

	return clients
}

// getProcessInfo gets information about a running process by name
func getProcessInfo(processName string) *ProcessInfo {
	// Read /proc to find the process
	procDirs, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil
	}

	for _, procDir := range procDirs {
		// Read cmdline to get the command
		cmdlineFile := filepath.Join(procDir, "cmdline")
		cmdlineBytes, err := os.ReadFile(cmdlineFile)
		if err != nil {
			continue
		}

		cmdline := string(cmdlineBytes)
		if strings.Contains(cmdline, processName) {
			// Extract PID from directory name
			pidStr := filepath.Base(procDir)
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			// Get process status
			statusFile := filepath.Join(procDir, "stat")
			statusBytes, err := os.ReadFile(statusFile)
			if err != nil {
				continue
			}

			// Parse stat file for uptime
			statFields := strings.Fields(string(statusBytes))
			if len(statFields) > 21 {
				startTimeJiffies, err := strconv.ParseUint(statFields[21], 10, 64)
				if err == nil {
					// Calculate uptime (simplified)
					uptimeSeconds := time.Now().Unix() - int64(startTimeJiffies/100)
					uptime := time.Duration(uptimeSeconds) * time.Second

					return &ProcessInfo{
						PID:    pid,
						Name:   processName,
						Status: "running",
						Uptime: uptime,
					}
				}
			}
		}
	}

	return nil
}

// getClientVersion gets the version of an installed client
func getClientVersion(clientName string) string {
	// This would typically be cached or retrieved from a version check
	// For now, return a placeholder
	switch clientName {
	case "geth":
		return "v1.13.5"
	case "reth":
		return "v0.2.0-beta"
	case "lighthouse":
		return "v4.5.0"
	case "prysm":
		return "v4.2.1"
	default:
		return "unknown"
	}
}

// getGethSyncStatus gets sync status from Geth's HTTP API
func getGethSyncStatus() SyncInfo {
	syncInfo := SyncInfo{IsSyncing: false, SyncPercent: 100.0}

	// Try to get sync status from Geth's HTTP API
	client := &http.Client{Timeout: 2 * time.Second}

	// Geth eth_syncing call
	payload := `{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}`
	resp, err := client.Post("http://localhost:8545", "application/json", strings.NewReader(payload))
	if err != nil {
		return syncInfo
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return syncInfo
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return syncInfo
	}

	if syncResult, ok := result["result"]; ok {
		if syncResult == false {
			// Not syncing = fully synced
			syncInfo.IsSyncing = false
			syncInfo.SyncPercent = 100.0
		} else if syncData, ok := syncResult.(map[string]interface{}); ok {
			// Currently syncing
			syncInfo.IsSyncing = true

			if currentHex, ok := syncData["currentBlock"].(string); ok {
				if current, err := parseHexInt(currentHex); err == nil {
					syncInfo.CurrentBlock = current
				}
			}

			if highestHex, ok := syncData["highestBlock"].(string); ok {
				if highest, err := parseHexInt(highestHex); err == nil {
					syncInfo.HighestBlock = highest
					if highest > 0 {
						syncInfo.SyncPercent = float64(syncInfo.CurrentBlock) / float64(highest) * 100
					}
				}
			}
		}
	}

	// Get peer count
	peerPayload := `{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":2}`
	peerResp, err := client.Post("http://localhost:8545", "application/json", strings.NewReader(peerPayload))
	if err == nil {
		defer peerResp.Body.Close()
		peerBody, _ := io.ReadAll(peerResp.Body)
		var peerResult map[string]interface{}
		if json.Unmarshal(peerBody, &peerResult) == nil {
			if peerHex, ok := peerResult["result"].(string); ok {
				if peers, err := parseHexInt(peerHex); err == nil {
					syncInfo.PeersCount = int(peers)
				}
			}
		}
	}

	return syncInfo
}

// getRethSyncStatus gets sync status from Reth's HTTP API
func getRethSyncStatus() SyncInfo {
	// Similar to Geth but Reth might have different endpoints
	return getGethSyncStatus() // For now, use same logic
}

// getLighthouseSyncStatus gets sync status from Lighthouse's HTTP API
func getLighthouseSyncStatus() SyncInfo {
	syncInfo := SyncInfo{IsSyncing: false, SyncPercent: 100.0}

	client := &http.Client{Timeout: 2 * time.Second}

	// Lighthouse HTTP API
	resp, err := client.Get("http://localhost:5052/eth/v1/node/syncing")
	if err != nil {
		return syncInfo
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return syncInfo
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return syncInfo
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		if isSyncing, ok := data["is_syncing"].(bool); ok {
			syncInfo.IsSyncing = isSyncing
		}

		if headSlot, ok := data["head_slot"].(string); ok {
			if head, err := strconv.ParseUint(headSlot, 10, 64); err == nil {
				syncInfo.CurrentBlock = head
			}
		}

		if syncDistance, ok := data["sync_distance"].(string); ok {
			if distance, err := strconv.ParseUint(syncDistance, 10, 64); err == nil {
				syncInfo.HighestBlock = syncInfo.CurrentBlock + distance
				if syncInfo.HighestBlock > 0 {
					syncInfo.SyncPercent = float64(syncInfo.CurrentBlock) / float64(syncInfo.HighestBlock) * 100
				}
			}
		}
	}

	return syncInfo
}

// getPrysmSyncStatus gets sync status from Prysm's HTTP API
func getPrysmSyncStatus() SyncInfo {
	syncInfo := SyncInfo{IsSyncing: false, SyncPercent: 100.0}

	client := &http.Client{Timeout: 2 * time.Second}

	// Prysm HTTP API
	resp, err := client.Get("http://localhost:5052/eth/v1/node/syncing")
	if err != nil {
		return syncInfo
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return syncInfo
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return syncInfo
	}

	// Similar to Lighthouse
	if data, ok := result["data"].(map[string]interface{}); ok {
		if isSyncing, ok := data["is_syncing"].(bool); ok {
			syncInfo.IsSyncing = isSyncing
		}

		if headSlot, ok := data["head_slot"].(string); ok {
			if head, err := strconv.ParseUint(headSlot, 10, 64); err == nil {
				syncInfo.CurrentBlock = head
			}
		}

		if syncDistance, ok := data["sync_distance"].(string); ok {
			if distance, err := strconv.ParseUint(syncDistance, 10, 64); err == nil {
				syncInfo.HighestBlock = syncInfo.CurrentBlock + distance
				if syncInfo.HighestBlock > 0 {
					syncInfo.SyncPercent = float64(syncInfo.CurrentBlock) / float64(syncInfo.HighestBlock) * 100
				}
			}
		}
	}

	return syncInfo
}

// GetEthereumMetrics gets blockchain metrics
func GetEthereumMetrics() EthereumMetrics {
	metrics := EthereumMetrics{
		NetworkName: "Mainnet",
		IsSyncing:   false,
		SyncPercent: 100.0,
	}

	client := &http.Client{Timeout: 2 * time.Second}

	// Get current block number
	blockPayload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`
	resp, err := client.Post("http://localhost:8545", "application/json", strings.NewReader(blockPayload))
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		if json.Unmarshal(body, &result) == nil {
			if blockHex, ok := result["result"].(string); ok {
				if block, err := parseHexInt(blockHex); err == nil {
					metrics.CurrentBlock = block
				}
			}
		}
	}

	// Get gas price
	gasPricePayload := `{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":3}`
	gasResp, err := client.Post("http://localhost:8545", "application/json", strings.NewReader(gasPricePayload))
	if err == nil {
		defer gasResp.Body.Close()
		gasBody, _ := io.ReadAll(gasResp.Body)
		var gasResult map[string]interface{}
		if json.Unmarshal(gasBody, &gasResult) == nil {
			if gasPriceHex, ok := gasResult["result"].(string); ok {
				if gasPrice, err := parseHexInt(gasPriceHex); err == nil {
					// Convert wei to gwei
					gweiPrice := float64(gasPrice) / 1e9
					metrics.GasPrice = fmt.Sprintf("%.1f gwei", gweiPrice)
				}
			}
		}
	}

	// Get peer count
	peerPayload := `{"jsonrpc":"2.0","method":"net_peerCount","params":[],"id":2}`
	peerResp, err := client.Post("http://localhost:8545", "application/json", strings.NewReader(peerPayload))
	if err == nil {
		defer peerResp.Body.Close()
		peerBody, _ := io.ReadAll(peerResp.Body)
		var peerResult map[string]interface{}
		if json.Unmarshal(peerBody, &peerResult) == nil {
			if peerHex, ok := peerResult["result"].(string); ok {
				if peers, err := parseHexInt(peerHex); err == nil {
					metrics.PeerCount = int(peers)
				}
			}
		}
	}

	return metrics
}

// GetLatestLogs gets the latest log entries from client log files
func GetLatestLogs(clientName string, lines int) []string {
	logDir := filepath.Join(pkg.InstallClientsDir, clientName, "logs")

	// Find the most recent log file
	files, err := filepath.Glob(filepath.Join(logDir, "*.log"))
	if err != nil || len(files) == 0 {
		return []string{fmt.Sprintf("No log files found for %s", clientName)}
	}

	// Get the most recent file
	var newestFile string
	var newestTime time.Time
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		if info.ModTime().After(newestTime) {
			newestTime = info.ModTime()
			newestFile = file
		}
	}

	if newestFile == "" {
		return []string{"No recent log files found"}
	}

	// Read the last few lines
	content, err := os.ReadFile(newestFile)
	if err != nil {
		return []string{fmt.Sprintf("Error reading log file: %v", err)}
	}

	logLines := strings.Split(string(content), "\n")
	startIdx := len(logLines) - lines - 1
	if startIdx < 0 {
		startIdx = 0
	}

	var result []string
	for i := startIdx; i < len(logLines) && len(result) < lines; i++ {
		if strings.TrimSpace(logLines[i]) != "" {
			result = append(result, logLines[i])
		}
	}

	return result
}

// parseHexInt parses a hex string to uint64
func parseHexInt(hexStr string) (uint64, error) {
	// Remove 0x prefix if present
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	return strconv.ParseUint(hexStr, 16, 64)
}

// IsProcessRunning checks if a process is running by PID
func IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// StopClient stops a client by name
func StopClient(clientName string) error {
	processInfo := getProcessInfo(clientName)
	if processInfo == nil {
		return fmt.Errorf("client %s is not running", clientName)
	}

	return pkg.StopProcess(processInfo.PID)
}

// RestartClient restarts a client
func RestartClient(clientName string) error {
	// First stop the client
	if err := StopClient(clientName); err != nil {
		// If it wasn't running, that's ok
		if !strings.Contains(err.Error(), "not running") {
			return err
		}
	}

	// Wait a bit for clean shutdown
	time.Sleep(2 * time.Second)

	// Load config to get client settings
	config, err := pkg.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Start based on client type
	switch clientName {
	case "geth":
		return clients.StartGeth(config.ExecutionCientSettings.ExecutionType, config.ExecutionCientSettings.Port)
	case "reth":
		return clients.StartReth(config.ExecutionCientSettings.ExecutionType, config.ExecutionCientSettings.Port)
	case "lighthouse":
		return clients.StartLightHouse(config.ConsensusCientSettings.Port...)
	case "prysm":
		return clients.StartPrsym(config.ConsensusCientSettings.Port...)
	default:
		return fmt.Errorf("unknown client: %s", clientName)
	}
}
