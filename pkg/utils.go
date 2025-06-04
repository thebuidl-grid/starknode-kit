package pkg

import (
	"buidlguidl-go/pkg/types"
	"fmt"
	"os"
	"path"
)

// ClientType represents an Ethereum client type
type ClientType string

const (
	ClientGeth       ClientType = "geth"
	ClientReth       ClientType = "reth"
	ClientLighthouse ClientType = "lighthouse"
	ClientPrysm      ClientType = "prysm"
)

var (
	InstallDir = path.Join(getHomeDir(), "starcknode-kit")

	InstallClientsDir = path.Join(InstallDir, "ethereum_clients")

	jwtDir  = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath = path.Join(jwtDir, "jwt.hex")
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
		return "", fmt.Errorf("Execution Client %s not supported", client)
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

func DefaultConfig() types.StarkNodeKitConfig {
	return types.StarkNodeKitConfig{
		ExecutionCientSettings: types.ClientConfig{
			Name:    "reth",
			Network: "sepolia",
			Port:    []string{"8545", "30303"},
		},
		ConsensusCientSettings: types.ClientConfig{
			Name:    "lighthouse",
			Network: "sepolia",
			Port:    []string{"8545", "30303"},
		},
	}
}
