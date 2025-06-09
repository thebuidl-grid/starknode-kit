package pkg

import (
	"fmt"
	"os"
	"path"
	t "starknode-kit/pkg/types"

	"gopkg.in/yaml.v3"
)

var (
	InstallDir = path.Join(getHomeDir(), "starknode-kit")

	InstallClientsDir = path.Join(InstallDir, "ethereum_clients")

	jwtDir         = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath        = path.Join(jwtDir, "jwt.hex")
	configDir      = path.Join(InstallDir, "config")
	yamlConfigPath = fmt.Sprintf("%s/starknode.yaml", configDir)
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

func getHomeDir() string {
	homeDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}

func IsInstalled(c t.ClientType) error {
	dir := path.Join(InstallClientsDir, string(c))
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return err
	}
	if !info.IsDir() {
		return err
	}
	return nil
}

func LoadConfig() (t.StarkNodeKitConfig, error) {
	var cfg t.StarkNodeKitConfig
	cfgByt, err := os.ReadFile(yamlConfigPath)
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
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to update config file: %w", err)
	}
	cfg, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(yamlConfigPath, cfg, 0600)
	if err != nil {
		return err
	}
	return nil
}

func CreateStarkNodeConfig() error {
	default_config := defaultConfig()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	cfg, err := yaml.Marshal(default_config)
	if err != nil {
		return err
	}
	err = os.WriteFile(yamlConfigPath, cfg, 0600)
	if err != nil {
		return err
	}
	return nil
}

func defaultConfig() t.StarkNodeKitConfig {
	return t.StarkNodeKitConfig{
		Network: "sepolia",
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
	}
}

func GetClientVersion(clientName string) string {
	// This would typically be cached or retrieved from a version check
	// For now, return a placeholder
	switch clientName {
	case "geth":
		return latestGethVersion
	case "reth":
		return latestRethVersion
	case "lighthouse":
		return latestLighthouseVersion
	case "prysm":
		return latestPrysmVersion
	default:
		return "unknown"
	}
}
