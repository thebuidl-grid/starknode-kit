package pkg

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	InstallDir = path.Join(getHomeDir(), "starcknode-kit")

	InstallClientsDir = path.Join(InstallDir, "ethereum_clients")

	jwtDir         = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath        = path.Join(jwtDir, "jwt.hex")
	configDir      = path.Join(InstallDir, "config")
	yamlConfigPath = fmt.Sprintf("%s/stacknode.yaml", configDir)
)

func GetExecutionClient(c string) (ClientType, error) {
	sprtClients := map[string]ClientType{
		"geth": ClientGeth,
		"reth": ClientReth,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("Execution Client %s not supported", client)
	}
	return client, nil
}
func GetConsensusClient(c string) (ClientType, error) {
	sprtClients := map[string]ClientType{
		"lighthouse": ClientLighthouse,
		"prysm":      ClientPrysm,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("Consensus Client %s not supported", client) // TODO change this
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

func LoadConfig() (StarkNodeKitConfig, error) {
	var cfg StarkNodeKitConfig
	cfgByt, err := os.ReadFile(yamlConfigPath)
	if err != nil {
		return StarkNodeKitConfig{}, err
	}
	err = yaml.Unmarshal(cfgByt, &cfg)
	if err != nil {
		return StarkNodeKitConfig{}, err
	}
	return cfg, nil
}

func UpdateStackNodeConfig(config StarkNodeKitConfig) error {
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

func CreateStackNodeConfig() error {
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

func defaultConfig() StarkNodeKitConfig {
	return StarkNodeKitConfig{
		ExecutionCientSettings: ClientSettings{
			Name:    ClientGeth,
			Network: "sepolia",
			Port:    []string{"8545", "30303"},
		},
		ConsensusCientSettings: ClientSettings{
			Name:    ClientPrysm,
			Network: "sepolia",
			Port:    []string{"8545", "30303"},
		},
	}
}
