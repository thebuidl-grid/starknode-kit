package clients

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/pflag"
)

// Mock configuration for testing
func createTestRethConfig() *RethConfig {
	tmpDir := "/tmp/test-reth"
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	return &RethConfig{
		InstallDir:        tmpDir,
		ExecutionType:     "full",
		ExecutionPeerPort: 30303,
		JWTPath:           filepath.Join(tmpDir, "ethereum_clients", "jwt", "jwt.hex"),
		LogFilePath:       filepath.Join(tmpDir, "ethereum_clients", "reth", "logs", fmt.Sprintf("reth_%s.log", timestamp)),
	}
}

func TestParseRethFlags(t *testing.T) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected *RethConfig
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			expected: &RethConfig{
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
			},
		},
		{
			name: "custom directory",
			args: []string{"cmd", "--directory", "/custom/path"},
			expected: &RethConfig{
				InstallDir:        "/custom/path",
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
			},
		},
		{
			name: "archive execution type",
			args: []string{"cmd", "--executiontype", "archive"},
			expected: &RethConfig{
				ExecutionType:     "archive",
				ExecutionPeerPort: 30303,
			},
		},
		{
			name: "custom peer port",
			args: []string{"cmd", "--executionpeerport", "30304"},
			expected: &RethConfig{
				ExecutionType:     "full",
				ExecutionPeerPort: 30304,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test args
			os.Args = tt.args

			// Create a new flag set for this test to avoid conflicts
			flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
			config := ParseRethFlagsWithFlagSet(flagSet)

			// Handle case where UserHomeDir fails
			if config == nil {
				t.Skip("Skipping test - UserHomeDir failed")
				return
			}

			if tt.expected.InstallDir != "" && config.InstallDir != tt.expected.InstallDir {
				t.Errorf("InstallDir = %v, want %v", config.InstallDir, tt.expected.InstallDir)
			}

			if config.ExecutionType != tt.expected.ExecutionType {
				t.Errorf("ExecutionType = %v, want %v", config.ExecutionType, tt.expected.ExecutionType)
			}

			if config.ExecutionPeerPort != tt.expected.ExecutionPeerPort {
				t.Errorf("ExecutionPeerPort = %v, want %v", config.ExecutionPeerPort, tt.expected.ExecutionPeerPort)
			}

			// Verify derived paths are set
			if config.JWTPath == "" {
				t.Error("JWTPath should not be empty")
			}

			if config.LogFilePath == "" {
				t.Error("LogFilePath should not be empty")
			}

			// Verify JWT path format
			expectedJWTPath := filepath.Join(config.InstallDir, "ethereum_clients", "jwt", "jwt.hex")
			if config.JWTPath != expectedJWTPath {
				t.Errorf("JWTPath = %v, want %v", config.JWTPath, expectedJWTPath)
			}

			// Verify log path contains timestamp
			if !strings.Contains(config.LogFilePath, "reth_") || !strings.Contains(config.LogFilePath, ".log") {
				t.Errorf("LogFilePath should contain timestamp: %v", config.LogFilePath)
			}
		})
	}
}

func TestGetRethCommand(t *testing.T) {
	tests := []struct {
		name       string
		installDir string
		goos       string
		expected   string
	}{
		{
			name:       "linux command",
			installDir: "/home/user",
			goos:       "linux",
			expected:   "/home/user/ethereum_clients/reth/reth",
		},
		{
			name:       "darwin command",
			installDir: "/Users/user",
			goos:       "darwin",
			expected:   "/Users/user/ethereum_clients/reth/reth",
		},
		{
			name:       "windows command",
			installDir: "C:\\Users\\user",
			goos:       "windows",
			expected:   "C:\\Users\\user\\ethereum_clients\\reth\\reth.exe",
		},
	}

	// We can't change runtime.GOOS in tests, so we'll test the current platform only
	defer func() {
		// No need to restore anything since we can't modify runtime.GOOS
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can only test the current platform's logic
			if runtime.GOOS == tt.goos {
				result := GetRethCommand(tt.installDir)
				if result != tt.expected {
					t.Errorf("GetRethCommand() = %v, want %v", result, tt.expected)
				}
			}
		})
	}

	// Test current platform
	t.Run("current platform", func(t *testing.T) {
		installDir := "/test/path"
		result := GetRethCommand(installDir)

		if runtime.GOOS == "windows" {
			expectedSuffix := "reth.exe"
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("On Windows, command should end with %s, got %s", expectedSuffix, result)
			}
		} else {
			expectedSuffix := "reth"
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("On non-Windows, command should end with %s, got %s", expectedSuffix, result)
			}
		}

		if !strings.Contains(result, installDir) {
			t.Errorf("Command should contain install directory %s, got %s", installDir, result)
		}
	})
}

func TestBuildRethArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   *RethConfig
		expected []string
	}{
		{
			name: "full sync mode",
			config: &RethConfig{
				InstallDir:        "/test",
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
				JWTPath:           "/test/jwt.hex",
			},
			expected: []string{
				"node",
				"--network", "mainnet",
				"--http",
				"--http.addr", "0.0.0.0",
				"--http.port", "8545",
				"--http.api", "eth,net,engine,admin",
				"--http.corsdomain", "*",
				"--authrpc.addr", "0.0.0.0",
				"--authrpc.port", "8551",
				"--authrpc.jwtsecret", "/test/jwt.hex",
				"--port", "30303",
				"--metrics", "0.0.0.0:6060",
				"--datadir", "/test/ethereum_clients/reth/database",
			},
		},
		{
			name: "archive mode",
			config: &RethConfig{
				InstallDir:        "/test",
				ExecutionType:     "archive",
				ExecutionPeerPort: 30304,
				JWTPath:           "/test/jwt.hex",
			},
			expected: []string{
				"node",
				"--network", "mainnet",
				"--http",
				"--http.addr", "0.0.0.0",
				"--http.port", "8545",
				"--http.api", "eth,net,engine,admin",
				"--http.corsdomain", "*",
				"--authrpc.addr", "0.0.0.0",
				"--authrpc.port", "8551",
				"--authrpc.jwtsecret", "/test/jwt.hex",
				"--port", "30304",
				"--metrics", "0.0.0.0:6060",
				"--archive",
				"--datadir", "/test/ethereum_clients/reth/database",
			},
		},
		{
			name: "custom peer port",
			config: &RethConfig{
				InstallDir:        "/custom",
				ExecutionType:     "full",
				ExecutionPeerPort: 31313,
				JWTPath:           "/custom/jwt.hex",
			},
			expected: []string{
				"node",
				"--network", "mainnet",
				"--http",
				"--http.addr", "0.0.0.0",
				"--http.port", "8545",
				"--http.api", "eth,net,engine,admin",
				"--http.corsdomain", "*",
				"--authrpc.addr", "0.0.0.0",
				"--authrpc.port", "8551",
				"--authrpc.jwtsecret", "/custom/jwt.hex",
				"--port", "31313",
				"--metrics", "0.0.0.0:6060",
				"--datadir", "/custom/ethereum_clients/reth/database",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildRethArgs(tt.config)

			if len(result) != len(tt.expected) {
				t.Errorf("BuildRethArgs() returned %d args, want %d", len(result), len(tt.expected))
				t.Errorf("Got: %v", result)
				t.Errorf("Expected: %v", tt.expected)
				return
			}

			for i, arg := range result {
				if i < len(tt.expected) && arg != tt.expected[i] {
					t.Errorf("BuildRethArgs() arg[%d] = %v, want %v", i, arg, tt.expected[i])
				}
			}
		})
	}
}

func TestRethConfig_Validation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := createTestRethConfig()

		if config.InstallDir == "" {
			t.Error("InstallDir should not be empty")
		}

		if config.ExecutionType == "" {
			t.Error("ExecutionType should not be empty")
		}

		if config.ExecutionPeerPort <= 0 {
			t.Error("ExecutionPeerPort should be positive")
		}

		if config.JWTPath == "" {
			t.Error("JWTPath should not be empty")
		}

		if config.LogFilePath == "" {
			t.Error("LogFilePath should not be empty")
		}
	})

	t.Run("execution types", func(t *testing.T) {
		tests := []struct {
			execType   string
			hasArchive bool
		}{
			{"full", false},
			{"archive", true},
		}

		for _, test := range tests {
			config := createTestRethConfig()
			config.ExecutionType = test.execType

			args := BuildRethArgs(config)
			argsStr := strings.Join(args, " ")

			hasArchiveFlag := strings.Contains(argsStr, "--archive")
			if hasArchiveFlag != test.hasArchive {
				t.Errorf("ExecutionType %s: expected archive flag = %v, got %v in args: %s",
					test.execType, test.hasArchive, hasArchiveFlag, argsStr)
			}
		}
	})
}

func TestStartReth_ErrorHandling(t *testing.T) {
	t.Run("missing binary", func(t *testing.T) {
		config := createTestRethConfig()
		config.InstallDir = "/nonexistent/path"

		err := StartReth(config)
		if err == nil {
			t.Error("StartReth should return error for missing binary")
		}

		if !strings.Contains(err.Error(), "reth binary not found") {
			t.Errorf("Error should mention missing binary, got: %v", err)
		}
	})

	t.Run("invalid log directory", func(t *testing.T) {
		config := createTestRethConfig()
		// Set log path to a file (not directory) to cause mkdir error
		config.LogFilePath = "/dev/null/invalid.log"

		err := StartReth(config)
		if err == nil {
			t.Error("StartReth should return error for invalid log directory")
		}

		if !strings.Contains(err.Error(), "error creating log directory") {
			t.Errorf("Error should mention log directory creation, got: %v", err)
		}
	})
}

func TestRethArgsComparison(t *testing.T) {
	t.Run("compare full vs archive args", func(t *testing.T) {
		configFull := createTestRethConfig()
		configFull.ExecutionType = "full"

		configArchive := createTestRethConfig()
		configArchive.ExecutionType = "archive"

		argsFull := BuildRethArgs(configFull)
		argsArchive := BuildRethArgs(configArchive)

		// Archive should have one extra argument (--archive)
		if len(argsArchive) != len(argsFull)+1 {
			t.Errorf("Archive mode should have exactly one more argument than full mode. Full: %d, Archive: %d",
				len(argsFull), len(argsArchive))
		}

		// Archive args should contain --archive flag
		argsArchiveStr := strings.Join(argsArchive, " ")
		if !strings.Contains(argsArchiveStr, "--archive") {
			t.Error("Archive mode should contain --archive flag")
		}

		// Full args should not contain --archive flag
		argsFullStr := strings.Join(argsFull, " ")
		if strings.Contains(argsFullStr, "--archive") {
			t.Error("Full mode should not contain --archive flag")
		}
	})
}

func TestRethCommandStructure(t *testing.T) {
	t.Run("verify command structure", func(t *testing.T) {
		config := createTestRethConfig()
		args := BuildRethArgs(config)

		// First argument should always be "node"
		if len(args) == 0 || args[0] != "node" {
			t.Error("First argument should be 'node'")
		}

		// Should contain network specification
		argsStr := strings.Join(args, " ")
		if !strings.Contains(argsStr, "--network mainnet") {
			t.Error("Should contain --network mainnet")
		}

		// Should contain HTTP API configuration
		requiredHTTPArgs := []string{
			"--http",
			"--http.addr 0.0.0.0",
			"--http.port 8545",
			"--http.api eth,net,engine,admin",
		}

		for _, required := range requiredHTTPArgs {
			if !strings.Contains(argsStr, required) {
				t.Errorf("Should contain %s", required)
			}
		}

		// Should contain AuthRPC configuration
		requiredAuthArgs := []string{
			"--authrpc.addr 0.0.0.0",
			"--authrpc.port 8551",
			"--authrpc.jwtsecret",
		}

		for _, required := range requiredAuthArgs {
			if !strings.Contains(argsStr, required) {
				t.Errorf("Should contain %s", required)
			}
		}

		// Should contain metrics
		if !strings.Contains(argsStr, "--metrics 0.0.0.0:6060") {
			t.Error("Should contain metrics configuration")
		}

		// Should contain data directory
		if !strings.Contains(argsStr, "--datadir") {
			t.Error("Should contain --datadir")
		}
	})
}

// Benchmark tests
func BenchmarkBuildRethArgs(b *testing.B) {
	config := createTestRethConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildRethArgs(config)
	}
}

func BenchmarkGetRethCommand(b *testing.B) {
	installDir := "/test/path"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetRethCommand(installDir)
	}
}

func BenchmarkParseRethFlags(b *testing.B) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"cmd", "--directory", "/test", "--executiontype", "full"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseRethFlags()
	}
}
