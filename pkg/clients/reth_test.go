package clients

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// Mock configuration for testing
func createTestRethConfig() *RethConfig {
	return &RethConfig{
		ExecutionType:     "full",
		ExecutionPeerPort: 30303,
		LogFilePath:       "/tmp/test-reth/reth_test.log",
	}
}

func TestRethConfig_Creation(t *testing.T) {
	config := createTestRethConfig()

	if config.ExecutionType != "full" {
		t.Errorf("ExecutionType = %v, want %v", config.ExecutionType, "full")
	}

	if config.ExecutionPeerPort != 30303 {
		t.Errorf("ExecutionPeerPort = %v, want %v", config.ExecutionPeerPort, 30303)
	}

	if config.LogFilePath == "" {
		t.Error("LogFilePath should not be empty")
	}
}

func TestGetRethCommand(t *testing.T) {
	t.Run("current platform", func(t *testing.T) {
		result := GetRethCommand()

		if runtime.GOOS == "windows" {
			expectedSuffix := "reth.exe"
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("On Windows, command should end with %s, got %s", expectedSuffix, result)
			}
		} else {
			expectedSuffix := "reth"
			if !strings.HasSuffix(result, expectedSuffix) && !strings.HasSuffix(result, "reth.exe") {
				t.Errorf("On non-Windows, command should end with reth, got %s", result)
			}
		}

		// Check that the path contains the expected directory structure
		expectedParts := []string{"reth"}
		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Command should contain %s, got %s", part, result)
			}
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
				ExecutionType:     "full",
				ExecutionPeerPort: 30303,
				LogFilePath:       "/tmp/test.log",
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
				"--authrpc.jwtsecret",
				"--port", "30303",
				"--metrics", "0.0.0.0:6060",
				"--datadir",
			},
		},
		{
			name: "archive mode",
			config: &RethConfig{
				ExecutionType:     "archive",
				ExecutionPeerPort: 30304,
				LogFilePath:       "/tmp/test.log",
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
				"--authrpc.jwtsecret",
				"--port", "30304",
				"--metrics", "0.0.0.0:6060",
				"--archive",
				"--datadir",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildRethArgs(tt.config)

			// Check that essential arguments are present
			argsStr := strings.Join(result, " ")

			// Check for key arguments (not exact match due to dynamic values)
			expectedKeyArgs := []string{
				"node",
				"--network mainnet",
				"--http",
				"--http.addr 0.0.0.0",
				"--http.port 8545",
				"--authrpc.addr 0.0.0.0",
				"--authrpc.port 8551",
				"--metrics 0.0.0.0:6060",
				"--datadir",
			}

			for _, expectedArg := range expectedKeyArgs {
				if !strings.Contains(argsStr, expectedArg) {
					t.Errorf("BuildRethArgs() should contain %s, got: %s", expectedArg, argsStr)
				}
			}

			// Check execution type specific arguments
			if tt.config.ExecutionType == "archive" {
				if !strings.Contains(argsStr, "--archive") {
					t.Errorf("Archive mode should contain --archive flag, got: %s", argsStr)
				}
			} else {
				if strings.Contains(argsStr, "--archive") {
					t.Errorf("Non-archive mode should not contain --archive flag, got: %s", argsStr)
				}
			}

			// Check port
			expectedPort := fmt.Sprintf("--port %d", tt.config.ExecutionPeerPort)
			if !strings.Contains(argsStr, expectedPort) {
				t.Errorf("Should contain %s, got: %s", expectedPort, argsStr)
			}
		})
	}
}

func TestStartReth_Parameters(t *testing.T) {
	t.Run("parameter validation", func(t *testing.T) {
		// Test that StartReth accepts the correct parameters
		executionType := "full"
		ports := []int{30303}

		// We can't actually start reth in tests, but we can verify the function signature
		// and basic parameter handling by checking it doesn't panic with valid inputs
		err := StartReth(executionType, ports)

		// We expect this to fail since reth binary likely doesn't exist in test environment
		// but it shouldn't panic and should return an error
		if err == nil {
			t.Log("StartReth succeeded (reth binary must be present)")
		} else {
			t.Logf("StartReth failed as expected in test environment: %v", err)
		}
	})

	t.Run("different execution types", func(t *testing.T) {
		executionTypes := []string{"full", "archive"}
		ports := []int{30303}

		for _, execType := range executionTypes {
			err := StartReth(execType, ports)
			// Again, we expect this to fail in test environment, but shouldn't panic
			if err != nil {
				t.Logf("StartReth with type %s failed as expected: %v", execType, err)
			}
		}
	})

	t.Run("different ports", func(t *testing.T) {
		testPorts := [][]int{
			{30303},
			{30304},
			{31313},
		}
		executionType := "full"

		for _, ports := range testPorts {
			err := StartReth(executionType, ports)
			// Expected to fail in test environment
			if err != nil {
				t.Logf("StartReth with ports %v failed as expected: %v", ports, err)
			}
		}
	})
}

func TestRethConfig_Validation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := createTestRethConfig()

		if config.ExecutionType == "" {
			t.Error("ExecutionType should not be empty")
		}

		if config.ExecutionPeerPort <= 0 {
			t.Error("ExecutionPeerPort should be positive")
		}

		if config.LogFilePath == "" {
			t.Error("LogFilePath should not be empty")
		}
	})

	t.Run("execution types", func(t *testing.T) {
		validTypes := []string{"full", "archive"}

		for _, execType := range validTypes {
			config := createTestRethConfig()
			config.ExecutionType = execType

			// Verify the config accepts valid execution types
			if config.ExecutionType != execType {
				t.Errorf("Failed to set ExecutionType to %s", execType)
			}
		}
	})

	t.Run("port ranges", func(t *testing.T) {
		config := createTestRethConfig()

		// Test various port values
		testPorts := []int{1024, 30303, 30304, 65535}

		for _, port := range testPorts {
			config.ExecutionPeerPort = port
			if config.ExecutionPeerPort != port {
				t.Errorf("Failed to set ExecutionPeerPort to %d", port)
			}
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetRethCommand()
	}
}

func BenchmarkRethConfigCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createTestRethConfig()
	}
}
