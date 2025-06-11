package clients

import (
	"fmt"
	"starknode-kit/pkg"
	"starknode-kit/pkg/process"
	"starknode-kit/pkg/types"
	"strings"
	"time"
)

func NewConsensusClient(cfg types.ClientConfig) (types.IClient, error) {
	switch cfg.Name {
	case "lighthouse":
		return &lightHouseConfig{consensusCheckpoint: cfg.ConsensusCheckpoint, port: cfg.Port}, nil
	case "prysm":
		return &prysmConfig{consensusCheckpoint: cfg.ConsensusCheckpoint, port: cfg.Port}, nil
	default:
		return nil, fmt.Errorf("unsupported consensus client: %s", cfg.Name)
	}
}

func NewExecutionClient(cfg types.ClientConfig) (types.IClient, error) {
	switch cfg.Name {
	case "geth":
		return &gethConfig{executionType: cfg.ExecutionType, port: cfg.Port[0]}, nil
	case "reth":
		return &rethConfig{executionType: cfg.ExecutionType, port: cfg.Port[0]}, nil
	default:
		return nil, fmt.Errorf("unsupported execution client: %s", cfg.Name)
	}
}

func RestartClient(clientName string) error {
	// First stop the client
	if err := process.StopClient(clientName); err != nil {
		// If it wasn't running, that's ok
		if !strings.Contains(err.Error(), "not running") {
			return err
		}
	}

	// Wait a bit for clean shutdown
	time.Sleep(2 * time.Second)

	// Load config to get client settings
	config, err := pkg.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if clientName == "geth" || clientName == "reth" {
		c, _ := NewExecutionClient(config.ExecutionCientSettings)
		err := c.Start()
		if err != nil {
			return err
		}
	} else if clientName == "prysm" || clientName == "lighthouse" {
		c, _ := NewConsensusClient(config.ConsensusCientSettings)
		err := c.Start()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown client: %s", clientName)
	}

	return nil
}
