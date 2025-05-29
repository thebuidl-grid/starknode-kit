package clients

import (
	"buidlguidl-go/pkg"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Configuration options for prysm
type prysmConfig struct {
	consensusPeerPorts  string
	consensusPeerPorts2 string
}

func (p prysmConfig) getPrysmCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "prsym", "prsym.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "prsym", "prsym")
}

// BuildGethArgs builds the arguments for the geth command
func (p *prysmConfig) buildPrysmArgs() []string {
	args := []string{
		"beacon-chain",
		"--mainnet",
		"--p2p-udp-port",
		p.consensusPeerPorts2,
		"--p2p-quic-port",
		p.consensusPeerPorts,
		"--p2p-tcp-port",
		p.consensusPeerPorts,
		"--execution-endpoint",
		"http://localhost:8551",
		"--grpc-gateway-host=0.0.0.0",
		"--grpc-gateway-port=5052",
		//		`--checkpoint-sync-url=${consensusCheckpoint}`,
		//	`--genesis-beacon-api-url=${consensusCheckpoint}`,
		"--accept-terms-of-use=true",
		"--jwt-secret",
		pkg.JWTPath,
		"--monitoring-host",
		"127.0.0.1",
		"--monitoring-port",
		"5054",
	}

	// TODO still too large
	// Add data directory
	dataDir := filepath.Join(pkg.InstallClientsDir, "prsym", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}

func StartPrsym(port ...string) error {
	config := prysmConfig{port[0], port[1]} // TODO change
	args := config.buildPrysmArgs()
	command := config.getPrysmCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		pkg.InstallClientsDir,
		"prysm",
		"logs",
		fmt.Sprintf("geth_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if err := pkg.StartProcess(command, logFile, args...); err != nil {
		return err
	}
	return nil
}

// TODO add log rotation
