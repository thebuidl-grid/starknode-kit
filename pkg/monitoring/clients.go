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
	"time"

	"starknode-kit/pkg"
	"starknode-kit/pkg/types"
)

// getGethSyncStatus gets sync status from Geth's HTTP API
func GetGethSyncStatus() types.SyncInfo {
	syncInfo := types.SyncInfo{IsSyncing: false, SyncPercent: 100.0}

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
func GetRethSyncStatus() types.SyncInfo {
	// Similar to Geth but Reth might have different endpoints
	return getGethSyncStatus() // For now, use same logic
}

// getLighthouseSyncStatus gets sync status from Lighthouse's HTTP API
func GetLighthouseSyncStatus() types.SyncInfo {
	syncInfo := types.SyncInfo{IsSyncing: false, SyncPercent: 100.0}

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
func GetPrysmSyncStatus() types.SyncInfo {
	syncInfo := types.SyncInfo{IsSyncing: false, SyncPercent: 100.0}

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
func GetEthereumMetrics() types.EthereumMetrics {
	metrics := types.EthereumMetrics{
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
