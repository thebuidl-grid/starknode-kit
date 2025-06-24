package clients

import (
	"runtime"
	"strings"
	"testing"
)

// Mock configuration for testing
func createTestRethConfig() *rethConfig {
	return &rethConfig{
		executionType: "full",
		port:          30303,
	}
}

func TestRethConfig_Creation(t *testing.T) {
	config := createTestRethConfig()

	if config.executionType != "full" {
		t.Errorf("executionType = %v, want %v", config.executionType, "full")
	}

	if config.port != 30303 {
		t.Errorf("port = %v, want %v", config.port, 30303)
	}
}

func TestGetRethCommand(t *testing.T) {
	t.Run("current platform", func(t *testing.T) {
		config := createTestRethConfig()
		result := config.getCommand()

		if runtime.GOOS == "windows" {
			if !strings.HasSuffix(result, "reth.exe") {
				t.Errorf("On Windows, command should end with reth.exe, got %s", result)
			}
		} else {
			if !strings.HasSuffix(result, "reth") {
				t.Errorf("On non-Windows, command should end with reth, got %s", result)
			}
		}

		if !strings.Contains(result, "reth") {
			t.Errorf("Command should contain 'reth', got %s", result)
		}
	})
}

func TestBuildRethArgs(t *testing.T) {
	tests := []struct {
		name     string
		config   *rethConfig
		expected string // partial string match for simplicity
	}{
		{
			name: "full sync mode",
			config: &rethConfig{
				executionType: "full",
				port:          30303,
			},
			expected: "--port 30303",
		},
		{
			name: "archive mode",
			config: &rethConfig{
				executionType: "archive",
				port:          30304,
			},
			expected: "--archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.config.buildArgs()
			argsStr := strings.Join(args, " ")

			if !strings.Contains(argsStr, tt.expected) {
				t.Errorf("Expected args to contain %s, got %s", tt.expected, argsStr)
			}

			if tt.config.executionType == "archive" && !strings.Contains(argsStr, "--archive") {
				t.Error("Expected --archive flag for archive mode")
			}
		})
	}
}

func TestStartReth_DoesNotPanic(t *testing.T) {
	t.Run("start simulation", func(t *testing.T) {
		config := createTestRethConfig()
		err := config.Start()
		if err != nil {
			t.Logf("Start failed as expected (e.g., binary/log path may be missing): %v", err)
		}
	})
}

func BenchmarkBuildRethArgs(b *testing.B) {
	config := createTestRethConfig()
	for i := 0; i < b.N; i++ {
		config.buildArgs()
	}
}

func BenchmarkGetRethCommand(b *testing.B) {
	config := createTestRethConfig()
	for i := 0; i < b.N; i++ {
		config.getCommand()
	}
}
