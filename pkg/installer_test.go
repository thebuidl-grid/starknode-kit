package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCompareClientVersions(t *testing.T) {
	installed := "1.2.3"

	// We're testing the reth client which has a hardcoded LatestRethVersion
	expectedVersion := LatestRethVersion

	isLatest, latest := CompareClientVersions("reth", installed)
	if compareVersions(installed, expectedVersion) >= 0 && !isLatest {
		t.Errorf("Expected latest, got not latest")
	}
	if compareVersions(installed, expectedVersion) < 0 && isLatest {
		t.Errorf("Expected not latest, got latest")
	}
	if latest != expectedVersion {
		t.Errorf("Expected latest %s, got %s", expectedVersion, latest)
	}
}

// -------------------- exec.Command Mocking --------------------

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	return cmd
}

// Test helper function that simulates binary output
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	args := os.Args
	if len(args) > 3 {
		if strings.Contains(args[3], "reth") {
			fmt.Fprint(os.Stdout, "reth Version: 1.2.3")
		} else if strings.Contains(args[3], "lighthouse") {
			fmt.Fprint(os.Stdout, "Lighthouse v2.3.4")
		} else if strings.Contains(args[3], "geth") {
			fmt.Fprint(os.Stdout, "geth version 3.4.5")
		} else if strings.Contains(args[3], "prysm") {
			fmt.Fprint(os.Stdout, "beacon-chain-v4.5.6-commit")
		}
	}

	os.Exit(0)
}

func TestGetVersionNumber(t *testing.T) {
	// Override execCommand with our mock
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []struct {
		client   string
		expected string
	}{
		{"reth", "1.2.3"},
		{"lighthouse", "2.3.4"},
		{"geth", "3.4.5"},
		{"prysm", "4.5.6"},
	}

	for _, tt := range tests {
		t.Run(tt.client, func(t *testing.T) {
			version := GetVersionNumber(tt.client)
			if version != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, version)
			}
		})
	}
}

func TestGetClientFileName(t *testing.T) {
	installer := NewInstaller("/tmp")

	tests := []struct {
		client  ClientType
		wantErr bool
	}{
		{ClientGeth, false},
		{ClientReth, false},
		{ClientLighthouse, false},
		{ClientPrysm, false},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.client), func(t *testing.T) {
			fileName, err := installer.GetClientFileName(tt.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClientFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && fileName == "" {
				t.Errorf("GetClientFileName() returned empty filename")
			}
		})
	}
}

func TestGetDownloadURL(t *testing.T) {
	installer := NewInstaller("/tmp")

	tests := []struct {
		client   ClientType
		fileName string
		wantErr  bool
	}{
		{ClientGeth, "geth-linux-amd64-1.15.10-2bf8a789", false},
		{ClientReth, "reth-v1.3.4-x86_64-unknown-linux-gnu", false},
		{ClientLighthouse, "lighthouse-v7.0.1-x86_64-unknown-linux-gnu", false},
		{ClientPrysm, "prysm.sh", false},
		{"unknown", "unknown", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.client), func(t *testing.T) {
			url, err := installer.getDownloadURL(tt.client, tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDownloadURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && url == "" {
				t.Errorf("GetDownloadURL() returned empty url")
			}
		})
	}
}

func TestNewInstaller(t *testing.T) {
	installDir := "/test/path"
	installer := NewInstaller(installDir)

	if installer == nil {
		t.Error("NewInstaller() returned nil")
		return
	}

	if installer.InstallDir != installDir {
		t.Errorf("NewInstaller() InstallDir = %v, want %v", installer.InstallDir, installDir)
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.2.3", 0},
		{"1.2.3", "1.2.4", -1},
		{"1.2.4", "1.2.3", 1},
		{"1.3.0", "1.2.9", 1},
		{"2.0.0", "1.9.9", 1},
		{"1.10.0", "1.9.0", 1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_vs_%s", tt.v1, tt.v2), func(t *testing.T) {
			result := compareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("compareVersions(%s, %s) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

// TestIsClientLatestVersion tests the IsClientLatestVersion method
func TestIsClientLatestVersion(t *testing.T) {
	installer := NewInstaller("/tmp")

	tests := []struct {
		client     ClientType
		version    string
		wantLatest bool
	}{
		{ClientReth, "0.1.0", false},                      // Older version
		{ClientReth, LatestRethVersion, true},             // Latest version
		{ClientReth, "999.999.999", true},                 // Future version
		{ClientGeth, "1.0.0", false},                      // Older version
		{ClientGeth, LatestGethVersion, true},             // Latest version
		{ClientLighthouse, "1.0.0", false},                // Older version
		{ClientLighthouse, LatestLighthouseVersion, true}, // Latest version
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.client, tt.version), func(t *testing.T) {
			isLatest, _ := installer.IsClientLatestVersion(tt.client, tt.version)
			if isLatest != tt.wantLatest {
				t.Errorf("IsClientLatestVersion() isLatest = %v, want %v", isLatest, tt.wantLatest)
			}
		})
	}
}

// MockInstaller represents a mock installer for testing
type MockInstaller struct {
	RemoveClientCalled   bool
	InstallClientCalled  bool
	SetupJWTSecretCalled bool
}

// Override RemoveClient for testing
func (m *MockInstaller) RemoveClient(client ClientType) error {
	m.RemoveClientCalled = true
	return nil
}

// TestDownloadFile tests the downloadFile function with a mock HTTP server
func TestDownloadFile(t *testing.T) {
	// Skip this test for now as it requires setting up a mock HTTP server
	t.Skip("Skipping downloadFile test as it requires a mock HTTP server")
}

// TestInstallClient tests the InstallClient method
func TestInstallClient(t *testing.T) {
	// Skip this test as it requires filesystem operations
	t.Skip("Skipping InstallClient test as it requires filesystem operations")
}

// TestSetupJWTSecret tests the SetupJWTSecret method
func TestSetupJWTSecret(t *testing.T) {
	// Skip this test as it requires filesystem operations
	t.Skip("Skipping SetupJWTSecret test as it requires filesystem operations")
}

// TestRemoveClient tests the RemoveClient method
func TestRemoveClient(t *testing.T) {
	// Skip this test as it requires filesystem operations
	t.Skip("Skipping RemoveClient test as it requires filesystem operations")
}
