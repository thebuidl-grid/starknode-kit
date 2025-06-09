package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"starknode-kit/pkg"
	"time"
)

// Configuration options for Geth
type GethConfig struct {
	ExecutionType     string
	ExecutionPeerPort int
	LogFilePath       string
}

// GetGethCommand returns the geth command path based on platform
func GetGethCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "geth", "geth.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "geth", "geth")
}

// BuildGethArgs builds the arguments for the geth command
func buildGethArgs(config *GethConfig) []string {
	args := []string{
		"--mainnet",
		fmt.Sprintf("--port=%d", config.ExecutionPeerPort),
		fmt.Sprintf("--discovery.port=%d", config.ExecutionPeerPort),
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
	if config.ExecutionType == "full" {
		args = append(args, "--syncmode=snap")
	} else if config.ExecutionType == "archive" {
		args = append(args, "--syncmode=full", "--gcmode=archive")
	}

	// Add data directory
	dataDir := filepath.Join(pkg.InstallClientsDir, "geth", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}
func StartGeth(executionType string, port []int) error {
	config := GethConfig{ExecutionPeerPort: port[0], ExecutionType: executionType}
	args := buildGethArgs(&config)
	command := GetGethCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		pkg.InstallClientsDir,
		"geth",
		"logs",
		fmt.Sprintf("geth_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if err := pkg.StartProcess("geth", command, logFile, args...); err != nil {
		return err
	}
	return nil
}
