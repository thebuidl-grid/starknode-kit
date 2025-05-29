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
type lightHouseConfig struct {
	consensusPeerPorts  string
	consensusPeerPorts2 string
}

func (p lightHouseConfig) getCommand() string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(pkg.InstallClientsDir, "prsym", "prsym.sh")
	}
	return filepath.Join(pkg.InstallClientsDir, "prsym", "prsym")
}

// BuildGethArgs builds the arguments for the geth command
func (p *lightHouseConfig) buildArgs() []string {
	args := []string{
		"bn",
		"--network",
		"mainnet",
		"--port",
		p.consensusPeerPorts,
		"--quic-port",
		p.consensusPeerPorts2,
		"--execution-endpoint",
		"http://localhost:8551",
		//"--checkpoint-sync-url",
		// consensusCheckpoint,
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
	dataDir := filepath.Join(pkg.InstallClientsDir, "prsym", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}

func StartLighHouse(port ...string) error {
	config := prysmConfig{port[0], port[1]} // TODO change
	args := config.buildArgs()
	command := config.getCommand()
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(
		pkg.InstallClientsDir,
		"lighhouse",
		"logs",
		fmt.Sprintf("lighhouse_%s.log", timestamp))
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
