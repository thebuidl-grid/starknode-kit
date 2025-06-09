package clients

import (
	"fmt"
	"starknode-kit/pkg/types"
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
