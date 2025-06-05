package pkg

type ClientType string

const (
	ClientGeth       ClientType = "geth"
	ClientReth       ClientType = "reth"
	ClientLighthouse ClientType = "lighthouse"
	ClientPrysm      ClientType = "prysm"
)

type StarkNodeKitConfig struct {
	ExecutionCientSettings ClientSettings `yaml:"execution_client"`
	ConsensusCientSettings ClientSettings `yaml:"consensus_client"`
}

type ClientSettings struct {
	Name    ClientType `yaml:"name"`
	Network string     `yaml:"network"`
	Port    []string   `yaml:"ports"`
}
