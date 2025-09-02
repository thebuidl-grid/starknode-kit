package types

type ClientType string

const (
	ClientGeth           ClientType = "geth"
	ClientReth           ClientType = "reth"
	ClientLighthouse     ClientType = "lighthouse"
	ClientPrysm          ClientType = "prysm"
	ClientJuno           ClientType = "juno"
	ClientStarkValidator ClientType = "starknet-staking-v2"
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
	case "starknet-staking-v2":
		return ClientJuno
	default:
		return ""
	}
}

type IClient interface {
	Start() error
}

type (
	StarkNodeKitConfig struct {
		Network                string         `yaml:"network"`
		Wallet                 WalletConfig   `wallet:"wallet,omitempty"`
		ExecutionCientSettings ClientConfig   `yaml:"execution_client"`
		ConsensusCientSettings ClientConfig   `yaml:"consensus_client"`
		JunoConfig             JunoConfig     `yaml:"juno_client,omitempty"`
		Wallets                []WalletConfig `yaml:"wallet"`
	}

	ClientConfig struct {
		ExecutionType       string     `yaml:"execution_type,omitempty"`
		Port                []int      `yaml:"ports"`
		ConsensusCheckpoint string     `yaml:"consensus_checkpoint,omitempty"`
		Name                ClientType `yaml:"name"`
	}

	JunoConfig struct {
		Port            int      `yaml:"port"`
		EthNode         string   `yaml:"eth_node"`
		Environment     []string `yaml:"environment"`
		IsValidatorNode bool     `yaml:"is_validator_node"`
	}

	WalletConfig struct {
		Name   string `yaml:"name"`
		Wallet Wallet
	}

	Wallet struct {
		Address    string `json:"address"`
		ClassHash  string `json:"class_hash"`
		Deployed   bool   `json:"deployed"`
		Legacy     bool   `json:"legacy"`
		PrivateKey string `json:"private_key"`
		PublicKey  string `json:"public_key"`
		Salt       string `json:"salt"`
	}
)
