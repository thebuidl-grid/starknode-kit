package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCompareClientVersions(t *testing.T) {
	expectedVersion := "1.2.4"
	installed := "1.2.3"

	oldFunc := getLatestGitHubRelease
	GetLatestGitHubRelease = func(owner, repo string) (string, error) {
		return expectedVersion, nil
	}
	defer func() { GetLatestGitHubRelease = oldFunc }()

	isLatest, latest := CompareClientVersions("reth", installed)
	if isLatest {
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

	installed_clients_dir = "./fakepath"

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
