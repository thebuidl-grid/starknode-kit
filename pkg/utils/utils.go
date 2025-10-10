package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/versions"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	starkutils "github.com/NethermindEth/starknet.go/utils"
	envsubt "github.com/emperorsixpacks/envsubst"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func GetExecutionClient(c string) (t.ClientType, error) {
	sprtClients := map[string]t.ClientType{
		"geth": t.ClientGeth,
		"reth": t.ClientReth,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("execution client %s not supported", c)
	}
	return client, nil
}
func GetConsensusClient(c string) (t.ClientType, error) {
	sprtClients := map[string]t.ClientType{
		"lighthouse": t.ClientLighthouse,
		"prysm":      t.ClientPrysm,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("consensus client %s not supported", c)
	}
	return client, nil
}

func IsInstalled(c t.ClientType) bool {
	client := strings.ToLower(string(c))
	dir := path.Join(constants.InstallClientsDir, client)
	if c == t.ClientStarkValidator || c == t.ClientJuno {
		dir = path.Join(constants.InstallStarknetDir, client)
	}
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	if !info.IsDir() {
		return false
	}
	return true
}

func LoadConfig() (t.StarkNodeKitConfig, error) {
	var cfg t.StarkNodeKitConfig
	cfgByt, err := os.ReadFile(constants.ConfigPath)
	if err != nil {
		return t.StarkNodeKitConfig{}, err
	}
	err = godotenv.Load(constants.EnvFIlePath)
	if err == nil {
		err = envsubt.Unmarshal(cfgByt, &cfg)
		if err != nil {
			return t.StarkNodeKitConfig{}, err
		}
		return cfg, nil
	}
	err = yaml.Unmarshal(cfgByt, &cfg)
	if err != nil {
		return t.StarkNodeKitConfig{}, err
	}
	return cfg, nil
}

func UpdateStarkNodeConfig(config t.StarkNodeKitConfig) error {
	if err := os.MkdirAll(constants.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}
	cfg, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(constants.ConfigPath, cfg, 0600)
	if err != nil {
		return err
	}
	return nil
}

func CreateStarkNodeConfig(cfg *types.StarkNodeKitConfig) error {
	var setupConfig *types.StarkNodeKitConfig
	if _, err := os.Stat(constants.ConfigPath); err == nil {
		fmt.Println(Yellow(fmt.Sprintf("Starknode-kit already initialized at %s", constants.ConfigDir)))
	}
	cfg.ConsensusCientSettings.ConsensusCheckpoint = fmt.Sprintf("https://beaconstate-%s.chainsafe.io", cfg.Network)

	if cfg == nil {
		setupConfig = defaultConfig()
	} else {
		setupConfig = cfg
	}
	if err := os.MkdirAll(constants.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	if setupConfig.Wallet.Name != "" {
		setupConfig.Wallet.Wallet.Normalize()
	}
	conigBytes, err := yaml.Marshal(*setupConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(constants.ConfigPath, conigBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func GetRunningClients() []types.ClientStatus {
	var clients []types.ClientStatus

	// Check for Geth
	if gethInfo := process.GetProcessInfo("geth"); gethInfo != nil {
		status := types.ClientStatus{
			Name:       "Geth",
			Status:     gethInfo.Status,
			PID:        gethInfo.PID,
			Uptime:     gethInfo.Uptime,
			Version:    versions.GetVersionNumber("geth"),
			SyncStatus: GetGethSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Reth
	if rethInfo := process.GetProcessInfo("reth"); rethInfo != nil {
		status := types.ClientStatus{
			Name:       "Reth",
			Status:     rethInfo.Status,
			PID:        rethInfo.PID,
			Uptime:     rethInfo.Uptime,
			Version:    versions.GetVersionNumber("reth"),
			SyncStatus: GetRethSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Lighthouse
	if lighthouseInfo := process.GetProcessInfo("lighthouse"); lighthouseInfo != nil {
		status := types.ClientStatus{
			Name:       "Lighthouse",
			Status:     lighthouseInfo.Status,
			PID:        lighthouseInfo.PID,
			Uptime:     lighthouseInfo.Uptime,
			Version:    versions.GetVersionNumber("lighthouse"),
			SyncStatus: GetLighthouseSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Prysm
	if prysmInfo := process.GetProcessInfo("prysm"); prysmInfo != nil {
		status := types.ClientStatus{
			Name:       "Prysm",
			Status:     prysmInfo.Status,
			PID:        prysmInfo.PID,
			Uptime:     prysmInfo.Uptime,
			Version:    versions.GetVersionNumber("prysm"),
			SyncStatus: GetPrysmSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Juno (Starknet client)
	if junoInfo := process.GetProcessInfo("juno"); junoInfo != nil {
		status := types.ClientStatus{
			Name:    "Juno",
			Status:  junoInfo.Status,
			PID:     junoInfo.PID,
			Uptime:  junoInfo.Uptime,
			Version: versions.GetVersionNumber("juno"),
		}
		clients = append(clients, status)
	}

	// Check for Starknet Validator
	if validatorInfo := process.GetProcessInfo("starknet-staking-v2"); validatorInfo != nil {
		status := types.ClientStatus{
			Name:    "Validator",
			Status:  validatorInfo.Status,
			PID:     validatorInfo.PID,
			Uptime:  validatorInfo.Uptime,
			Version: versions.GetVersionNumber("starknet-staking-v2"),
		}
		clients = append(clients, status)
	}

	return clients
}

func ParseHexInt(hexStr string) (uint64, error) {
	// Remove 0x prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")
	return strconv.ParseUint(hexStr, 16, 64)
}

func SetNetwork(cfg *t.StarkNodeKitConfig, network string) error {
	switch network {
	case "mainnet":
		cfg.Network = "mainnet"
		cfg.ConsensusCientSettings.ConsensusCheckpoint = "https://mainnet-checkpoint-sync.stakely.io/"
		return nil
	case "sepolia":
		cfg.Network = "sepolia"
		cfg.ConsensusCientSettings.ConsensusCheckpoint = "https://sepolia-checkpoint-sync.stakely.io/"
		return nil
	default:
		return fmt.Errorf("Network %v not supported", network)
	}
}

func GetStarknetClient(c string) (t.ClientType, error) {
	sprtClients := map[string]t.ClientType{
		"juno":                t.ClientJuno,
		"starknet-staking-v2": t.ClientStarkValidator,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("starknet client %s not supported", c)
	}
	return client, nil
}

func ViewConfig() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	fmt.Println(cfg)
	return nil

}

// PadZerosInFelt pads a felt value to 66 characters with leading zeros
// This ensures consistent formatting for Starknet addresses and hashes
func PadZerosInFelt(hexFelt *felt.Felt) string {
	const targetLength = 66
	hexStr := hexFelt.String()

	// Check if the hex value is already of the desired length
	if len(hexStr) >= targetLength {
		return hexStr
	}

	// Extract the hex value without the "0x" prefix
	hexValue := hexStr[2:]

	// Pad zeros after the "0x" prefix
	paddedHexValue := fmt.Sprintf("%0*s", targetLength-2, hexValue)

	// Add back the "0x" prefix to the padded hex value
	return "0x" + paddedHexValue
}

// FormatStarknetAddress formats a felt address with proper padding
func FormatStarknetAddress(addr *felt.Felt) string {
	return PadZerosInFelt(addr)
}

// FormatTransactionHash formats a transaction hash with proper padding
func FormatTransactionHash(hash *felt.Felt) string {
	return PadZerosInFelt(hash)
}

func CheckRPCStatus(rpcURL, method string) (string, error) {
	payload := map[string]any{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  []any{},
		"id":      1,
	}
	body, _ := json.Marshal(payload)

	start := time.Now()
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(body))
	latency := time.Since(start)

	if err != nil {
		return "❌ Disconnected", err
	}
	defer resp.Body.Close()

	// Determine status based on latency
	if latency > 1*time.Second {
		return "⚠️ Slow", nil
	}
	return "✅ Connected", nil
}

func EstimateGasFee(accnt *account.Account, callData []rpc.FunctionCall) (*rpc.BroadcastInvokeTxnV3, []rpc.FeeEstimation, error) {
	nonce, err := accnt.Nonce(context.Background())
	if err != nil {
		return nil, nil, err
	}
	calldata, err := accnt.FmtCalldata(callData)
	if err != nil {
		return nil, nil, err
	}

	InvokeTx := starkutils.BuildInvokeTxn(
		accnt.Address,
		nonce,
		calldata,
		&rpc.ResourceBoundsMapping{
			L1Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L1DataGas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
			L2Gas: rpc.ResourceBounds{
				MaxAmount:       "0x0",
				MaxPricePerUnit: "0x0",
			},
		},
		nil,
	)

	err = accnt.SignInvokeTransaction(context.Background(), InvokeTx)
	if err != nil {
		return nil, nil, err
	}

	feeRes, err := accnt.Provider.EstimateFee(
		context.Background(),
		[]rpc.BroadcastTxn{InvokeTx},
		[]rpc.SimulationFlag{},
		rpc.WithBlockTag(rpc.BlockTagLatest),
	)
	if err != nil {
		return nil, nil, err
	}
	return InvokeTx, feeRes, nil
}
