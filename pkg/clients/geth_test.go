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
func createTestGethConfig() *GethConfig {
	tmpDir := "/tmp/test-geth"
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	return &GethConfig{
		InstallDir:        tmpDir,
		ExecutionType:     "full",
		ExecutionPeerPort: 30303,
		JWTPath:           filepath.Join(tmpDir, "ethereum_clients", "jwt", "jwt.hex"),
		LogFilePath:       filepath.Join(tmpDir, "ethereum_clients", "geth", "logs", fmt.Sprintf("geth_%s.log", timestamp)),
	}
}

func TestParseGethFlags(t *testing.T) {
	// Save original args
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	tests := []struct {
		name     string
		args     []string
		expected *GethConfig
	}{
		{
			name: "default values",
			args: []string{"cmd"},
			expected: &GethConfig{
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
			},
		},
		{
			name: "custom directory",
			args: []string{"cmd", "--directory", "/custom/path"},
			expected: &GethConfig{
				InstallDir:        "/custom/path",
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
			},
		},
		{
			name: "archive execution type",
			args: []string{"cmd", "--executiontype", "archive"},
			expected: &GethConfig{
				ExecutionType:     "archive",
				ExecutionPeerPort: 30303,
			},
		},
		{
			name: "custom peer port",
			args: []string{"cmd", "--executionpeerport", "30304"},
			expected: &GethConfig{
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
			config := ParseGethFlagsWithFlagSet(flagSet)

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
			if !strings.Contains(config.LogFilePath, "geth_") || !strings.Contains(config.LogFilePath, ".log") {
				t.Errorf("LogFilePath should contain timestamp: %v", config.LogFilePath)
			}
		})
	}
}

func TestGetGethCommand(t *testing.T) {
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
			expected:   "/home/user/ethereum_clients/geth/geth",
		},
		{
			name:       "darwin command",
			installDir: "/Users/user",
			goos:       "darwin",
			expected:   "/Users/user/ethereum_clients/geth/geth",
		},
		{
			name:       "windows command",
			installDir: "C:\\Users\\user",
			goos:       "windows",
			expected:   "C:\\Users\\user\\ethereum_clients\\geth\\geth.exe",
		},
	}

	// Note: We can't change runtime.GOOS in tests, so we'll test the current platform

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can only test the current platform's logic
			if runtime.GOOS == tt.goos {
				result := GetGethCommand(tt.installDir)
				if result != tt.expected {
					t.Errorf("GetGethCommand() = %v, want %v", result, tt.expected)
				}
			}
		})
	}

	// Test current platform
	t.Run("current platform", func(t *testing.T) {
		installDir := "/test/path"
		result := GetGethCommand(installDir)

		if runtime.GOOS == "windows" {
			expectedSuffix := "geth.exe"
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("On Windows, command should end with %s, got %s", expectedSuffix, result)
			}
		} else {
			expectedSuffix := "geth"
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("On non-Windows, command should end with %s, got %s", expectedSuffix, result)
			}
		}

		if !strings.Contains(result, installDir) {
			t.Errorf("Command should contain install directory %s, got %s", installDir, result)
		}
	})
}

func TestBuildGethArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   *GethConfig
		expected []string
	}{
		{
			name: "full sync mode",
			config: &GethConfig{
				InstallDir:        "/test",
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
				JWTPath:           "/test/jwt.hex",
			},
			expected: []string{
				"--mainnet",
				"--port=30303",
				"--discovery.port=30303",
				"--http",
				"--http.api=eth,net,engine,admin",
				"--http.corsdomain=*",
				"--http.addr=0.0.0.0",
				"--http.port=8545",
				"--authrpc.jwtsecret=/test/jwt.hex",
				"--authrpc.addr=0.0.0.0",
				"--authrpc.port=8551",
				"--authrpc.vhosts=*",
				"--metrics",
				"--metrics.addr=0.0.0.0",
				"--metrics.port=6060",
				"--syncmode=snap",
				"--datadir=/test/ethereum_clients/geth/database",
			},
		},
		{
			name: "archive mode",
			config: &GethConfig{
				InstallDir:        "/test",
				ExecutionType:     "archive",
				ExecutionPeerPort: 30304,
				JWTPath:           "/test/jwt.hex",
			},
			expected: []string{
				"--mainnet",
				"--port=30304",
				"--discovery.port=30304",
				"--http",
				"--http.api=eth,net,engine,admin",
				"--http.corsdomain=*",
				"--http.addr=0.0.0.0",
				"--http.port=8545",
				"--authrpc.jwtsecret=/test/jwt.hex",
				"--authrpc.addr=0.0.0.0",
				"--authrpc.port=8551",
				"--authrpc.vhosts=*",
				"--metrics",
				"--metrics.addr=0.0.0.0",
				"--metrics.port=6060",
				"--syncmode=full",
				"--gcmode=archive",
				"--datadir=/test/ethereum_clients/geth/database",
			},
		},
		{
			name: "custom peer port",
			config: &GethConfig{
				InstallDir:        "/custom",
				ExecutionType:     "full",
				ExecutionPeerPort: 31313,
				JWTPath:           "/custom/jwt.hex",
			},
			expected: []string{
				"--mainnet",
				"--port=31313",
				"--discovery.port=31313",
				"--http",
				"--http.api=eth,net,engine,admin",
				"--http.corsdomain=*",
				"--http.addr=0.0.0.0",
				"--http.port=8545",
				"--authrpc.jwtsecret=/custom/jwt.hex",
				"--authrpc.addr=0.0.0.0",
				"--authrpc.port=8551",
				"--authrpc.vhosts=*",
				"--metrics",
				"--metrics.addr=0.0.0.0",
				"--metrics.port=6060",
				"--syncmode=snap",
				"--datadir=/custom/ethereum_clients/geth/database",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildGethArgs(tt.config)

			if len(result) != len(tt.expected) {
				t.Errorf("BuildGethArgs() returned %d args, want %d", len(result), len(tt.expected))
			}

			for i, arg := range result {
				if i < len(tt.expected) && arg != tt.expected[i] {
					t.Errorf("BuildGethArgs() arg[%d] = %v, want %v", i, arg, tt.expected[i])
				}
			}
		})
	}
}

func TestGethConfig_Validation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := createTestGethConfig()

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
		validTypes := []string{"full", "archive"}

		for _, execType := range validTypes {
			config := createTestGethConfig()
			config.ExecutionType = execType

			args := BuildGethArgs(config)
			argsStr := strings.Join(args, " ")

			switch execType {
			case "full":
				if !strings.Contains(argsStr, "--syncmode=snap") {
					t.Errorf("Full sync should contain --syncmode=snap, got %s", argsStr)
				}
			case "archive":
				if !strings.Contains(argsStr, "--syncmode=full") || !strings.Contains(argsStr, "--gcmode=archive") {
					t.Errorf("Archive sync should contain --syncmode=full and --gcmode=archive, got %s", argsStr)
				}
			}
		}
	})
}

func TestStartGeth_ErrorHandling(t *testing.T) {
	t.Run("missing binary", func(t *testing.T) {
		config := createTestGethConfig()
		config.InstallDir = "/nonexistent/path"

		err := StartGeth(config)
		if err == nil {
			t.Error("StartGeth should return error for missing binary")
		}

		if !strings.Contains(err.Error(), "geth binary not found") {
			t.Errorf("Error should mention missing binary, got: %v", err)
		}
	})

	t.Run("invalid log directory", func(t *testing.T) {
		config := createTestGethConfig()
		// Set log path to a file (not directory) to cause mkdir error
		config.LogFilePath = "/dev/null/invalid.log"

		err := StartGeth(config)
		if err == nil {
			t.Error("StartGeth should return error for invalid log directory")
		}

		if !strings.Contains(err.Error(), "error creating log directory") {
			t.Errorf("Error should mention log directory creation, got: %v", err)
		}
	})
}

// Benchmark tests
func BenchmarkBuildGethArgs(b *testing.B) {
	config := createTestGethConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BuildGethArgs(config)
	}
}

func BenchmarkGetGethCommand(b *testing.B) {
	installDir := "/test/path"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetGethCommand(installDir)
	}
}
