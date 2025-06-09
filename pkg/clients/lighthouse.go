package clients

import (
	"starknode-kit/pkg"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var consensusCheckpoint = "https://mainnet-checkpoint-sync.stakely.io/"

// Configuration options for prysm
type lightHouseConfig struct {
	consensusPeerPorts  int
	consensusPeerPorts2 int
}

func (p lightHouseConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "lighthouse", "lighthouse.exe")
	}
	return filepath.Join(pkg.InstallClientsDir, "lighthouse", "lighthouse")
}

// BuildGethArgs builds the arguments for the geth command
func (p *lightHouseConfig) buildArgs() []string {
	args := []string{
		"bn",
		"--network",
		"mainnet",
		fmt.Sprintf("--port=%d", p.consensusPeerPorts),
		fmt.Sprintf("--quic-port=%d", p.consensusPeerPorts2),
		"--execution-endpoint",
		"http://localhost:8551",
		"--checkpoint-sync-url",
		consensusCheckpoint,
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

func StartLightHouse(port ...int) error {
	config := lightHouseConfig{port[0], port[1]} // TODO change
	args := config.buildArgs()
	command := config.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		pkg.InstallClientsDir,
		"lighthouse",
		"logs",
		fmt.Sprintf("lighthouse_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if err := pkg.StartProcess("lighthouse",command, logFile, args...); err != nil {
		return err
	}
	return nil
}

// TODO add log rotation
