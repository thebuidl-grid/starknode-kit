package pkg

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

var (
	config     StarkNodeKitConfig
	_          = loadConfig()
	InstallDir = path.Join(getHomeDir(), "starcknode-kit")

	InstallClientsDir = path.Join(InstallDir, "ethereum_clients")

	jwtDir         = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath        = path.Join(jwtDir, "jwt.hex")
	config_dir     = path.Join(InstallDir, "config")
	yamlConfigPath = fmt.Sprintf("%s/stacknode.yaml", config_dir)
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
		return "", fmt.Errorf("Execution Client %s not supported", client) // TODO change this
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

func loadConfig() error {
	_, err := os.ReadFile(yamlConfigPath)
	if err != nil {
		config = defaultConfig()
	}
	return nil
}

func CreateStackNodeConfig() error {
	default_config := defaultConfig()
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
		ExecutionCientSettings: clientConfig{
			Name:    ClientGeth,
			Network: "sepolia",
			Port:    []string{"8545", "30303"},
		},
		ConsensusCientSettings: clientConfig{
			Name:    ClientPrysm,
			Network: "sepolia",
			Port:    []string{"8545", "30303"},
		},
	}
}
