package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	t "github.com/thebuidl-grid/starknode-kit/pkg/types"

	"github.com/joho/godotenv"
)

func GetGethSyncStatus() t.SyncInfo {
	syncInfo := t.SyncInfo{IsSyncing: false, SyncPercent: 100.0}

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

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return syncInfo
	}
	if syncResult, ok := result["result"]; ok {
		if syncResult == false {
			// Not syncing = fully synced
			syncInfo.IsSyncing = false
			syncInfo.SyncPercent = 100.0
		} else if syncData, ok := syncResult.(map[string]any); ok {
			// Currently syncing
			syncInfo.IsSyncing = true

			if currentHex, ok := syncData["currentBlock"].(string); ok {
				if current, err := ParseHexInt(currentHex); err == nil {
					syncInfo.CurrentBlock = current
				}
			}

			if highestHex, ok := syncData["highestBlock"].(string); ok {
				if highest, err := ParseHexInt(highestHex); err == nil {
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
		var peerResult map[string]any
		if json.Unmarshal(peerBody, &peerResult) == nil {
			if peerHex, ok := peerResult["result"].(string); ok {
				if peers, err := ParseHexInt(peerHex); err == nil {
					syncInfo.PeersCount = int(peers)
				}
			}
		}
	}

	return syncInfo
}

// getRethSyncStatus gets sync status from Reth's HTTP API
func GetRethSyncStatus() t.SyncInfo {
	// Similar to Geth but Reth might have different endpoints
	return GetGethSyncStatus() // For now, use same logic
}

// getLighthouseSyncStatus gets sync status from Lighthouse's HTTP API
func GetLighthouseSyncStatus() t.SyncInfo {
	syncInfo := t.SyncInfo{IsSyncing: false, SyncPercent: 100.0}

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

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return syncInfo
	}

	if data, ok := result["data"].(map[string]any); ok {
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
func GetPrysmSyncStatus() t.SyncInfo {
	syncInfo := t.SyncInfo{IsSyncing: false, SyncPercent: 100.0}

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

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return syncInfo
	}

	// Similar to Lighthouse
	if data, ok := result["data"].(map[string]any); ok {
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

func GetJunoMetrics(network string) t.EthereumMetrics {
	metrics := t.EthereumMetrics{
		NetworkName: network,
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

func defaultConfig() *t.StarkNodeKitConfig {
	return &t.StarkNodeKitConfig{
		// WalletAddress: "${STARKNET_WALLET}",
		// PrivateKey:    "${STARKNET_PRIVATE_KEY}",
		Network: "mainnet",
		ExecutionCientSettings: t.ClientConfig{
			Name:          t.ClientGeth,
			Port:          []int{30303},
			ExecutionType: "full",
		},
		ConsensusCientSettings: t.ClientConfig{
			Name:                t.ClientPrysm,
			Port:                []int{5052, 9000},
			ConsensusCheckpoint: "https://mainnet-checkpoint-sync.stakely.io/",
		},
		JunoConfig: t.JunoConfig{
			Port:    6060,
			EthNode: "wss://eth.drpc.org",
			Environment: []string{
				"JUNO_HTTP_PORT=6060",
				"JUNO_HTTP_HOST=0.0.0.0",
			},
		},
	}
}

func writeToENV(ks map[string]string) error {
	var file *os.File

	_, err := os.Stat(constants.EnvFIlePath)
	if err := os.MkdirAll(constants.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	if os.IsNotExist(err) {
		_, err := os.Create(constants.EnvFIlePath)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	err = godotenv.Write(ks, constants.EnvFIlePath)
	return err
}

// TODO change juno get metrics url
