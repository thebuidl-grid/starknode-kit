package clients

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"starknode-kit/pkg"
)

// JunoConfig represents the configuration for a Juno node
type JunoConfig struct {
	Network     string   `yaml:"network"`
	Port        string   `yaml:"port"`
	UseSnapshot bool     `yaml:"use_snapshot"`
	DataDir     string   `yaml:"data_dir"`
	EthNode     string   `yaml:"eth_node"`
	Environment []string `yaml:"environment"`
}

// DefaultJunoConfig returns a default configuration for Juno
func DefaultJunoConfig() *JunoConfig {
	return &JunoConfig{
		Network:     "mainnet",
		Port:        "6060",
		UseSnapshot: true,
		DataDir:     pkg.JunoDataDir,
		EthNode:     "ws://localhost:8546",
		Environment: []string{
			"JUNO_NETWORK=mainnet",
			"JUNO_HTTP_PORT=6060",
			"JUNO_HTTP_HOST=0.0.0.0",
		},
	}
}

// JunoClient represents a client for interacting with a local Juno node
type JunoClient struct {
	config   *JunoConfig
	process  *os.Process
	junoPath string
	logFile  *os.File
}

// NewJunoClient creates a new Juno client instance
func NewJunoClient(config *JunoConfig) (*JunoClient, error) {
	if config == nil {
		config = DefaultJunoConfig()
	}

	// Get Juno binary path
	junoPath := getJunoPath()
	if junoPath == "" {
		return nil, fmt.Errorf("Juno is not installed. Please install it first using 'starknode add -s juno'")
	}

	return &JunoClient{
		config:   config,
		junoPath: junoPath,
	}, nil
}

// getJunoPath returns the path to the Juno binary
func getJunoPath() string {
	// Check if Juno is installed in the starknode-kit directory
	junoDir := filepath.Join(pkg.InstallClientsDir, "juno")
	junoPath := filepath.Join(junoDir, "juno")

	if _, err := os.Stat(junoPath); err == nil {
		return junoPath
	}

	// Check if Juno is available globally
	if path, err := exec.LookPath("juno"); err == nil {
		return path
	}

	return ""
}

// StartNode starts a local Juno node
func (c *JunoClient) StartNode(ctx context.Context) error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(c.config.DataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create logs directory
	logsDir := filepath.Join(filepath.Dir(c.config.DataDir), "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logsDir, fmt.Sprintf("juno_%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	c.logFile = logFile

	// Build Juno command arguments
	args := c.buildJunoArgs()

	// Start Juno process
	cmd := exec.CommandContext(ctx, c.junoPath, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.Dir = c.config.DataDir

	// Set environment variables
	cmd.Env = append(os.Environ(), c.config.Environment...)

	fmt.Printf("Starting Juno node with command: %s %v\n", c.junoPath, args)

	if err := cmd.Start(); err != nil {
		logFile.Close()
		return fmt.Errorf("failed to start Juno node: %w", err)
	}

	c.process = cmd.Process

	// Wait a bit for the process to start
	time.Sleep(2 * time.Second)

	// Check if process is still running
	if c.process == nil || c.process.Pid == 0 {
		logFile.Close()
		return fmt.Errorf("Juno process failed to start")
	}

	fmt.Printf("Juno node started with PID: %d\n", c.process.Pid)
	return nil
}

// buildJunoArgs builds the command line arguments for Juno
func (c *JunoClient) buildJunoArgs() []string {
	args := []string{
		"--http",
		fmt.Sprintf("--http-port=%s", c.config.Port),
		"--http-host=0.0.0.0",
		fmt.Sprintf("--db-path=%s", c.config.DataDir),
		fmt.Sprintf("--eth-node=%s", c.config.EthNode),
	}

	// Add network configuration
	if c.config.Network == "mainnet" {
		args = append(args, "--network=mainnet")
	} else if c.config.Network == "sepolia" {
		args = append(args, "--network=sepolia")
	} else if c.config.Network == "sepolia-integration" {
		args = append(args, "--network=sepolia-integration")
	}

	// Add snapshot flag if enabled
	if c.config.UseSnapshot {
		args = append(args, "--snapshot")
	}

	// Add metrics endpoint
	args = append(args, "--metrics", "--metrics-port=6060")

	return args
}

// StopNode stops the running Juno node
func (c *JunoClient) StopNode(ctx context.Context) error {
	if c.process == nil {
		return nil
	}

	fmt.Printf("Stopping Juno node (PID: %d)...\n", c.process.Pid)

	// Send SIGTERM first
	if err := c.process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to send interrupt signal: %w", err)
	}

	// Wait for graceful shutdown
	done := make(chan error, 1)
	go func() {
		_, err := c.process.Wait()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("process wait error: %w", err)
		}
	case <-time.After(10 * time.Second):
		// Force kill if graceful shutdown takes too long
		fmt.Println("Force killing Juno process...")
		if err := c.process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}

	// Close log file
	if c.logFile != nil {
		c.logFile.Close()
	}

	c.process = nil
	fmt.Println("Juno node stopped successfully")
	return nil
}

// GetNodeStatus returns the status of the Juno node
func (c *JunoClient) GetNodeStatus(ctx context.Context) (string, error) {
	if c.process == nil {
		return "not running", nil
	}

	// Check if process is still running
	if err := c.process.Signal(os.Signal(nil)); err != nil {
		return "stopped", nil
	}

	// Try to get status via HTTP API
	status, err := c.getHTTPStatus()
	if err != nil {
		return "running (status check failed)", nil
	}

	return status, nil
}

// getHTTPStatus tries to get status via HTTP API
func (c *JunoClient) getHTTPStatus() (string, error) {
	// This is a placeholder - in a real implementation, you would make an HTTP request
	// to the Juno API to get the actual status
	return "running", nil
}

// Close cleans up resources
func (c *JunoClient) Close() error {
	if c.logFile != nil {
		return c.logFile.Close()
	}
	return nil
}
