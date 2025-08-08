package monitoring

import (
	// "context"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/thebuidl-grid/starknode-kit/pkg"
	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	// "regexp"

	"strings"
	"time"
)

// utils.ParseHexInt parses a hex string to uint64

// GetEthereumMetrics gets blockchain metrics
func GetEthereumMetrics() t.EthereumMetrics {
	config, _ := utils.LoadConfig()
	metrics := t.EthereumMetrics{
		NetworkName: config.Network,
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
		var result map[string]any
		if json.Unmarshal(body, &result) == nil {
			if blockHex, ok := result["result"].(string); ok {
				if block, err := utils.ParseHexInt(blockHex); err == nil {
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
		var gasResult map[string]any
		if json.Unmarshal(gasBody, &gasResult) == nil {
			if gasPriceHex, ok := gasResult["result"].(string); ok {
				if gasPrice, err := utils.ParseHexInt(gasPriceHex); err == nil {
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
		var peerResult map[string]any
		if json.Unmarshal(peerBody, &peerResult) == nil {
			if peerHex, ok := peerResult["result"].(string); ok {
				if peers, err := utils.ParseHexInt(peerHex); err == nil {
					metrics.PeerCount = int(peers)
				}
			}
		}
	}

	return metrics
}

func GetJunoMetrics() t.EthereumMetrics {
	config, _ := utils.LoadConfig()
	metrics := t.EthereumMetrics{
		NetworkName: config.Network,
		IsSyncing:   false,
		SyncPercent: 100.0,
	}

	client := &http.Client{Timeout: 2 * time.Second}

	// Get current block number
	blockPayload := `{"jsonrpc":"2.0","method":"starknet_blockNumber","params":[],"id":1}`
	resp, err := client.Post("http://localhost:6060", "application/json", strings.NewReader(blockPayload))
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var result map[string]any
		if json.Unmarshal(body, &result) == nil {
			if block, ok := result["result"].(float64); ok {
				metrics.CurrentBlock = uint64(block)
			}
		}
	}

	// Get gas price
	gasPricePayload := `{"jsonrpc":"2.0","method":"starknet_syncing","params":[],"id":3}`
	gasResp, err := client.Post("http://localhost:6060", "application/json", strings.NewReader(gasPricePayload))
	if err == nil {
		defer gasResp.Body.Close()
		syncBody, _ := io.ReadAll(gasResp.Body)
		var syncResult map[string]any
		if json.Unmarshal(syncBody, &syncResult) == nil {
			if result, ok := syncResult["result"].(map[string]any); ok {
				currentBlock, ok := result["current_block_num"].(float64)
				if !ok {
					return t.EthereumMetrics{}
				}
				hightestBlock, ok := result["highest_block_num"].(float64)
				if !ok {
					return t.EthereumMetrics{}
				}
				metrics.IsSyncing = hightestBlock > currentBlock
				metrics.CurrentBlock = uint64(currentBlock)
				metrics.SyncPercent = (currentBlock / hightestBlock) * 100
			}
		}
	}
	return metrics
}

// GetLatestLogs gets the latest log entries from client log files
func GetLatestLogs(clientName string, lines int) []string {
	logDir := filepath.Join(pkg.InstallClientsDir, clientName, "logs")

	// NOTE minor fix
	if clientName == "juno" {
		logDir = filepath.Join(pkg.InstallStarknetDir, clientName, "logs")
	}

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
