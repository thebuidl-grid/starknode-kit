package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"starknode-kit/pkg"
	"starknode-kit/pkg/process"
	"time"
)

// Configuration options for Reth
type rethConfig struct {
	port          int
	executionType string
}

// GetRethCommand returns the reth command path based on platform
func (_ rethConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "reth", "reth.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "reth", "reth")
}

// BuildRethArgs builds the arguments for the reth command
func (config *rethConfig) buildArgs() []string {
	// Build common arguments
	args := []string{
		"node",
		"--network", "mainnet",
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.port", "8545",
		"--http.api", "eth,net,engine,admin",
		"--http.corsdomain", "*",
		"--authrpc.addr", "0.0.0.0",
		"--authrpc.port", "8551",
		"--authrpc.jwtsecret", pkg.JWTPath,
		"--port", fmt.Sprintf("%d", config.port),
		"--metrics", "0.0.0.0:6060",
	}

	// Add execution type specific arguments
	if config.executionType == "archive" {
		args = append(args, "--archive")
	}

	// Add data directory
	dataDir := filepath.Join(pkg.InstallClientsDir, "ethereum_clients", "reth", "database")
	args = append(args, "--datadir", dataDir)

	return args
}

func (c *rethConfig) Start() error {
	args := c.buildArgs()
	command := c.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		pkg.InstallClientsDir,
		"geth",
		"logs",
		fmt.Sprintf("geth_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	return process.StartClient("lighthouse", command, logFile, args...)
}
