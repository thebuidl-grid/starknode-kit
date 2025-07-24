package types

type ClientType string

const (
	ClientGeth       ClientType = "geth"
	ClientReth       ClientType = "reth"
	ClientLighthouse ClientType = "lighthouse"
	ClientPrysm      ClientType = "prysm"
	ClientJuno       ClientType = "juno"
)

func GetClientType(client string) ClientType {
	switch client {
	case "geth":
		return ClientGeth
	case "reth":
		return ClientReth
	case "lighthouse":
		return ClientLighthouse
	case "prysm":
		return ClientPrysm
	case "juno":
		return ClientJuno
	default:
		return ""
	}
}

type IClient interface {
	Start() error
}

type StarkNodeKitConfig struct {
	WalletAddress          string       `yaml:"wallet_address"`
	PrivateKey             string       `yaml:"private_key"`
	Network                string       `yaml:"network"`
	ExecutionCientSettings ClientConfig `yaml:"execution_client"`
	ConsensusCientSettings ClientConfig `yaml:"consensus_client"`
	JunoConfig             JunoConfig   `yaml:"juno_client"`
}

type ClientConfig struct {
	Name                ClientType `yaml:"name"`
	ExecutionType       string     `yaml:"execution_type,omitempty"`
	Port                []int      `yaml:"ports"`
	ConsensusCheckpoint string     `yaml:"consensus_checkpoint,omitempty"`
}

type JunoConfig struct {
	Port        int      `yaml:"port"`
	EthNode     string   `yaml:"eth_node"`
	Environment []string `yaml:"environment"`
}
