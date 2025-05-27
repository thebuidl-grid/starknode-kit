package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "regexp"
    "runtime"
    "strconv"
    "strings"
)

// Version constants
const (
    LatestGethVersion       = "1.15.10"
    LatestRethVersion       = "1.3.4"
    LatestLighthouseVersion = "7.0.1"
)

// GethHash maps Geth versions to their commit hashes
var GethHash = map[string]string{
    "1.14.3":  "ab48ba42",
    "1.14.12": "293a300d",
    "1.15.10": "2bf8a789",
}

// ClientType represents an Ethereum client type
type ClientType string

const (
    ClientGeth       ClientType = "geth"
    ClientReth       ClientType = "reth"
    ClientLighthouse ClientType = "lighthouse"
    ClientPrysm      ClientType = "prysm"
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
func NewInstaller(installDir string) *Installer {
    return &Installer{InstallDir: installDir}
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
    clientDir := filepath.Join(i.InstallDir, "ethereum_clients", string(client))
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
    jwtDir := filepath.Join(i.InstallDir, "ethereum_clients", "jwt")
    jwtPath := filepath.Join(jwtDir, "jwt.hex")
    
    // Check if JWT already exists
    if _, err := os.Stat(jwtPath); err == nil {
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
    if err := os.WriteFile(jwtPath, jwt, 0600); err != nil {
        return fmt.Errorf("failed to write JWT secret: %w", err)
    }
    
    return nil
}

// GetClientVersion returns the installed version of a client
func (i *Installer) GetClientVersion(client ClientType) (string, error) {
    var clientPath string
    var args []string
    
    // Build client command and arguments
    if client == ClientPrysm {
        clientPath = filepath.Join(i.InstallDir, "ethereum_clients", string(client), "prysm.sh")
        args = []string{"beacon-chain", "--version"}
    } else {
        clientPath = filepath.Join(i.InstallDir, "ethereum_clients", string(client), string(client))
        args = []string{"--version"}
    }
    
    // Check if client exists
    if _, err := os.Stat(clientPath); os.IsNotExist(err) {
        return "", fmt.Errorf("client binary not found at %s", clientPath)
    }
    
    // Run version command
    var cmd *exec.Cmd
    if client == ClientPrysm {
        cmd = exec.Command("sh", append([]string{clientPath}, args...)...)
    } else {
        cmd = exec.Command(clientPath, args...)
    }
    
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("error getting version: %w", err)
    }
    
    // Parse version
    version := parseVersionFromOutput(client, string(output))
    if version == "" {
        return "", fmt.Errorf("unable to parse version from output")
    }
    
    return version, nil
}

// parseVersionFromOutput extracts version numbers from client output
func parseVersionFromOutput(client ClientType, output string) string {
    output = strings.TrimSpace(output)
    var re *regexp.Regexp
    
    switch client {
    case ClientGeth:
        re = regexp.MustCompile(`geth version (\d+\.\d+\.\d+)`)
    case ClientReth:
        re = regexp.MustCompile(`reth Version: (\d+\.\d+\.\d+)`)
    case ClientLighthouse:
        re = regexp.MustCompile(`Lighthouse v(\d+\.\d+\.\d+)`)
    case ClientPrysm:
        re = regexp.MustCompile(`beacon-chain-v(\d+\.\d+\.\d+)-`)
    default:
        return ""
    }
    
    matches := re.FindStringSubmatch(output)
    if len(matches) >= 2 {
        return matches[1]
    }
    
    return ""
}

// IsClientLatestVersion checks if a client is at the latest version
func (i *Installer) IsClientLatestVersion(client ClientType, installedVersion string) (bool, string) {
    var latestVersion string
    
    switch client {
    case ClientGeth:
        latestVersion = LatestGethVersion
    case ClientReth:
        latestVersion = LatestRethVersion
    case ClientLighthouse:
        latestVersion = LatestLighthouseVersion
    default:
        return false, ""
    }
    
    return compareVersions(installedVersion, latestVersion) >= 0, latestVersion
}

// compareVersions compares two semver strings
func compareVersions(v1, v2 string) int {
    parts1 := strings.Split(v1, ".")
    parts2 := strings.Split(v2, ".")
    
    for i := 0; i < 3; i++ {
        // Handle potential index out of bounds
        if i >= len(parts1) {
            if i >= len(parts2) {
                return 0 // Equal at this point
            }
            return -1 // v2 has more parts, so it's newer
        }
        if i >= len(parts2) {
            return 1 // v1 has more parts, so it's newer
        }
        
        // Convert to integers and compare
        n1, _ := strconv.Atoi(parts1[i])
        n2, _ := strconv.Atoi(parts2[i])
        
        if n1 > n2 {
            return 1
        }
        if n1 < n2 {
            return -1
        }
    }
    
    return 0 // Versions are equal
}

// RemoveClient removes a client's installation
func (i *Installer) RemoveClient(client ClientType) error {
    clientDir := filepath.Join(i.InstallDir, "ethereum_clients", string(client))
    
    if _, err := os.Stat(clientDir); err == nil {
        fmt.Printf("Removing %s installation.\n", client)
        return os.RemoveAll(clientDir)
    }
    
    return nil
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

func main() {
    // Create installer with default home directory
    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("Error getting home directory: %v\n", err)
        os.Exit(1)
    }
    
    installer := NewInstaller(homeDir)
    cli := &CommandLine{Installer: installer}
    
    if err := cli.Run(os.Args); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}