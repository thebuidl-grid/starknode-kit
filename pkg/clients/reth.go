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
type RethConfig struct {
	ExecutionType     string
	ExecutionPeerPort int
	LogFilePath       string
}

// GetRethCommand returns the reth command path based on platform
func GetRethCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "reth", "reth.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "reth", "reth")
}

// BuildRethArgs builds the arguments for the reth command
func BuildRethArgs(config *RethConfig) []string {
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
		"--port", fmt.Sprintf("%d", config.ExecutionPeerPort),
		"--metrics", "0.0.0.0:6060",
	}

	// Add execution type specific arguments
	if config.ExecutionType == "archive" {
		args = append(args, "--archive")
	}

	// Add data directory
	dataDir := filepath.Join(pkg.InstallClientsDir, "ethereum_clients", "reth", "database")
	args = append(args, "--datadir", dataDir)

	return args
}

func StartReth(executionType string, port []int) error {
	config := GethConfig{ExecutionPeerPort: port[0], ExecutionType: executionType}
	args := buildGethArgs(&config)
	command := GetGethCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		pkg.InstallClientsDir,
		"reth",
		"logs",
		fmt.Sprintf("reth_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if err := process.StartProcess("reth", command, logFile, args...); err != nil {
		return err
	}
	return nil
}
