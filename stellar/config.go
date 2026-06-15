package stellar

// Network identifiers
const (
	NetworkPubnet  = "pubnet"
	NetworkTestnet = "testnet"
)

// DefaultHTTPPort is the stellar-core local admin HTTP port
const DefaultHTTPPort = 11626

// Config holds all settings needed to run a stellar-core node
type Config struct {
	// Network is "pubnet" (mainnet) or "testnet"
	Network string `yaml:"network"`
	// ConfigPath is the path to stellar-core.cfg; defaults to <InstallDir>/stellar-core.cfg
	ConfigPath string `yaml:"config_path,omitempty"`
	// HTTPPort is the local admin HTTP port stellar-core listens on (default 11626)
	HTTPPort int `yaml:"http_port,omitempty"`
}

func (c *Config) httpPort() int {
	if c.HTTPPort == 0 {
		return DefaultHTTPPort
	}
	return c.HTTPPort
}

func (c *Config) configPath() string {
	if c.ConfigPath != "" {
		return c.ConfigPath
	}
	return installDir() + "/stellar-core.cfg"
}
