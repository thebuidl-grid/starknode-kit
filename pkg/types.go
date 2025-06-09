package pkg

type ClientType string

const (
	ClientGeth       ClientType = "geth"
	ClientReth       ClientType = "reth"
	ClientLighthouse ClientType = "lighthouse"
	ClientPrysm      ClientType = "prysm"
)

type StarkNodeKitConfig struct {
	Network                string       `yaml:"network"`
	ExecutionCientSettings ClientConfig `yaml:"execution_client"`
	ConsensusCientSettings ClientConfig `yaml:"consensus_client"`
}

type ClientConfig struct {
	Name          ClientType `yaml:"name"`
	ExecutionType string     `yaml:"execution_type,omitempty"`
	Port          []int      `yaml:"ports"`
}

type Process struct {
	Processes []process `yaml:"processes"`
}

type process struct {
	Name ClientType `yaml:"name"`
	Type string     `yaml:"type"`
	Pid  int        `yaml:"pid"`
	Path string     `yaml:"path"`
}
