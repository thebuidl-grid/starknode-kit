package utils

import (
	"fmt"
	"os"
	"path"
	"starknode-kit/pkg"
	"starknode-kit/pkg/process"
	"starknode-kit/pkg/types"
	t "starknode-kit/pkg/types"
	"starknode-kit/pkg/versions"
	"strconv"
	"strings"

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

	return clients
}

func ParseHexInt(hexStr string) (uint64, error) {
	// Remove 0x prefix if present
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	return strconv.ParseUint(hexStr, 16, 64)
}
