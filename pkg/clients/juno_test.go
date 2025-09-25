package clients

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
)

func TestDefaultJunoConfig(t *testing.T) {
	config := DefaultJunoConfig()

	// Test default values
	if config.Network != "mainnet" {
		t.Errorf("Expected Network to be 'mainnet', got '%s'", config.Network)
	}

	if config.Port != "6060" {
		t.Errorf("Expected Port to be '6060', got '%s'", config.Port)
	}

	if !config.UseSnapshot {
		t.Error("Expected UseSnapshot to be true")
	}

	if config.DataDir != "./juno-data" {
		t.Errorf("Expected DataDir to be './juno-data', got '%s'", config.DataDir)
	}

	if config.EthNode != "ws://localhost:8546" {
		t.Errorf("Expected EthNode to be 'ws://localhost:8546', got '%s'", config.EthNode)
	}

	// Test environment variables
	expectedEnv := []string{
		"JUNO_NETWORK=mainnet",
		"JUNO_HTTP_PORT=6060",
		"JUNO_HTTP_HOST=0.0.0.0",
	}

	if len(config.Environment) != len(expectedEnv) {
		t.Errorf("Expected %d environment variables, got %d", len(expectedEnv), len(config.Environment))
	}

	for i, expected := range expectedEnv {
		if config.Environment[i] != expected {
			t.Errorf("Expected environment variable %d to be '%s', got '%s'", i, expected, config.Environment[i])
		}
	}
}

func TestNewJunoClient(t *testing.T) {
	// Mock the presence of the Juno binary
	tempDir := t.TempDir()
	junoDir := filepath.Join(tempDir, "juno")
	if err := os.MkdirAll(junoDir, 0755); err != nil {
		t.Fatalf("Failed to create juno dir: %v", err)
	}
	junoPath := filepath.Join(junoDir, "juno")
	if err := os.WriteFile(junoPath, []byte("#!/bin/sh\necho 'juno version 0.14.6'\n"), 0755); err != nil {
		t.Fatalf("Failed to create dummy juno binary: %v", err)
	}
	// Save and override InstallClientsDir
	origInstallClientsDir := constants.InstallClientsDir
	constants.InstallClientsDir = tempDir
	defer func() { constants.InstallClientsDir = origInstallClientsDir }()

	// Test with nil config
	client, err := NewJunoClient(nil)
	if err != nil {
		t.Errorf("Expected no error when config is nil, got: %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created when config is nil")
	}

	// Test with custom config
	customConfig := &JunoConfig{
		Network:     "sepolia",
		Port:        "6061",
		UseSnapshot: false,
		DataDir:     "/custom/path",
		EthNode:     "ws://custom:8546",
		Environment: []string{"CUSTOM_VAR=value"},
	}

	client, err = NewJunoClient(customConfig)
	if err != nil {
		t.Errorf("Expected no error with custom config, got: %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created with custom config")
	}

	if client.config.Network != "sepolia" {
		t.Errorf("Expected Network to be 'sepolia', got '%s'", client.config.Network)
	}

	if client.config.Port != "6061" {
		t.Errorf("Expected Port to be '6061', got '%s'", client.config.Port)
	}

	if client.config.UseSnapshot {
		t.Error("Expected UseSnapshot to be false")
	}

	if client.config.DataDir != "/custom/path" {
		t.Errorf("Expected DataDir to be '/custom/path', got '%s'", client.config.DataDir)
	}

	if client.config.EthNode != "ws://custom:8546" {
		t.Errorf("Expected EthNode to be 'ws://custom:8546', got '%s'", client.config.EthNode)
	}
}

func TestGetJunoPath(t *testing.T) {
	// Test when Juno is not installed
	path := getJunoPath()
	if path != "" {
		t.Errorf("Expected empty path when Juno is not installed, got '%s'", path)
	}
}

func TestBuildJunoArgs(t *testing.T) {
	config := &JunoConfig{
		Network:     "mainnet",
		Port:        "6060",
		UseSnapshot: true,
		DataDir:     "/test/data",
		EthNode:     "ws://localhost:8546",
	}

	client := &JunoClient{config: config}
	args := client.buildJunoArgs()

	// Test required arguments
	expectedArgs := []string{
		"--http",
		"--http-port=6060",
		"--http-host=0.0.0.0",
		"--db-path=/test/data",
		"--eth-node=ws://localhost:8546",
		"--network=mainnet",
		"--snapshot",
		"--metrics",
		"--metrics-port=6060",
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}
}

func TestBuildJunoArgsSepolia(t *testing.T) {
	config := &JunoConfig{
		Network:     "sepolia",
		Port:        "6061",
		UseSnapshot: false,
		DataDir:     "/test/sepolia",
		EthNode:     "ws://sepolia:8546",
	}

	client := &JunoClient{config: config}
	args := client.buildJunoArgs()

	// Test sepolia network arguments
	expectedArgs := []string{
		"--http",
		"--http-port=6061",
		"--http-host=0.0.0.0",
		"--db-path=/test/sepolia",
		"--eth-node=ws://sepolia:8546",
		"--network=sepolia",
		"--metrics",
		"--metrics-port=6060",
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}
}

func TestBuildJunoArgsSepoliaIntegration(t *testing.T) {
	config := &JunoConfig{
		Network:     "sepolia-integration",
		Port:        "6062",
		UseSnapshot: true,
		DataDir:     "/test/sepolia-integration",
		EthNode:     "ws://sepolia-integration:8546",
	}

	client := &JunoClient{config: config}
	args := client.buildJunoArgs()

	// Test sepolia-integration network arguments
	expectedArgs := []string{
		"--http",
		"--http-port=6062",
		"--http-host=0.0.0.0",
		"--db-path=/test/sepolia-integration",
		"--eth-node=ws://sepolia-integration:8546",
		"--network=sepolia-integration",
		"--snapshot",
		"--metrics",
		"--metrics-port=6060",
	}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d arguments, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected argument %d to be '%s', got '%s'", i, expected, args[i])
		}
	}
}

func TestJunoClientStartNode(t *testing.T) {
	// Create temporary directory for test
	tempDir := t.TempDir()

	config := &JunoConfig{
		Network:     "mainnet",
		Port:        "6060",
		UseSnapshot: true,
		DataDir:     tempDir,
		EthNode:     "ws://localhost:8546",
	}

	client := &JunoClient{config: config}

	// Mock the junoPath to avoid actual binary execution
	client.junoPath = "echo" // Use echo as a mock command

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// The dummy binary (echo) will start successfully, so we do not expect an error
	err := client.StartNode(ctx)
	if err != nil {
		t.Errorf("Expected no error when starting with dummy binary, got: %v", err)
	}

	// Test that directories were created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Expected data directory to be created")
	}

	logsDir := filepath.Join(filepath.Dir(tempDir), "logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		t.Error("Expected logs directory to be created")
	}
}

func TestJunoClientStopNode(t *testing.T) {
	client := &JunoClient{}

	// Test stopping when no process is running
	_ = client.StopNode(context.Background()) // Should not panic

	// Test with mock process
	client.process = &os.Process{Pid: 99999}  // Non-existent PID
	_ = client.StopNode(context.Background()) // Should not panic
}

func TestJunoClientGetNodeStatus(t *testing.T) {
	client := &JunoClient{}

	// Test when no process is running
	status, err := client.GetNodeStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error when getting status of non-existent process, got: %v", err)
	}
	if status != "not running" {
		t.Errorf("Expected status 'not running', got '%s'", status)
	}

	// Test with mock process
	client.process = &os.Process{Pid: 99999} // Non-existent PID
	status, err = client.GetNodeStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error when getting status of invalid process, got: %v", err)
	}
	if status != "stopped" {
		t.Errorf("Expected status 'stopped', got '%s'", status)
	}
}

func TestJunoClientClose(t *testing.T) {
	client := &JunoClient{}

	// Test closing when no log file is open
	err := client.Close()
	if err != nil {
		t.Errorf("Expected no error when closing without log file, got: %v", err)
	}

	// Test closing with log file
	tempFile, err := os.CreateTemp("", "juno_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	client.logFile = tempFile
	err = client.Close()
	if err != nil {
		t.Errorf("Expected no error when closing with log file, got: %v", err)
	}
}

func TestJunoConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *JunoConfig
		wantErr bool
	}{
		{
			name: "valid mainnet config",
			config: &JunoConfig{
				Network:     "mainnet",
				Port:        "6060",
				UseSnapshot: true,
				DataDir:     "/test/data",
				EthNode:     "ws://localhost:8546",
			},
			wantErr: false,
		},
		{
			name: "valid sepolia config",
			config: &JunoConfig{
				Network:     "sepolia",
				Port:        "6061",
				UseSnapshot: false,
				DataDir:     "/test/sepolia",
				EthNode:     "ws://sepolia:8546",
			},
			wantErr: false,
		},
		{
			name: "valid sepolia-integration config",
			config: &JunoConfig{
				Network:     "sepolia-integration",
				Port:        "6062",
				UseSnapshot: true,
				DataDir:     "/test/sepolia-integration",
				EthNode:     "ws://sepolia-integration:8546",
			},
			wantErr: false,
		},
		{
			name: "invalid network",
			config: &JunoConfig{
				Network:     "invalid",
				Port:        "6060",
				UseSnapshot: true,
				DataDir:     "/test/data",
				EthNode:     "ws://localhost:8546",
			},
			wantErr: false, // We don't validate network in buildJunoArgs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &JunoClient{config: tt.config}
			args := client.buildJunoArgs()

			// Check that required arguments are present
			requiredArgs := []string{
				"--http",
				"--http-host=0.0.0.0",
				"--metrics",
				"--metrics-port=6060",
			}

			for _, required := range requiredArgs {
				found := false
				for _, arg := range args {
					if arg == required {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Required argument '%s' not found in args", required)
				}
			}

			// Check that eth-node is present (substring match)
			ethNodeFound := false
			for _, arg := range args {
				if strings.Contains(arg, "--eth-node=") {
					ethNodeFound = true
					break
				}
			}
			if !ethNodeFound {
				t.Error("--eth-node argument not found in args")
			}
		})
	}
}

func TestJunoClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directory for test
	tempDir := t.TempDir()

	config := &JunoConfig{
		Network:     "mainnet",
		Port:        "6060",
		UseSnapshot: true,
		DataDir:     tempDir,
		EthNode:     "ws://localhost:8546",
	}

	client := &JunoClient{config: config}

	// Test that client can be created
	if client == nil {
		t.Fatal("Failed to create Juno client")
	}

	// Test that config is properly set
	if client.config.Network != "mainnet" {
		t.Errorf("Expected Network to be 'mainnet', got '%s'", client.config.Network)
	}

	// Test that arguments can be built
	args := client.buildJunoArgs()
	if len(args) == 0 {
		t.Error("Expected non-empty arguments list")
	}

	// Test that required arguments are present
	hasHttp := false
	hasEthNode := false
	for _, arg := range args {
		if arg == "--http" {
			hasHttp = true
		}
		if strings.Contains(arg, "--eth-node=") {
			hasEthNode = true
		}
	}

	if !hasHttp {
		t.Error("Expected --http argument to be present")
	}
	if !hasEthNode {
		t.Error("Expected --eth-node argument to be present")
	}
}

// Benchmark tests
func BenchmarkBuildJunoArgs(b *testing.B) {
	config := &JunoConfig{
		Network:     "mainnet",
		Port:        "6060",
		UseSnapshot: true,
		DataDir:     "/test/data",
		EthNode:     "ws://localhost:8546",
	}

	client := &JunoClient{config: config}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.buildJunoArgs()
	}
}

func BenchmarkNewJunoClient(b *testing.B) {
	config := &JunoConfig{
		Network:     "mainnet",
		Port:        "6060",
		UseSnapshot: true,
		DataDir:     "/test/data",
		EthNode:     "ws://localhost:8546",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewJunoClient(config)
		if err != nil {
			b.Fatalf("Failed to create client: %v", err)
		}
	}
}
