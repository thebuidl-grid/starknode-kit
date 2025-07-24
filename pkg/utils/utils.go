package utils

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/thebuidl-grid/starknode-kit/pkg"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/versions"

	"github.com/NethermindEth/juno/core/felt"
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
	dir := path.Join(pkg.InstallClientsDir, string(c))
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
	cfgByt, err := os.ReadFile(pkg.ConfigPath)
	if err != nil {
		return t.StarkNodeKitConfig{}, err
	}
	err = godotenv.Load(pkg.EnvFIlePath)
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
	if err := os.MkdirAll(pkg.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}
	cfg, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(pkg.ConfigPath, cfg, 0600)
	if err != nil {
		return err
	}
	return nil
}

func CreateStarkNodeConfig() error {
	if _, err := os.Stat(pkg.ConfigDir); err == nil {
		return fmt.Errorf("Starknode-kit already initialized at %s", pkg.ConfigDir)
	}

	default_config := defaultConfig()
	if err := os.MkdirAll(pkg.ConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	cfg, err := yaml.Marshal(default_config)
	if err != nil {
		return err
	}
	err = os.WriteFile(pkg.ConfigPath, cfg, 0600)
	if err != nil {
		return err
	}
	return nil
}

func GetClientVersion(clientName string) string {
	// This would typically be cached or retrieved from a version check
	// For now, return a placeholder
	switch clientName {
	case "geth":
		return versions.LatestGethVersion
	case "reth":
		return versions.LatestRethVersion
	case "lighthouse":
		return versions.LatestLighthouseVersion
	case "prysm":
		return "" //versions.LatestPrysmVersion
	default:
		return "unknown"
	}
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
			Version:    GetClientVersion("geth"),
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
			Version:    GetClientVersion("reth"),
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
			Version:    GetClientVersion("lighthouse"),
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
			Version:    GetClientVersion("prysm"),
			SyncStatus: GetPrysmSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Juno (Starknet client)
	if junoInfo := process.GetProcessInfo("juno"); junoInfo != nil {
		status := types.ClientStatus{
			Name:       "Juno",
			Status:     junoInfo.Status,
			PID:        junoInfo.PID,
			Uptime:     junoInfo.Uptime,
			Version:    GetClientVersion("juno"),
			SyncStatus: GetJunoSyncStatus(),
		}
		clients = append(clients, status)
	}

	return clients
}

func ParseHexInt(hexStr string) (uint64, error) {
	// Remove 0x prefix if present
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
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
		"juno": t.ClientJuno,
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
