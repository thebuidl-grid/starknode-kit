package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
)

// JunoClient represents a client for interacting with a local Juno node
type JunoClient struct {
	config          types.JunoConfig
	isValidatorNode bool
	network         string
}

// getJunoPath returns the path to the Juno binary
func getJunoPath() string {
	// Check if Juno is installed in the github.com/thebuidl-grid/starknode-kit directory
	junoDir := filepath.Join(constants.InstallStarknetDir, "juno")
	junoPath := filepath.Join(junoDir, "juno", "build", "juno")

	if _, err := os.Stat(junoPath); err == nil {
		return junoPath
	}

	return ""
}

// StartNode starts a local Juno node
func (c *JunoClient) Start() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		constants.InstallStarknetDir,
		"juno",
		"logs",
		fmt.Sprintf("juno_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	args := c.buildJunoArgs()
	return process.StartClient("juno", getJunoPath(), logFile, args...)
}

// buildJunoArgs builds the command line arguments for Juno
func (c *JunoClient) buildJunoArgs() []string {
	args := []string{
		"--http",
		fmt.Sprintf("--http-port=%d", c.config.Port),
		"--http-host=0.0.0.0",
		fmt.Sprintf("--db-path=%s", filepath.Join(constants.InstallStarknetDir, "juno", "database")),
		fmt.Sprintf("--eth-node=%s", c.config.EthNode),
		fmt.Sprintf("--ws=%t", c.isValidatorNode),
		"--ws-port=6061",
		"--ws-host=0.0.0.0",
	}

	// Add network configuration
	if c.network == "mainnet" {
		args = append(args, "--network=mainnet")
	} else if c.network == "sepolia" {
		args = append(args, "--network=sepolia")
	} else if c.network == "sepolia-integration" {
		args = append(args, "--network=sepolia-integration")
	}

	return args
}
