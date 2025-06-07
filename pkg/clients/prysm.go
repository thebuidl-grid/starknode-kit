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
	consensusPeerPorts  int
	consensusPeerPorts2 int
}

func (p prysmConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "prysm", "prysm.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "prysm", "prysm.sh")
}

// BuildGethArgs builds the arguments for the geth command
func (p *prysmConfig) buildArgs() []string {
	args := []string{
		"beacon-chain",
		"--mainnet",
		"--p2p-udp-port",
		fmt.Sprintf("--p2p-udp-port=%d", p.consensusPeerPorts2),
		fmt.Sprintf("--p2p-quic-port=%d", p.consensusPeerPorts),
		fmt.Sprintf("--p2p-tcp-port=%d", p.consensusPeerPorts),
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

func StartPrsym(port ...int) error {
	config := prysmConfig{port[0], port[1]} // TODO change
	args := config.buildArgs()
	command := config.getCommand()
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
