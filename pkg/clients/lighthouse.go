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

// Configuration options for prysm
type lightHouseConfig struct {
	port                []int // [quic/tcp, udp]
	consensusCheckpoint string
	network             string
}

func (_ lightHouseConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "lighthouse", "lighthouse.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "lighthouse", "lighthouse")
}

// BuildGethArgs builds the arguments for the geth command
func (c *lightHouseConfig) buildArgs() []string {
	args := []string{
		"bn",
		"--network",
		c.network,
		fmt.Sprintf("--port=%d", c.port[0]),
		fmt.Sprintf("--quic-port=%d", c.port[1]),
		"--execution-endpoint",
		"http://localhost:8551",
		"--checkpoint-sync-url",
		c.consensusCheckpoint,
		"--checkpoint-sync-url-timeout",
		"1200",
		"--disable-deposit-contract-sync",
		"--execution-jwt",
		pkg.JWTPath,
		"--metrics",
		"--metrics-address",
		"127.0.0.1",
		"--metrics-port",
		"5054",
		"--http",
		"--disable-upnp", // There is currently a bug in the p2p-lib that causes panics with this enabled
	}

	// TODO still too large
	// Add data directory
	dataDir := filepath.Join(pkg.InstallClientsDir, "lighthouse", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}

func (c *lightHouseConfig) Start() error {
	args := c.buildArgs()
	command := c.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(pkg.InstallClientsDir, "lighthouse", "logs", fmt.Sprintf("lighthouse_%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	return process.StartClient("lighthouse", command, logFile, args...)
}
