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
		Network                string          `yaml:"network"`
		Wallet                 WalletConfig    `yaml:"wallet"`
		ExecutionCientSettings ClientConfig    `yaml:"execution_client"`
		ConsensusCientSettings ClientConfig    `yaml:"consensus_client"`
		JunoConfig             JunoConfig      `yaml:"juno_client,omitempty"`
		ValidatorConfig        ValidatorConfig `yaml:"validator_config"`
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
		Name          string `yaml:"name"`
		RewardAddress string `yaml:"reward_address"`
		Wallet        Wallet `yaml:"wallet"`
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

	ValidatorConfig struct {
		ProviderConfig struct {
			JunoRPC string `json:"http" yaml:"juno_rpc_http"`
			JunoWS  string `json:"ws" yaml:"juno_rpc_ws"`
		} `json:"provider" yaml:"provider_config"`
		SignerConfig struct {
			OperationalAddress string `json:"operational_address"`
			WalletPrivateKey   string `json:"privateKey"`
		} `json:"signer" yaml:"signer"`
	}
)
