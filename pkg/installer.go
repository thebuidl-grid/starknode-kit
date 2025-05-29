package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Version constants
const (
	LatestGethVersion       = "1.15.10"
	LatestRethVersion       = "1.3.4"
	LatestLighthouseVersion = "7.0.1"
)

var (
	// execCommand allows mocking exec.Command in tests
	execCommand = exec.Command

	// This can be overridden in tests

	// GethHash maps Geth versions to their commit hashes
	GethHash = map[string]string{
		"1.14.3":  "ab48ba42",
		"1.14.12": "293a300d",
		"1.15.10": "2bf8a789",
	}
)

// ClientConfig holds the download configuration for a client
type ClientConfig struct {
	FileName    string
	DownloadURL string
	BinaryPath  string
}

// Installer manages Ethereum client installation
type Installer struct {
	InstallDir string
}

// NewInstaller creates a new installer instance
func NewInstaller() *Installer {
	return &Installer{InstallDir: InstallClientsDir}
}

// GetClientFileName returns the file name based on platform and architecture
func (i *Installer) GetClientFileName(client ClientType) (string, error) {
	// Get current OS and architecture
	goos := runtime.GOOS     // "darwin", "linux", "windows"
	goarch := runtime.GOARCH // "amd64", "arm64"

	// Convert Go arch to client-specific arch names
	var archName string
	switch goarch {
	case "amd64":
		if goos == "darwin" {
			archName = "x86_64-apple-darwin"
		} else if goos == "linux" {
			archName = "x86_64-unknown-linux-gnu"
		} else {
			return "", fmt.Errorf("unsupported OS: %s", goos)
		}
	case "arm64":
		if goos == "darwin" {
			archName = "aarch64-apple-darwin"
		} else if goos == "linux" {
			archName = "aarch64-unknown-linux-gnu"
		} else {
			return "", fmt.Errorf("unsupported OS: %s", goos)
		}
	default:
		return "", fmt.Errorf("unsupported architecture: %s", goarch)
	}

	// Determine filename based on client
	var fileName string
	switch client {
	case ClientGeth:
		// Map Go arch back to geth arch names
		gethArch := "amd64"
		if goarch == "arm64" {
			gethArch = "arm64"
		}
		fileName = fmt.Sprintf("geth-%s-%s-%s-%s",
			goos, gethArch, LatestGethVersion, GethHash[LatestGethVersion])
	case ClientReth:
		fileName = fmt.Sprintf("reth-v%s-%s", LatestRethVersion, archName)
	case ClientLighthouse:
		fileName = fmt.Sprintf("lighthouse-v%s-%s", LatestLighthouseVersion, archName)
	case ClientPrysm:
		fileName = "prysm.sh"
	default:
		return "", fmt.Errorf("unknown client: %s", client)
	}

	return fileName, nil
}

// GetDownloadURL returns the appropriate URL for downloading a client
func (i *Installer) GetDownloadURL(client ClientType, fileName string) (string, error) {
	switch client {
	case ClientGeth:
		return fmt.Sprintf("https://gethstore.blob.core.windows.net/builds/%s.tar.gz", fileName), nil
	case ClientReth:
		return fmt.Sprintf("https://github.com/paradigmxyz/reth/releases/download/v%s/%s.tar.gz",
			LatestRethVersion, fileName), nil
	case ClientLighthouse:
		return fmt.Sprintf("https://github.com/sigp/lighthouse/releases/download/v%s/%s.tar.gz",
			LatestLighthouseVersion, fileName), nil
	case ClientPrysm:
		return "https://raw.githubusercontent.com/prysmaticlabs/prysm/master/prysm.sh", nil
	default:
		return "", fmt.Errorf("unknown client: %s", client)
	}
}

// InstallClient installs the specified Ethereum client
func (i *Installer) InstallClient(client ClientType) error {
	// Get client file name
	fileName, err := i.GetClientFileName(client)
	if err != nil {
		return err
	}

	// Create client directory paths
	clientDir := filepath.Join(i.InstallDir, string(client))
	databaseDir := filepath.Join(clientDir, "database")
	logsDir := filepath.Join(clientDir, "logs")

	// Determine the path to the client binary/script
	var clientPath string
	if client == ClientPrysm {
		clientPath = filepath.Join(clientDir, "prysm.sh")
	} else {
		clientPath = filepath.Join(clientDir, string(client))
	}

	// Check if client is already installed
	if _, err := os.Stat(clientPath); err == nil {
		fmt.Printf("%s is already installed.\n", client)
		return nil
	}

	// Create directories
	fmt.Printf("Creating '%s'\n", clientDir)
	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Get download URL
	downloadURL, err := i.GetDownloadURL(client, fileName)
	if err != nil {
		return err
	}

	// Handle installation differently based on client
	if client == ClientPrysm {
		fmt.Println("Downloading Prysm.")
		if err := downloadFile(downloadURL, clientPath); err != nil {
			return err
		}

		// Make executable
		if err := os.Chmod(clientPath, 0755); err != nil {
			return fmt.Errorf("error making prysm.sh executable: %w", err)
		}
	} else {
		// Standard client installation (geth, reth, lighthouse)
		archivePath := filepath.Join(clientDir, fmt.Sprintf("%s.tar.gz", fileName))

		// Download file
		fmt.Printf("Downloading %s.\n", client)
		if err := downloadFile(downloadURL, archivePath); err != nil {
			return err
		}

		// Extract archive
		fmt.Printf("Uncompressing %s.\n", client)
		cmd := exec.Command("tar", "-xzvf", archivePath, "-C", clientDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error extracting archive: %w", err)
		}

		// For Geth, we need to move the binary from the extracted folder
		if client == ClientGeth {
			extractedDir := filepath.Join(clientDir, fileName)
			mvCmd := exec.Command("mv", filepath.Join(extractedDir, "geth"), clientDir)
			if err := mvCmd.Run(); err != nil {
				return fmt.Errorf("error moving geth binary: %w", err)
			}

			// Remove extracted directory
			if err := os.RemoveAll(extractedDir); err != nil {
				return fmt.Errorf("error removing extracted directory: %w", err)
			}
		}

		// Clean up archive
		fmt.Printf("Cleaning up %s directory.\n", client)
		if err := os.Remove(archivePath); err != nil {
			return fmt.Errorf("error removing archive: %w", err)
		}
	}

	fmt.Printf("%s installed successfully.\n", client)
	return nil
}

// SetupJWTSecret creates a JWT secret file for client authentication
func (i *Installer) SetupJWTSecret() error {

	// Check if JWT already exists
	if _, err := os.Stat(jwtDir); err == nil {
		return nil
	}

	// Create JWT directory
	if err := os.MkdirAll(jwtDir, 0755); err != nil {
		return fmt.Errorf("failed to create JWT directory: %w", err)
	}

	// Generate random JWT secret (32 bytes)
	cmd := exec.Command("openssl", "rand", "-hex", "32")
	jwt, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	// Write JWT to file
	if err := os.WriteFile(JWTPath, jwt, 0600); err != nil {
		return fmt.Errorf("failed to write JWT secret: %w", err)
	}

	return nil
}

// RemoveClient removes a client's installation
func (i *Installer) RemoveClient(client ClientType) error {
	clientDir := filepath.Join(i.InstallDir, string(client))

	if _, err := os.Stat(clientDir); err == nil {
		fmt.Printf("Removing %s installation.\n", client)
		return os.RemoveAll(clientDir)
	}

	return nil
}

// GetClientVersion gets the installed version of a client
func (i *Installer) GetClientVersion(client ClientType) (string, error) {
	clientDir := filepath.Join(i.InstallDir, string(client))

	// Check if client is installed
	clientPath := filepath.Join(clientDir, string(client))
	if client == ClientPrysm {
		clientPath = filepath.Join(clientDir, "prysm.sh")
	}

	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		return "", fmt.Errorf("%s is not installed", client)
	}

	// Get the current directory to return to it later
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Change to the installation directory
	if err := os.Chdir(i.InstallDir); err != nil {
		return "", fmt.Errorf("failed to change to installation directory: %w", err)
	}
	defer os.Chdir(currentDir) // Return to original directory when done

	version := GetVersionNumber(string(client))
	if version == "" {
		return "", fmt.Errorf("failed to get version for %s", client)
	}

	return version, nil
}

// IsClientLatestVersion checks if the installed client is the latest version
func (i *Installer) IsClientLatestVersion(client ClientType, version string) (bool, string) {
	isLatest, latestVersion := CompareClientVersions(string(client), version)
	return isLatest, latestVersion
}

// downloadFile downloads a file from a URL to a local path
func downloadFile(url, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// CommandLine represents the command line program
type CommandLine struct {
	Installer *Installer
}

// Run parses and executes command line arguments
func (c *CommandLine) Run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments")
	}

	// Parse flags
	var clientStr, directory string
	var isRemove, isVersion bool

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--client", "-c":
			if i+1 < len(args) {
				clientStr = args[i+1]
				i++
			}
		case "--directory", "-d":
			if i+1 < len(args) {
				directory = args[i+1]
				i++
			}
		case "--remove", "-r":
			isRemove = true
		case "--version", "-v":
			isVersion = true
		}
	}

	// Validate client
	if clientStr == "" {
		return fmt.Errorf("client must be specified with --client")
	}

	var client ClientType
	switch clientStr {
	case "geth":
		client = ClientGeth
	case "reth":
		client = ClientReth
	case "lighthouse":
		client = ClientLighthouse
	case "prysm":
		client = ClientPrysm
	default:
		return fmt.Errorf("unknown client: %s", clientStr)
	}

	// Update installer directory if specified
	if directory != "" {
		c.Installer.InstallDir = directory
	}

	// Execute requested action
	if isRemove {
		return c.Installer.RemoveClient(client)
	} else if isVersion {
		version, err := c.Installer.GetClientVersion(client)
		if err != nil {
			return err
		}

		isLatest, latestVersion := c.Installer.IsClientLatestVersion(client, version)
		fmt.Printf("%s version: %s\n", client, version)

		if !isLatest {
			fmt.Printf("Update available: %s (latest: %s)\n", version, latestVersion)
		} else {
			fmt.Println("You have the latest version.")
		}

		return nil
	} else {
		// Install JWT secret if needed
		if err := c.Installer.SetupJWTSecret(); err != nil {
			return err
		}

		// Install the client
		return c.Installer.InstallClient(client)
	}
}

func GetVersionNumber(client string) string {
	var argument string

	switch client {
	case "reth", "lighthouse", "geth":
		argument = "--version"
	case "prysm":
		argument = "beacon-chain --version"
	default:
		fmt.Printf("Unknown client: %s\n", client)
		return ""
	}

	var clientCommand string
	switch runtime.GOOS {
	case "darwin", "linux":
		if client == "prysm" {
			clientCommand = filepath.Join(InstallClientsDir, client, fmt.Sprintf("%s.sh", client))
		} else {
			clientCommand = filepath.Join(InstallClientsDir, client, client)
		}
	case "windows":
		fmt.Println("getVersionNumber() for windows is not yet implemented")
		os.Exit(1)
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		return ""
	}

	cmdParts := strings.Split(argument, " ")
	cmd := execCommand(clientCommand, cmdParts...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command for %s: %v\n", client, err)
		return ""
	}

	versionOutput := strings.TrimSpace(string(output))
	var versionMatch []string

	switch client {
	case "reth":
		versionMatch = regexp.MustCompile(`reth Version: (\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	case "lighthouse":
		versionMatch = regexp.MustCompile(`Lighthouse v(\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	case "geth":
		versionMatch = regexp.MustCompile(`geth version (\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	case "prysm":
		versionMatch = regexp.MustCompile(`beacon-chain-v(\d+\.\d+\.\d+)-`).FindStringSubmatch(versionOutput)
	}

	if len(versionMatch) > 1 {
		return versionMatch[1]
	}

	fmt.Printf("Unable to parse version number for %s\n", client)
	return ""
}

func CompareClientVersions(client, installedVersion string) (bool, string) {
	var latestVersion string
	switch client {
	case "reth":
		latestVersion = LatestRethVersion
	case "geth":
		latestVersion = LatestGethVersion
	case "lighthouse":
		latestVersion = LatestLighthouseVersion
	case "prysm":
		// Just use a hard-coded latest version for Prysm
		latestVersion = "4.0.5" // Replace with an appropriate version
	default:
		fmt.Printf("Unknown client: %s\n", client)
		return false, ""
	}

	if compareVersions(installedVersion, latestVersion) < 0 {
		return false, latestVersion
	}
	return true, latestVersion
}

func compareVersions(v1, v2 string) int {
	split := func(v string) []int {
		parts := strings.Split(v, ".")
		ints := make([]int, len(parts))
		for i, p := range parts {
			fmt.Sscanf(p, "%d", &ints[i])
		}
		return ints
	}

	a := split(v1)
	b := split(v2)

	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
