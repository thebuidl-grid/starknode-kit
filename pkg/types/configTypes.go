package types

type ClientType string

func (c ClientType) String() string {
	return string(c)
}

const (
	ClientGeth       ClientType = "geth"
	ClientReth       ClientType = "reth"
	ClientLighthouse ClientType = "lighthouse"
	ClientPrysm      ClientType = "prysm"
	ClientJuno       ClientType = "juno"
)

type IClient interface {
	Start() error
}

type StarkNodeKitConfig struct {
	WalletAddress          string       `yaml:"wallet_address"`
	PrivateKey             string       `yaml:"private_key"`
	Network                string       `yaml:"network"`
	ExecutionCientSettings ClientConfig `yaml:"execution_client"`
	ConsensusCientSettings ClientConfig `yaml:"consensus_client"`
}

type ClientConfig struct {
	Name                ClientType `yaml:"name"`
	Network             string     `yaml:"network,omitempty"`
	ExecutionType       string     `yaml:"execution_type,omitempty"`
	Port                []int      `yaml:"ports"`
	ConsensusCheckpoint string     `yaml:"consensus_checkpoint,omitempty"`
}
