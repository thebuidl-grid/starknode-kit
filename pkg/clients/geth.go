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

// Configuration options for Geth
type gethConfig struct {
	port          int
	executionType string
}

// GetGethCommand returns the geth command path based on platform
func (_ gethConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "geth", "geth.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "geth", "geth")
}

// BuildGethArgs builds the arguments for the geth command
func (c *gethConfig) buildArgs() []string {
	args := []string{
		"--mainnet",
		fmt.Sprintf("--port=%d", c.port),
		fmt.Sprintf("--discovery.port=%d", c.port),
		"--http",
		"--http.api=eth,net,engine,admin",
		"--http.corsdomain=*",
		"--http.addr=0.0.0.0",
		"--http.port=8545",
		"--authrpc.jwtsecret=" + pkg.JWTPath,
		"--authrpc.addr=0.0.0.0",
		"--authrpc.port=8551",
		"--authrpc.vhosts=*",
		"--metrics",
		"--metrics.addr=0.0.0.0",
		"--metrics.port=6060",
	}

	// Add execution type specific arguments
	if c.executionType == "full" {
		args = append(args, "--syncmode=snap")
	} else if c.executionType == "archive" {
		args = append(args, "--syncmode=full", "--gcmode=archive")
	}

	// Add data directory
	dataDir := filepath.Join(pkg.InstallClientsDir, "geth", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}

func (c *gethConfig) Start() error {
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

	return process.StartProcess("lighthouse", command, logFile, args...)
}
