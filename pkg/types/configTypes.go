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

type (
	StarkNodeKitConfig struct {
		Network                string       `yaml:"network"`
		IsValidatorNode        bool         `yaml:"is_validator_node"`
		Wallet                 WalletConfig `wallet:"wallet"`
		ExecutionCientSettings ClientConfig `yaml:"execution_client"`
		ConsensusCientSettings ClientConfig `yaml:"consensus_client"`
		JunoConfig             JunoConfig   `yaml:"juno_client"`
	}

	ClientConfig struct {
		ExecutionType       string     `yaml:"execution_type,omitempty"`
		Port                []int      `yaml:"ports"`
		ConsensusCheckpoint string     `yaml:"consensus_checkpoint,omitempty"`
		Name                ClientType `yaml:"name"`
	}

	JunoConfig struct {
		Port        int      `yaml:"port"`
		EthNode     string   `yaml:"eth_node"`
		Environment []string `yaml:"environment"`
	}

	WalletConfig struct {
		WalletAddress string `yaml:"wallet_address"`
		PrivateKey    string `yaml:"private_key"`
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
