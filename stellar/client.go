package stellar

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// Client manages a running stellar-core process
type Client struct {
	cfg Config
}

// New returns a Client for the given config.
// Returns an error if stellar-core is not installed.
func New(cfg Config) (*Client, error) {
	if BinaryPath() == "" {
		return nil, fmt.Errorf("stellar-core is not installed; run 'stellar install' first")
	}
	return &Client{cfg: cfg}, nil
}

// Start launches stellar-core run in the background, writing logs to the managed logs dir
func (c *Client) Start() error {
	if err := os.MkdirAll(filepath.Join(installDir(), "logs"), 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(installDir(), "logs", fmt.Sprintf("stellar-core_%s.log", timestamp))

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	cmd := exec.Command(BinaryPath(), "run", "--conf", c.cfg.configPath())
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start stellar-core: %w", err)
	}

	fmt.Printf("stellar-core started (pid %d) — logs: %s\n", cmd.Process.Pid, logPath)
	return nil
}

// HTTPCommand sends an admin command to the running stellar-core HTTP endpoint
// and returns the raw response body. Common commands: info, peers, quorum, scp, metrics
func (c *Client) HTTPCommand(command string) (string, error) {
	url := fmt.Sprintf("http://localhost:%d/%s", c.cfg.httpPort(), command)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("http-command %q: %w", command, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}
	return string(body), nil
}

// Info returns the output of the stellar-core 'info' HTTP command
func (c *Client) Info() (string, error) {
	return c.HTTPCommand("info")
}

// BinaryPath returns the path to the installed stellar-core binary, or "" if absent
func BinaryPath() string {
	p := filepath.Join(installDir(), "stellar-core")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}

// installDir returns the root managed directory for stellar-core
func installDir() string {
	home, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, "starknode-kit", "stellar_clients", "stellar-core")
}
