package pkg

import (
	"fmt"
	"os"
	"path"

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

func GetExecutionClient(c string) (ClientType, error) {
	sprtClients := map[string]ClientType{
		"geth": ClientGeth,
		"reth": ClientReth,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("execution client %s not supported", c)
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
		return "", fmt.Errorf("consensus client %s not supported", c)
	}
	return client, nil
}

func GetStarknetClient(c string) (ClientType, error) {
	sprtClients := map[string]ClientType{
		"juno": ClientJuno,
	}
	client, ok := sprtClients[c]
	if !ok {
		return "", fmt.Errorf("starknet client %s not supported", c)
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

func IsInstalled(c ClientType) error {
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

func UpdateStarkNodeConfig(config StarkNodeKitConfig) error {
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
	fmt.Printf("StarkNodeKitConfig created successfully")
	return nil
}

func defaultConfig() StarkNodeKitConfig {
	return StarkNodeKitConfig{
		ExecutionCientSettings: ClientConfig{
			Name:          ClientGeth,
			Network:       "sepolia",
			Port:          []int{30303},
			ExecutionType: "full",
		},
		ConsensusCientSettings: ClientConfig{
			Name:    ClientPrysm,
			Network: "sepolia",
			Port:    []int{5052, 9000},
		},
	}
}

func ViewConfig() error {
	cfg, err := LoadConfig()
	if err != nil {
		return err
	}
	fmt.Println(cfg)
	return nil
}
