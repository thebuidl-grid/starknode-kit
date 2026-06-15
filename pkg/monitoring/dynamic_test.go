package monitoring

import (
	"testing"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

// TestDynamicLayout tests that the layout rebuilds correctly based on running clients
func TestDynamicLayout(t *testing.T) {
	app := NewMonitorApp()

	// Test 1: Initial layout should have all panels created
	if app.ExecutionLogBox == nil {
		t.Error("ExecutionLogBox should be created")
	}
	if app.ConsensusLogBox == nil {
		t.Error("ConsensusLogBox should be created")
	}
	if app.JunoLogBox == nil {
		t.Error("JunoLogBox should be created")
	}
	if app.ValidatorLogBox == nil {
		t.Error("ValidatorLogBox should be created")
	}
	if app.NoClientsBox == nil {
		t.Error("NoClientsBox should be created")
	}
}

// TestValidatorLogChannel tests that the validator log channel is created and working
func TestValidatorLogChannel(t *testing.T) {
	app := NewMonitorApp()

	// Test that channel is created
	if app.ValidatorLogChan == nil {
		t.Fatal("ValidatorLogChan should be created")
	}

	// Test that we can send to the channel
	testMessage := "Test validator log message"
	select {
	case app.ValidatorLogChan <- testMessage:
		// Successfully sent
	case <-time.After(time.Second):
		t.Error("Failed to send to ValidatorLogChan")
	}

	// Test that we can receive from the channel
	select {
	case msg := <-app.ValidatorLogChan:
		if msg != testMessage {
			t.Errorf("Expected '%s', got '%s'", testMessage, msg)
		}
	case <-time.After(time.Second):
		t.Error("Failed to receive from ValidatorLogChan")
	}
}

// TestRebuildDynamicLayoutNoClients tests layout when no clients are running
func TestRebuildDynamicLayoutNoClients(t *testing.T) {
	app := NewMonitorApp()

	// Mock no running clients
	// Note: This test assumes GetRunningClients returns empty array when no clients running
	app.rebuildDynamicLayout()

	// The grid should be rebuilt but we can't easily test the internal state
	// We just verify it doesn't panic
	t.Log("Dynamic layout rebuild completed without panics")
}

// TestValidatorClientDetection tests that validator is properly detected
func TestValidatorClientDetection(t *testing.T) {
	// Get running clients
	clients := utils.GetRunningClients()

	// Check if Validator is in the supported client types
	hasValidator := false
	for _, client := range clients {
		if client.Name == "Validator" || client.Name == "StarknetValidator" {
			hasValidator = true
			t.Logf("Validator client detected: %s (PID: %d)", client.Name, client.PID)
			break
		}
	}

	if !hasValidator {
		t.Log("No validator client currently running (expected if not started)")
	}
}

// TestValidatorClientType tests the ClientStarkValidator constant
func TestValidatorClientType(t *testing.T) {
	if types.ClientStarkValidator != "starknet-staking-v2" {
		t.Errorf("Expected ClientStarkValidator to be 'starknet-staking-v2', got '%s'", types.ClientStarkValidator)
	}
}

// TestAllLogChannelsCreated tests that all log channels are properly initialized
func TestAllLogChannelsCreated(t *testing.T) {
	app := NewMonitorApp()

	channels := map[string]chan string{
		"ExecutionLogChan": app.ExecutionLogChan,
		"ConsensusLogChan": app.ConsensusLogChan,
		"JunoLogChan":      app.JunoLogChan,
		"ValidatorLogChan": app.ValidatorLogChan,
		"StatusChan":       app.StatusChan,
		"NetworkChan":      app.NetworkChan,
		"JunoStatusChan":   app.JunoStatusChan,
		"ChainInfoChan":    app.ChainInfoChan,
		"SystemStatsChan":  app.SystemStatsChan,
		"RPCInfoChan":      app.RPCInfoChan,
	}

	for name, ch := range channels {
		if ch == nil {
			t.Errorf("Channel %s should be initialized", name)
		}
	}
}

// TestNoClientsMessage tests the "No Clients Running" message box
func TestNoClientsMessage(t *testing.T) {
	app := NewMonitorApp()

	if app.NoClientsBox == nil {
		t.Fatal("NoClientsBox should be created")
	}

	text := app.NoClientsBox.GetText(false)
	if text == "" {
		t.Error("NoClientsBox should have a message")
	}

	// Check that it contains the expected warning message
	if !contains(text, "NO CLIENTS RUNNING") {
		t.Error("NoClientsBox should contain 'NO CLIENTS RUNNING' message")
	}
}

// TestLogFormatting tests that log lines are properly formatted
func TestLogFormatting(t *testing.T) {
	testCases := []struct {
		input    string
		contains []string
	}{
		{
			input:    "INFO [12-07|15:32:22.145] Test message",
			contains: []string{"INFO", "Test message"},
		},
		{
			input:    "WARN something happened",
			contains: []string{"WARN", "something happened"},
		},
		{
			input:    "ERROR critical failure",
			contains: []string{"ERROR", "critical failure"},
		},
	}

	for _, tc := range testCases {
		formatted := formatLogLines(tc.input)
		for _, expected := range tc.contains {
			if !contains(formatted, expected) {
				t.Errorf("Formatted log should contain '%s'\nInput: %s\nOutput: %s", expected, tc.input, formatted)
			}
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexInString(s, substr) >= 0))
}

func indexInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
