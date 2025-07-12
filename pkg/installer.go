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
	"starknode-kit/pkg/types"
	"starknode-kit/pkg/versions"
	"strings"
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

	// Add a new variable for Starknet clients dir
	StarknetClientsDir = filepath.Join(InstallDir, "starknet_clients")
)

// installer manages Ethereum client installation
type installer struct {
	InstallDir string
}

// Newinstaller creates a new installer instance
func NewInstaller(Installpath string) *installer {
	if err := setupJWTSecret(); err != nil {
		panic(err)
	}
	return &installer{InstallDir: Installpath}
}

// GetClientFileName returns the file name based on platform and architecture
func (i *installer) getClientFileName(client types.ClientType) (string, error) {
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
	case types.ClientGeth:
		// Map Go arch back to geth arch names
		gethArch := "amd64"
		if goarch == "arm64" {
			gethArch = "arm64"
		}
		fileName = fmt.Sprintf("geth-%s-%s-%s-%s",
			goos, gethArch, versions.LatestGethVersion, GethHash[versions.LatestGethVersion])
	case types.ClientReth:
		fileName = fmt.Sprintf("reth-v%s-%s", versions.LatestRethVersion, archName)
	case types.ClientLighthouse:
		fileName = fmt.Sprintf("lighthouse-v%s-%s", versions.LatestLighthouseVersion, archName)
	case types.ClientPrysm:
		fileName = "prysm.sh"
	case types.ClientJuno:
		// Juno is a Go binary, we'll build it from source
		fileName = fmt.Sprintf("juno-%s", versions.LatestJunoVersion)
	default:
		return "", fmt.Errorf("unknown client: %s", client)
	}

	return fileName, nil
}

// getDownloadURL returns the appropriate URL for downloading a client
func (i *installer) getDownloadURL(client types.ClientType, fileName string) (string, error) {
	switch client {
	case types.ClientGeth:
		return fmt.Sprintf("https://gethstore.blob.core.windows.net/builds/%s.tar.gz", fileName), nil
	case types.ClientReth:
		return fmt.Sprintf("https://github.com/paradigmxyz/reth/releases/download/v%s/%s.tar.gz",
			versions.LatestRethVersion, fileName), nil
	case types.ClientLighthouse:
		return fmt.Sprintf("https://github.com/sigp/lighthouse/releases/download/v%s/%s.tar.gz",
			versions.LatestLighthouseVersion, fileName), nil
	case types.ClientPrysm:
		return "https://raw.githubusercontent.com/prysmaticlabs/prysm/master/prysm.sh", nil
	case types.ClientJuno:
		return fmt.Sprintf("https://github.com/NethermindEth/juno/archive/refs/tags/%s.tar.gz", versions.LatestJunoVersion), nil
	default:
		return "", fmt.Errorf("unknown client: %s", client)
	}
}

// InstallClient installs the specified Ethereum client
func (i *installer) InstallClient(client types.ClientType) error {
	// Handle Juno installation separately (npm-based)
	if client == types.ClientJuno {
		return i.installJuno()
	}

	// Get client file name
	fileName, err := i.getClientFileName(client)
	if err != nil {
		return err
	}

	// Create client directory paths
	clientDir := filepath.Join(i.InstallDir, string(client))
	databaseDir := filepath.Join(clientDir, "database")
	logsDir := filepath.Join(clientDir, "logs")

	// Determine the path to the client binary/script
	var clientPath string
	if client == types.ClientPrysm {
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
	downloadURL, err := i.getDownloadURL(client, fileName)
	if err != nil {
		return err
	}

	// Handle installation differently based on client
	if client == types.ClientPrysm {
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
		if client == types.ClientGeth {
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

// installJuno installs Juno by building from Go source code
func (i *installer) installJuno() error {
	//  Create Juno directory
	junoBaseDir := filepath.Join(StarknetClientsDir, string(types.ClientJuno))
	junoRepoDir := filepath.Join(junoBaseDir, "juno")
	databaseDir := filepath.Join(junoBaseDir, "database")
	logsDir := filepath.Join(junoBaseDir, "logs")

	// Check if Juno is already installed
	junoPath := filepath.Join(junoRepoDir, "juno")
	if _, err := os.Stat(junoPath); err == nil {
		fmt.Printf("%s is already installed.\n", types.ClientJuno)
		return nil
	}

	// Create directories
	fmt.Printf("Creating '%s'\n", junoBaseDir)
	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Check if Go is available
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("Go is not installed. Please install Go 1.24 or higher first: %w", err)
	}

	// Check if make is available
	if _, err := exec.LookPath("make"); err != nil {
		return fmt.Errorf("make is not installed. Please install make first: %w", err)
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git is not installed. Please install git first: %w", err)
	}

	// Clone or pull Juno repository
	if _, err := os.Stat(junoRepoDir); os.IsNotExist(err) {
		fmt.Printf("Cloning Juno repository...\n")
		cloneCmd := exec.Command("git", "clone", "https://github.com/NethermindEth/juno.git", "juno")
		cloneCmd.Dir = junoBaseDir
		cloneCmd.Stdout = os.Stdout
		cloneCmd.Stderr = os.Stderr
		if err := cloneCmd.Run(); err != nil {
			return fmt.Errorf("failed to clone Juno repository: %w", err)
		}
	} else {
		fmt.Printf("Juno repository already exists, pulling latest...\n")
		pullCmd := exec.Command("git", "pull")
		pullCmd.Dir = junoRepoDir
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr
		if err := pullCmd.Run(); err != nil {
			// Always try to recover by checking out main or master
			fmt.Printf("git pull failed, attempting to recover by checking out main or master...\n")
			fetchCmd := exec.Command("git", "fetch")
			fetchCmd.Dir = junoRepoDir
			fetchCmd.Stdout = os.Stdout
			fetchCmd.Stderr = os.Stderr
			if err := fetchCmd.Run(); err != nil {
				return fmt.Errorf("failed to fetch in Juno repo: %w", err)
			}
			checkoutMainCmd := exec.Command("git", "checkout", "main")
			checkoutMainCmd.Dir = junoRepoDir
			checkoutMainCmd.Stdout = os.Stdout
			checkoutMainCmd.Stderr = os.Stderr
			if err := checkoutMainCmd.Run(); err == nil {
				pullMainCmd := exec.Command("git", "pull", "origin", "main")
				pullMainCmd.Dir = junoRepoDir
				pullMainCmd.Stdout = os.Stdout
				pullMainCmd.Stderr = os.Stderr
				if err := pullMainCmd.Run(); err != nil {
					return fmt.Errorf("failed to pull origin main in Juno repo: %w", err)
				}
			} else {
				// Try master as fallback
				checkoutMasterCmd := exec.Command("git", "checkout", "master")
				checkoutMasterCmd.Dir = junoRepoDir
				checkoutMasterCmd.Stdout = os.Stdout
				checkoutMasterCmd.Stderr = os.Stderr
				if err := checkoutMasterCmd.Run(); err == nil {
					pullMasterCmd := exec.Command("git", "pull", "origin", "master")
					pullMasterCmd.Dir = junoRepoDir
					pullMasterCmd.Stdout = os.Stdout
					pullMasterCmd.Stderr = os.Stderr
					if err := pullMasterCmd.Run(); err != nil {
						return fmt.Errorf("failed to pull origin master in Juno repo: %w", err)
					}
				} else {
					return fmt.Errorf("failed to pull Juno repository and could not recover by checking out main or master: %w", err)
				}
			}
		}
	}

	// Checkout specific version
	fmt.Printf("Checking out version %s...\n", versions.LatestJunoVersion)
	checkoutCmd := exec.Command("git", "checkout", versions.LatestJunoVersion)
	checkoutCmd.Dir = junoRepoDir
	checkoutCmd.Stdout = os.Stdout
	checkoutCmd.Stderr = os.Stderr
	if err := checkoutCmd.Run(); err != nil {
		return fmt.Errorf("failed to checkout version %s: %w", versions.LatestJunoVersion, err)
	}

	// Install dependencies based on platform
	fmt.Printf("Installing dependencies...\n")
	if runtime.GOOS == "darwin" {
		// macOS dependencies
		brewCmd := exec.Command("brew", "install", "jemalloc", "pkg-config")
		brewCmd.Stdout = os.Stdout
		brewCmd.Stderr = os.Stderr
		if err := brewCmd.Run(); err != nil {
			// If brew install fails, try with arch -arm64 for Apple Silicon
			fmt.Printf("brew install failed, trying with arch -arm64...\n")
			brewCmdARM := exec.Command("arch", "-arm64", "brew", "install", "jemalloc", "pkg-config")
			brewCmdARM.Stdout = os.Stdout
			brewCmdARM.Stderr = os.Stderr
			if err := brewCmdARM.Run(); err != nil {
				return fmt.Errorf("failed to install macOS dependencies: %w", err)
			}
		}
	} else if runtime.GOOS == "linux" {
		// Linux dependencies
		aptCmd := exec.Command("sudo", "apt-get", "install", "-y", "libjemalloc-dev", "libjemalloc2", "pkg-config", "libbz2-dev")
		aptCmd.Stdout = os.Stdout
		aptCmd.Stderr = os.Stderr
		if err := aptCmd.Run(); err != nil {
			return fmt.Errorf("failed to install Linux dependencies: %w", err)
		}
	}

	// Install Go dependencies
	fmt.Printf("Installing Go dependencies...\n")
	installDepsCmd := exec.Command("make", "install-deps")
	installDepsCmd.Dir = junoRepoDir
	installDepsCmd.Stdout = os.Stdout
	installDepsCmd.Stderr = os.Stderr
	if err := installDepsCmd.Run(); err != nil {
		return fmt.Errorf("failed to install Go dependencies: %w", err)
	}

	// Build Juno
	fmt.Printf("Building Juno...\n")
	buildCmd := exec.Command("make", "juno")
	buildCmd.Dir = junoRepoDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build Juno: %w", err)
	}

	// Move the built binary to the correct location
	buildPath := filepath.Join(junoRepoDir, "build", "juno")
	if err := os.Rename(buildPath, junoPath); err != nil {
		return fmt.Errorf("failed to move Juno binary: %w", err)
	}

	// Make the binary executable
	if err := os.Chmod(junoPath, 0755); err != nil {
		return fmt.Errorf("failed to make Juno binary executable: %w", err)
	}

	// Clean up build artifacts
	fmt.Printf("Cleaning up build artifacts...\n")
	if err := os.RemoveAll(filepath.Join(junoRepoDir, "build")); err != nil {
		return fmt.Errorf("failed to clean up build directory: %w", err)
	}

	fmt.Printf("%s installed successfully.\n", types.ClientJuno)
	return nil
}

// setupJWTSecret creates a JWT secret file for client authentication
func setupJWTSecret() error {

	// Check if JWT already exists
	if _, err := os.Stat(JwtDir); err == nil {
		return nil
	}

	// Create JWT directory
	if err := os.MkdirAll(JwtDir, 0755); err != nil {
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
func (i *installer) RemoveClient(client types.ClientType) error {
	var clientDir string
	if client == types.ClientJuno {
		clientDir = filepath.Join(StarknetClientsDir, string(types.ClientJuno))
	} else {
		clientDir = filepath.Join(i.InstallDir, string(client))
	}

	if _, err := os.Stat(clientDir); err == nil {
		fmt.Printf("Removing %s installation.\n", client)

		// For Juno, we need to clean up Go build artifacts
		if client == types.ClientJuno {
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			if err := os.Chdir(clientDir); err != nil {
				return fmt.Errorf("failed to change to Juno directory: %w", err)
			}
			defer os.Chdir(currentDir)

			if err := os.RemoveAll(".git"); err != nil {
				return fmt.Errorf("failed to remove git repository: %w", err)
			}
			if err := os.RemoveAll("build"); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove build directory: %w", err)
			}
		}

		return os.RemoveAll(clientDir)
	}

	return nil
}

// GetClientVersion gets the installed version of a client
func (i *installer) GetClientVersion(client types.ClientType) (string, error) {
	var clientDir string
	if client == types.ClientJuno {
		clientDir = filepath.Join(StarknetClientsDir, string(types.ClientJuno))
	} else {
		clientDir = filepath.Join(i.InstallDir, string(client))
	}

	// Check if client is installed
	clientPath := filepath.Join(clientDir, string(client))
	if client == types.ClientPrysm {
		clientPath = filepath.Join(clientDir, "prysm.sh")
	} else if client == types.ClientJuno {
		clientPath = filepath.Join(clientDir, "juno")
	}

	if _, err := os.Stat(clientPath); os.IsNotExist(err) {
		return "", fmt.Errorf("%s is not installed", client)
	}

	// Handle Juno version checking differently (npm-based)
	if client == types.ClientJuno {
		return i.getJunoVersion(clientDir)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(i.InstallDir); err != nil {
		return "", fmt.Errorf("failed to change to installation directory: %w", err)
	}
	defer os.Chdir(currentDir)

	version := GetVersionNumber(string(client))
	if version == "" {
		return "", fmt.Errorf("failed to get version for %s", client)
	}

	return version, nil
}

// getJunoVersion gets the installed version of Juno via the binary
func (i *installer) getJunoVersion(junoDir string) (string, error) {
	// junoDir will be passed as StarknetClientsDir/juno by GetClientVersion
	junoPath := filepath.Join(junoDir, "juno")

	if _, err := os.Stat(junoPath); os.IsNotExist(err) {
		return "", fmt.Errorf("Juno binary not found")
	}

	cmd := exec.Command(junoPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get Juno version: %w", err)
	}

	versionOutput := strings.TrimSpace(string(output))
	versionMatch := regexp.MustCompile(`juno version (\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	if len(versionMatch) > 1 {
		return versionMatch[1], nil
	}

	gitCmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	gitCmd.Dir = junoDir
	gitOutput, err := gitCmd.Output()
	if err == nil {
		gitVersion := strings.TrimSpace(string(gitOutput))
		if strings.HasPrefix(gitVersion, "v") {
			gitVersion = gitVersion[1:]
		}
		return gitVersion, nil
	}

	return "", fmt.Errorf("unable to determine Juno version")
}

// IsClientLatestVersion checks if the installed client is the latest version
func (i *installer) IsClientLatestVersion(client types.ClientType, version string) (bool, string) {
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

func CompareClientVersions(client, installedVersion string) (isUpToDate bool, latestVersion string) {
	switch client {
	case "reth":
		latestVersion = versions.LatestRethVersion
	case "geth":
		latestVersion = versions.LatestGethVersion
	case "lighthouse":
		latestVersion = versions.LatestLighthouseVersion
	case "prysm":
		latestVersion = versions.LatestPrysmVersion
	case "juno":
		latestVersion = versions.LatestJunoVersion
	default:
		fmt.Printf("Unknown client: %s\n", client)
		return false, ""
	}

	isUpToDate = compareVersions(installedVersion, latestVersion) >= 0
	return // Implicit return of named values
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
