package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/spf13/pflag"
)

// Configuration options for Geth
type GethConfig struct {
	InstallDir        string
	ExecutionType     string
	ExecutionPeerPort int
	JWTPath           string
	LogFilePath       string
}

// StartGeth starts the Geth client with the given configuration
func StartGeth(config *GethConfig) error {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(config.LogFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("error creating log directory: %w", err)
	}

	// Open log file
	logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}
	defer logFile.Close()

	// Get geth command path based on platform
	gethCmd := GetGethCommand(config.InstallDir)
	if _, err := os.Stat(gethCmd); os.IsNotExist(err) {
		return fmt.Errorf("geth binary not found at %s. Please install geth first", gethCmd)
	}

	// Build geth command arguments
	args := BuildGethArgs(config)

	// Print the command that will be executed
	fmt.Printf("Launching geth with command: %s %s\n", gethCmd, strings.Join(args, " "))
	fmt.Printf("Logs will be written to: %s\n", config.LogFilePath)

	// Create the command
	cmd := exec.Command(gethCmd, args...)

	// Start geth in a pseudo-terminal
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("error starting geth: %w", err)
	}
	defer ptmx.Close()

	// Set up a channel to capture interrupt signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Handle terminal resize
	go func() {
		for {
			time.Sleep(time.Second)
			if ptmx == nil {
				return
			}
		}
	}()

	// Handle process output
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := ptmx.Read(buffer)
			if err != nil {
				fmt.Printf("Error reading from geth: %v\n", err)
				return
			}

			output := string(buffer[:n])

			// Write to log file
			if _, err := logFile.WriteString(output); err != nil {
				fmt.Printf("Error writing to log file: %v\n", err)
			}

			// Write to stdout
			fmt.Print(output)
		}
	}()

	// Wait for a signal to terminate
	<-sigs
	fmt.Println("\nShutting down geth...")

	// Kill the process
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Error sending SIGTERM to geth: %v\n", err)
		cmd.Process.Kill()
	}

	// Wait for process to exit
	cmd.Wait()

	return nil
}

// ParseGethFlags parses command line flags for Geth configuration
func ParseGethFlags() *GethConfig {
	return ParseGethFlagsWithFlagSet(pflag.CommandLine)
}

// ParseGethFlagsWithFlagSet parses command line flags for Geth configuration using a specific flag set
func ParseGethFlagsWithFlagSet(flagSet *pflag.FlagSet) *GethConfig {
	// Default config
	config := &GethConfig{
		ExecutionType:     "full",
		ExecutionPeerPort: 30303,
	}

	// Get default home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	config.InstallDir = home

	// Define flags only if they don't exist
	if flagSet.Lookup("directory") == nil {
		flagSet.StringVar(&config.InstallDir, "directory", config.InstallDir, "Installation directory")
	}
	if flagSet.Lookup("executiontype") == nil {
		flagSet.StringVar(&config.ExecutionType, "executiontype", config.ExecutionType, "Execution type (full or archive)")
	}
	if flagSet.Lookup("executionpeerport") == nil {
		flagSet.IntVar(&config.ExecutionPeerPort, "executionpeerport", config.ExecutionPeerPort, "Execution peer port")
	}

	// Parse flags
	flagSet.Parse(os.Args[1:])

	// Update config from parsed flags
	if dir, err := flagSet.GetString("directory"); err == nil {
		config.InstallDir = dir
	}
	if execType, err := flagSet.GetString("executiontype"); err == nil {
		config.ExecutionType = execType
	}
	if port, err := flagSet.GetInt("executionpeerport"); err == nil {
		config.ExecutionPeerPort = port
	}

	// Set derived config values
	config.JWTPath = filepath.Join(config.InstallDir, "ethereum_clients", "jwt", "jwt.hex")

	// Create a timestamped log file name
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	config.LogFilePath = filepath.Join(
		config.InstallDir,
		"ethereum_clients",
		"geth",
		"logs",
		fmt.Sprintf("geth_%s.log", timestamp),
	)

	return config
}

// GetGethCommand returns the geth command path based on platform
func GetGethCommand(installDir string) string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(installDir, "ethereum_clients", "geth", "geth.exe")
	}
	return filepath.Join(installDir, "ethereum_clients", "geth", "geth")
}

// BuildGethArgs builds the arguments for the geth command
func BuildGethArgs(config *GethConfig) []string {
	args := []string{
		"--mainnet",
		fmt.Sprintf("--port=%d", config.ExecutionPeerPort),
		fmt.Sprintf("--discovery.port=%d", config.ExecutionPeerPort),
		"--http",
		"--http.api=eth,net,engine,admin",
		"--http.corsdomain=*",
		"--http.addr=0.0.0.0",
		"--http.port=8545",
		"--authrpc.jwtsecret=" + config.JWTPath,
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
	dataDir := filepath.Join(config.InstallDir, "ethereum_clients", "geth", "database")
	args = append(args, "--datadir="+dataDir)

	return args
}
