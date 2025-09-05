package clients

import (
	"fmt"
	"strings"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

func NewConsensusClient(cfg types.ClientConfig, network string) (types.IClient, error) {
	switch cfg.Name {
	case "lighthouse":
		return &lightHouseConfig{consensusCheckpoint: cfg.ConsensusCheckpoint, port: cfg.Port, network: network}, nil
	case "prysm":
		return &prysmConfig{consensusCheckpoint: cfg.ConsensusCheckpoint, port: cfg.Port, network: network}, nil
	default:
		return nil, fmt.Errorf("unsupported consensus client: %s", cfg.Name)
	}
}

func NewExecutionClient(cfg types.ClientConfig, network string) (types.IClient, error) {
	switch cfg.Name {
	case "geth":
		return &gethConfig{executionType: cfg.ExecutionType, port: cfg.Port[0], network: network}, nil
	case "reth":
		return &rethConfig{executionType: cfg.ExecutionType, port: cfg.Port[0], network: network}, nil
	default:
		return nil, fmt.Errorf("unsupported execution client: %s", cfg.Name)
	}
}

func NewJunoClient(config types.JunoConfig, network string) (types.IClient, error) {

	// Get Juno binary path
	junoPath := getJunoPath()
	if junoPath == "" {
		return nil, fmt.Errorf("Juno is not installed. Please install it first using 'starknode-kit add -s juno'")
	}

	return &JunoClient{
		config:  config,
		network: network,
	}, nil
}

func NewValidatorClient(config types.ValidatorConfig) (types.IClient, error) {

	return &StakingValidator{
		Provider: stakingValidatorProviderConfig{
			starknetHttp: config.ProviderConfig.JunoRPC,
			starkentWS:   config.ProviderConfig.JunoWS,
		},
		Wallet: stakingValidatorWalletConfig{
			address:    config.SignerConfig.OperationalAddress,
			privatekey: config.SignerConfig.WalletPrivateKey,
		},
	}, nil
}
func RestartClient(pid int) error {
	// First stop the client
	if err := process.StopClient(pid); err != nil {
		// If it wasn't running, that's ok
		if !strings.Contains(err.Error(), "not running") {
			return err
		}
	}

	// Wait a bit for clean shutdown
	time.Sleep(2 * time.Second)

	// Load config to get client settings
	config, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	e, _ := NewExecutionClient(config.ExecutionCientSettings, config.Network)
	err = e.Start()
	if err != nil {
		return err
	}
	c, _ := NewConsensusClient(config.ConsensusCientSettings, config.Network)
	err = c.Start()
	if err != nil {
		return err
	}

	j, _ := NewJunoClient(config.JunoConfig, config.Network)
	err = j.Start()
	if err != nil {
		return err
	}

	return nil
}
