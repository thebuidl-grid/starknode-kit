package monitoring_test

import (
	"starknode-kit/pkg/monitoring"
	"starknode-kit/pkg/utils"
	"testing"
	"time"
)

func TestNewMonitorApp(t *testing.T) {
	app := monitoring.NewMonitorApp()
	if app == nil {
		t.Fatal("NewMonitorApp returned nil")
	}

	if app.App == nil {
		t.Error("App is nil")
	}
	if app.Grid == nil {
		t.Error("Grid is nil")
	}
	if app.UpdateRate != 2*time.Second {
		t.Error("UpdateRate is not 2 seconds")
	}
}

func TestMonitorAppStartStop(t *testing.T) {
	app := monitoring.NewMonitorApp()

	// Test that we can create the app without errors
	if app == nil {
		t.Fatal("Failed to create monitor app")
	}

	// Test that we can initialize the channels and components
	if app.SystemChan == nil {
		t.Error("SystemChan is nil")
	}
	if app.ClientsChan == nil {
		t.Error("ClientsChan is nil")
	}
	if app.LogsChan == nil {
		t.Error("LogsChan is nil")
	}
	if app.PeersChan == nil {
		t.Error("PeersChan is nil")
	}
	if app.ChainChan == nil {
		t.Error("ChainChan is nil")
	}
	if app.StopChan == nil {
		t.Error("StopChan is nil")
	}

	// Note: We don't test Start() method as it would start the TUI interface
	t.Log("Monitor app components initialized successfully")
}

func TestGetRunningClients(t *testing.T) {
	// This test just ensures the function doesn't panic
	clients := utils.GetRunningClients()
	// It's ok if there are no running clients, we just test that the function works
	// and returns a valid slice (even if empty)
	if clients == nil {
		t.Log("GetRunningClients returned nil (no clients found, which is expected)")
	} else {
		t.Logf("GetRunningClients returned %d clients", len(clients))
	}
}

func TestGetEthereumMetrics(t *testing.T) {
	// This test just ensures the function doesn't panic
	metrics := monitoring.GetEthereumMetrics()
	// We can check if the struct has valid data
	if metrics.NetworkName == "" {
		t.Log("No network name returned (expected if no client is running)")
	}
}
