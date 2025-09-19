package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
)

// Configuration options for prysm
type prysmConfig struct {
	port                []int // [quic/tcp, udp]
	consensusCheckpoint string
	network             string
}

func (_ prysmConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(constants.InstallClientsDir, "prysm", "prysm.exe")
	}
	return filepath.Join(constants.InstallClientsDir, "prysm", "prysm.sh")
}

// BuildGethArgs builds the arguments for the geth command
func (c *prysmConfig) buildArgs() []string {
	args := []string{
		"beacon-chain",
		fmt.Sprintf("--%s", c.network),
		fmt.Sprintf("--p2p-udp-port=%d", c.port[1]),
		fmt.Sprintf("--p2p-quic-port=%d", c.port[0]),
		fmt.Sprintf("--p2p-tcp-port=%d", c.port[0]),
		"--execution-endpoint",
		"http://localhost:8551",
		"--grpc-gateway-host=0.0.0.0",
		"--grpc-gateway-port=5052",
		fmt.Sprintf("--checkpoint-sync-url=%s", c.consensusCheckpoint),
		fmt.Sprintf("--genesis-beacon-api-url=%s", c.consensusCheckpoint),
		"--accept-terms-of-use=true",
		"--jwt-secret",
		constants.JWTPath,
		"--monitoring-host",
		"127.0.0.1",
		"--monitoring-port",
		"5054",
	}

	// TODO still too large
	// Add data directory
	dataDir := filepath.Join(constants.InstallClientsDir, "prsym", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}

// prysm.go
func (c *prysmConfig) Start() error {
	args := c.buildArgs()
	command := c.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(constants.InstallClientsDir, "prysm", "logs", fmt.Sprintf("prysm_%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	return process.StartClient("prysm", command, logFile, args...)
}

