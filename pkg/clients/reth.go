package clients

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
)

// Configuration options for Reth
type rethConfig struct {
	port          int
	executionType string
	network       string
}

// GetRethCommand returns the reth command path based on platform
func (_ rethConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(constants.InstallClientsDir, "reth", "reth.exe")
	}
	return filepath.Join(constants.InstallClientsDir, "reth", "reth")
}

// BuildRethArgs builds the arguments for the reth command
func (config *rethConfig) buildArgs() []string {
	// Build common arguments
	args := []string{
		"node",
		"--chain", config.network,
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.port", "8545",
		"--http.api", "eth,net,admin",
		"--http.corsdomain", "*",
		"--authrpc.addr", "0.0.0.0",
		"--authrpc.port", "8551",
		"--authrpc.jwtsecret", constants.JWTPath,
		"--port", fmt.Sprintf("%d", config.port),
		"--metrics", "0.0.0.0:7878",
	}

	// Add execution type specific arguments
	if config.executionType == "archive" {
		args = append(args, "--archive")
	}

	// Add data directory
	dataDir := filepath.Join(constants.InstallClientsDir, "reth", "database")
	args = append(args, "--datadir", dataDir)

	return args
}

func (c *rethConfig) Start() error {
	args := c.buildArgs()
	command := c.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		constants.InstallClientsDir,
		"reth",
		"logs",
		fmt.Sprintf("reth_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	return process.StartClient("reth", command, multiWriter, args...)
}