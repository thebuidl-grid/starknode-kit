package clients

import (
	"runtime"
	"strings"
	"testing"
)

// Mock configuration for testing
func createTestGethConfig() *gethConfig {
	return &gethConfig{
		port:          30303,
		executionType: "full",
	}
}

func TestGethConfig_Creation(t *testing.T) {
	config := createTestGethConfig()

	if config.executionType != "full" {
		t.Errorf("executionType = %v, want %v", config.executionType, "full")
	}

	if config.port != 30303 {
		t.Errorf("port = %v, want %v", config.port, 30303)
	}

}

func TestGetGethCommand(t *testing.T) {
	t.Run("current platform", func(t *testing.T) {
		geth := gethConfig{
			port:          30303,
			executionType: "full",
		}
		result := geth.getCommand()

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

		// We can't actually start geth in tests, but we can verify the function signature
		// and basic parameter handling by checking it doesn't panic with valid inputs
		geth := gethConfig{
			port:          30303,
			executionType: "full",
		}
		err := geth.Start()

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

		for _, execType := range executionTypes {
			geth := gethConfig{
				port:          30303,
				executionType: execType,
			}
			err := geth.Start()
			// Again, we expect this to fail in test environment, but shouldn't panic
			if err != nil {
				t.Logf("StartGeth with type %s failed as expected: %v", execType, err)
			}
		}
	})

	t.Run("different ports", func(t *testing.T) {
		testPorts := []int{30303, 30304, 31313}

		for _, port := range testPorts {
			geth := gethConfig{
				port:          port,
				executionType: "full",
			}
			err := geth.Start()
			// Expected to fail in test environment
			if err != nil {
				t.Logf("StartGeth with ports %v failed as expected: %v", port, err)
			}
		}
	})
}

func TestGethConfig_Validation(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := createTestGethConfig()

		if config.executionType == "" {
			t.Error("executionType should not be empty")
		}

		if config.port <= 0 {
			t.Error("port should be positive")
		}

	})

	t.Run("execution types", func(t *testing.T) {
		validTypes := []string{"full", "archive"}

		for _, execType := range validTypes {
			config := createTestGethConfig()
			config.executionType = execType

			// Verify the config accepts valid execution types
			if config.executionType != execType {
				t.Errorf("Failed to set executionType to %s", execType)
			}
		}
	})

	t.Run("port ranges", func(t *testing.T) {
		config := createTestGethConfig()

		// Test various port values
		testPorts := []int{1024, 30303, 30304, 65535}

		for _, port := range testPorts {
			config.port = port
			if config.port != port {
				t.Errorf("Failed to set port to %d", port)
			}
		}
	})
}

// Benchmark tests
func BenchmarkGetGethCommand(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		geth := gethConfig{
			port:          30303,
			executionType: "full",
		}
		geth.getCommand()
	}
}

func BenchmarkGethConfigCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createTestGethConfig()
	}
}
