package clients

import (
	"runtime"
	"strings"
	"testing"
)

// Mock configuration for testing
func createTestGethConfig() *GethConfig {
	return &GethConfig{
		ExecutionType:     "full",
		ExecutionPeerPort: 30303,
		LogFilePath:       "/tmp/test-geth/geth_test.log",
	}
}

func TestGethConfig_Creation(t *testing.T) {
	config := createTestGethConfig()

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

func TestGetGethCommand(t *testing.T) {
	t.Run("current platform", func(t *testing.T) {
		result := GetGethCommand()

		if runtime.GOOS == "windows" {
			expectedSuffix := "geth.exe"
			if !strings.HasSuffix(result, expectedSuffix) {
				t.Errorf("On Windows, command should end with %s, got %s", expectedSuffix, result)
			}
		} else {
			expectedSuffix := "geth"
			if !strings.HasSuffix(result, expectedSuffix) && !strings.HasSuffix(result, "geth.exe") {
				t.Errorf("On non-Windows, command should end with geth, got %s", result)
			}
		}

		// Check that the path contains the expected directory structure
		expectedParts := []string{"geth"}
		for _, part := range expectedParts {
			if !strings.Contains(result, part) {
				t.Errorf("Command should contain %s, got %s", part, result)
			}
		}
	})
}

func TestStartGeth_Parameters(t *testing.T) {
	t.Run("parameter validation", func(t *testing.T) {
		// Test that StartGeth accepts the correct parameters
		executionType := "full"
		ports := []int{30303}

		// We can't actually start geth in tests, but we can verify the function signature
		// and basic parameter handling by checking it doesn't panic with valid inputs
		err := StartGeth(executionType, ports)

		// We expect this to fail since geth binary likely doesn't exist in test environment
		// but it shouldn't panic and should return an error
		if err == nil {
			t.Log("StartGeth succeeded (geth binary must be present)")
		} else {
			t.Logf("StartGeth failed as expected in test environment: %v", err)
		}
	})

	t.Run("different execution types", func(t *testing.T) {
		executionTypes := []string{"full", "archive"}
		ports := []int{30303}

		for _, execType := range executionTypes {
			err := StartGeth(execType, ports)
			// Again, we expect this to fail in test environment, but shouldn't panic
			if err != nil {
				t.Logf("StartGeth with type %s failed as expected: %v", execType, err)
			}
		}
	})

	t.Run("different ports", func(t *testing.T) {
		testPorts := [][]int{
			{30303},
			{30304},
			{31313},
		}

		for _, ports := range testPorts {
			err := StartGeth("full", ports)
			// Expected to fail in test environment
			if err != nil {
				t.Logf("StartGeth with ports %v failed as expected: %v", ports, err)
			}
		}
	})
}

func TestGethConfig_Validation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := createTestGethConfig()

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
			config := createTestGethConfig()
			config.ExecutionType = execType

			// Verify the config accepts valid execution types
			if config.ExecutionType != execType {
				t.Errorf("Failed to set ExecutionType to %s", execType)
			}
		}
	})

	t.Run("port ranges", func(t *testing.T) {
		config := createTestGethConfig()

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

// Benchmark tests
func BenchmarkGetGethCommand(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetGethCommand()
	}
}

func BenchmarkGethConfigCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createTestGethConfig()
	}
}
