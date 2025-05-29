package clients 

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

// Configuration options for Reth
type RethConfig struct {
	InstallDir        string
	ExecutionType     string
	ExecutionPeerPort int
	JWTPath           string
	LogFilePath       string
}

// StartReth starts the Reth client with the given configuration
func StartReth(config *RethConfig) error {
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

	// Get reth command path based on platform
	rethCmd := GetRethCommand(config.InstallDir)
	if _, err := os.Stat(rethCmd); os.IsNotExist(err) {
		return fmt.Errorf("reth binary not found at %s. Please install reth first", rethCmd)
	}

	// Build reth command arguments
	args := BuildRethArgs(config)

	// Print the command that will be executed
	fmt.Printf("Launching reth with command: %s %s\n", rethCmd, strings.Join(args, " "))
	fmt.Printf("Logs will be written to: %s\n", config.LogFilePath)

	// Create the command
	cmd := exec.Command(rethCmd, args...)

	// Start reth in a pseudo-terminal
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("error starting reth: %w", err)
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
				fmt.Printf("Error reading from reth: %v\n", err)
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
	fmt.Println("\nShutting down reth...")

	// Kill the process
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Error sending SIGTERM to reth: %v\n", err)
		cmd.Process.Kill()
	}

	// Wait for process to exit
	cmd.Wait()

	return nil
}

// ParseRethFlags parses command line flags for Reth configuration
func ParseRethFlags() *RethConfig {
	return ParseRethFlagsWithFlagSet(pflag.CommandLine)
}

// ParseRethFlagsWithFlagSet parses command line flags for Reth configuration using a specific flag set
func ParseRethFlagsWithFlagSet(flagSet *pflag.FlagSet) *RethConfig {
	// Default config
	config := &RethConfig{
		ExecutionType:     "full",
		ExecutionPeerPort: 30303,
	}

	// Get default home directory
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return nil // Return nil config instead of os.Exit(1)
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
		"reth",
		"logs",
		fmt.Sprintf("reth_%s.log", timestamp),
	)

	return config
}

// GetRethCommand returns the reth command path based on platform
func GetRethCommand(installDir string) string {
	platform := runtime.GOOS
	if platform == "windows" {
		return filepath.Join(installDir, "ethereum_clients", "reth", "reth.exe")
	}
	return filepath.Join(installDir, "ethereum_clients", "reth", "reth")
}

// BuildRethArgs builds the arguments for the reth command
func BuildRethArgs(config *RethConfig) []string {
	// Build common arguments
	args := []string{
		"node",
		"--network", "mainnet",
		"--http",
		"--http.addr", "0.0.0.0",
		"--http.port", "8545",
		"--http.api", "eth,net,engine,admin",
		"--http.corsdomain", "*",
		"--authrpc.addr", "0.0.0.0",
		"--authrpc.port", "8551",
		"--authrpc.jwtsecret", config.JWTPath,
		"--port", fmt.Sprintf("%d", config.ExecutionPeerPort),
		"--metrics", "0.0.0.0:6060",
	}

	// Add execution type specific arguments
	if config.ExecutionType == "archive" {
		args = append(args, "--archive")
	}

	// Add data directory
	dataDir := filepath.Join(config.InstallDir, "ethereum_clients", "reth", "database")
	args = append(args, "--datadir", dataDir)

	return args
}
