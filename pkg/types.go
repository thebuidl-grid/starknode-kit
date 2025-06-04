package pkg

type ClientType string

const (
	ClientGeth       ClientType = "geth"
	ClientReth       ClientType = "reth"
	ClientLighthouse ClientType = "lighthouse"
	ClientPrysm      ClientType = "prysm"
)

type StarkNodeKitConfig struct {
	ExecutionCientSettings ClientConfig `yaml:"execution_client"`
	ConsensusCientSettings ClientConfig `yaml:"consensus_client"`
}

type clientConfig struct {
	Name    ClientType `yaml:"name"`
	Network string     `yaml:"network"`
	Port    []string   `yaml:"ports"`
}
