package types

type StarkNodeKitConfig struct {
	ExecutionCientSettings ClientConfig `yaml:"execution_client"`
	ConsensusCientSettings ClientConfig `yaml:"consensus_client"`
}

type ClientConfig struct {
	Name    string   `yaml:"name"`
	Network string   `yaml:"network"`
	Port    []string `yaml:"ports"`
}
