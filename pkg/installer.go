package pkg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/versions"
)

var (
	// execCommand allows mocking exec.Command in tests
	execCommand = exec.Command
)

type GitHubTag struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

func getGethHash(version string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/ethereum/go-ethereum/git/ref/tags/%s", version)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tag info: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tagInfo GitHubTag
	if err := json.Unmarshal(body, &tagInfo); err != nil {
		return "", err
	}

	return tagInfo.Object.SHA, nil
}

func getDistro() (string, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.SplitN(line, "=", 2)[1], "\""), nil
		}
	}
	return "", fmt.Errorf("could not determine Linux distro")
}

// installer manages Ethereum client installation
type installer struct {
	InstallDir string
}

// Newinstaller creates a new installer instance
func NewInstaller() installer {
	if err := setupJWTSecret(); err != nil {
		panic(err)
	}
	return installer{InstallDir: constants.InstallClientsDir}
}

func (installer) GetInsalledClients(dir string) ([]types.ClientType, error) {
	clients := make([]types.ClientType, 0)
	validClients := []string{string(types.ClientGeth), string(types.ClientReth), string(types.ClientJuno), string(types.ClientPrysm), string(types.ClientLighthouse)}
	dirclient, err := readFoldersWithReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, i := range dirclient {

		if !slices.Contains(validClients, string(i)) {
			continue
		}
		clients = append(clients, i)
	}

	return clients, nil

}

// GetClientFileName returns the file name based on platform and architecture
func (i *installer) getClientFileName(client types.ClientType, version string) (string, error) {
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
		gethHash, err := getGethHash("v" + version)
		if err != nil {
			return "", err
		}
		fileName = fmt.Sprintf("geth-%s-%s-%s-%s",
			goos, gethArch, version, gethHash[:8])
	case types.ClientReth:
		fileName = fmt.Sprintf("reth-v%s-%s", version, archName)
	case types.ClientLighthouse:
		fileName = fmt.Sprintf("lighthouse-v%s-%s", version, archName)
	case types.ClientPrysm:
		fileName = "prysm.sh"
	case types.ClientJuno:
		fileName = fmt.Sprintf("juno-%s", version)
	case types.ClientStarkValidator:
		if goarch == "amd64" {
			goarch = "x86_64"
		}
		fileName = fmt.Sprintf("starknet-staking-v2_v%s_%s_%s", version, stringTile(goos), goarch)
	default:
		return "", fmt.Errorf("unknown client: %s", client)
	}

	return fileName, nil
}

// TODO use ClientReleaseUrls map
// getDownloadURL returns the appropriate URL for downloading a client
func (i *installer) getDownloadURL(client types.ClientType, fileName, version string) (string, error) {

	switch client {
	case types.ClientGeth:
		return fmt.Sprintf("https://gethstore.blob.core.windows.net/builds/%s.tar.gz", fileName), nil
	case types.ClientReth:
		return fmt.Sprintf("https://github.com/paradigmxyz/reth/releases/download/v%s/%s.tar.gz",
			version, fileName), nil
	case types.ClientLighthouse:
		return fmt.Sprintf("https://github.com/sigp/lighthouse/releases/download/v%s/%s.tar.gz",
			version, fileName), nil
	case types.ClientPrysm:
		return "https://raw.githubusercontent.com/prysmaticlabs/prysm/master/prysm.sh", nil
	case types.ClientJuno:
		return fmt.Sprintf("https://github.com/NethermindEth/juno/archive/refs/tags/%s.tar.gz", version), nil
	case types.ClientStarkValidator:
		return fmt.Sprintf("https://github.com/NethermindEth/starknet-staking-v2/releases/download/v%s/%s.tar.gz", version, fileName), nil
	default:
		return "", fmt.Errorf("unknown client: %s", client)
	}
}

func (i *installer) installClient(client types.ClientType, clientPath, clientDir string) error {
	// Get client file name

	version, err := versions.FetchOnlineVersion(string(client))

	if err != nil {
		return err
	}

	fileName, err := i.getClientFileName(client, version)
	if err != nil {
		return err
	}

	// Get download URL
	downloadURL, err := i.getDownloadURL(client, fileName, version)
	if err != nil {
		return err
	}

	// Install the client
	if err := i.installClientBinary(client, clientDir, clientPath, downloadURL, fileName); err != nil {
		return err
	}

	fmt.Printf("%s installed successfully.\n", client)
	return nil
}

// InstallClient installs the specified Ethereum client
func (i *installer) InstallClient(client types.ClientType) error {

	// Setup client directories
	clientDir := i.getClientDirectory(client)
	if err := i.setupClientDirectories(clientDir); err != nil {
		return err
	}
	clientPath := i.getClientPath(client, clientDir)
	if i.isClientInstalled(clientPath, client) {
		return ErrClientIsInstalled
	}
	return i.installClient(client, clientPath, clientDir)
}

func (i *installer) UpdateClient(client types.ClientType) error {
	clientDir := i.getClientDirectory(client)
	clientPath := i.getClientPath(client, clientDir)
	return i.installClient(client, clientPath, clientDir)
}

// getClientDirectory returns the appropriate directory for the client
func (i *installer) getClientDirectory(client types.ClientType) string {
	if client == types.ClientJuno || client == types.ClientStarkValidator {
		return filepath.Join(constants.InstallStarknetDir, string(client))
	}
	return filepath.Join(constants.InstallClientsDir, string(client))
}

// setupClientDirectories creates the necessary directories for a client
func (i *installer) setupClientDirectories(clientDir string) error {
	databaseDir := filepath.Join(clientDir, "database")
	logsDir := filepath.Join(clientDir, "logs")

	fmt.Printf("Creating '%s'\n", clientDir)
	if err := os.MkdirAll(databaseDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}
	return nil
}

// getClientPath returns the path to the client binary/script
func (i *installer) getClientPath(client types.ClientType, clientDir string) string {
	switch client {
	case types.ClientPrysm:
		return filepath.Join(clientDir, "prysm.sh")
	case types.ClientJuno:
		return filepath.Join(clientDir, "juno")
	default:
		return filepath.Join(clientDir, string(client))
	}
}

// isClientInstalled checks if the client is already installed
func (i *installer) isClientInstalled(clientPath string, client types.ClientType) bool {
	if _, err := os.Stat(clientPath); err == nil {
		fmt.Printf("%s is already installed.\n", client)
		return true
	}
	return false
}

// installClientBinary handles the actual installation of the client binary
func (i *installer) installClientBinary(client types.ClientType, clientDir, clientPath, downloadURL, fileName string) error {
	switch client {
	case types.ClientPrysm:
		return i.installPrysmClient(downloadURL, clientPath, client)
	case types.ClientJuno:
		return i.installJunoClient(client, clientDir, downloadURL, fileName)
	default:
		return i.installStandardClient(client, clientDir, downloadURL, fileName)
	}
}

// installPrysmClient handles Prysm-specific installation (script download)
func (i *installer) installPrysmClient(downloadURL, clientPath string, client types.ClientType) error {
	fmt.Printf("Downloading %s.\n", client)
	if err := downloadFile(downloadURL, clientPath); err != nil {
		return err
	}

	// Make executable
	if err := os.Chmod(clientPath, 0755); err != nil {
		return fmt.Errorf("error making prysm.sh executable: %w", err)
	}
	return nil
}

// installJunoClient handles Juno installation (tar.gz download and extraction)
func (i *installer) installJunoClient(client types.ClientType, clientDir, downloadURL, fileName string) error {
	// Install platform-specific dependencies first
	if err := i.installJunoDependencies(); err != nil {
		return fmt.Errorf("failed to install Juno dependencies: %w", err)
	}

	// Download and extract like other standard clients
	archivePath := filepath.Join(clientDir, fmt.Sprintf("%s.tar.gz", fileName))

	// Download file
	fmt.Printf("Downloading %s.\n", client)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return err
	}

	// Extract archive
	fmt.Printf("Uncompressing %s.\n", client)
	if err := i.extractArchive(archivePath, clientDir); err != nil {
		return err
	}

	// Handle Juno-specific post-extraction (move binary to correct location)
	if err := i.handleJunoPostExtraction(clientDir, fileName); err != nil {
		return err
	}

	// Clean up archive
	fmt.Printf("Cleaning up %s directory.\n", client)
	if err := os.Remove(archivePath); err != nil {
		return fmt.Errorf("error removing archive: %w", err)
	}

	return nil
}

// installStandardClient handles standard client installation (geth, reth, lighthouse)
func (i *installer) installStandardClient(client types.ClientType, clientDir, downloadURL, fileName string) error {
	archivePath := filepath.Join(clientDir, fmt.Sprintf("%s.tar.gz", fileName))

	// Download file
	fmt.Printf("Downloading %s.\n", client)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return err
	}

	// Extract archive
	fmt.Printf("Uncompressing %s.\n", client)
	if err := i.extractArchive(archivePath, clientDir); err != nil {
		return err
	}

	// Handle Geth-specific post-extraction
	if client == types.ClientGeth {
		if err := i.handleGethPostExtraction(clientDir, fileName); err != nil {
			return err
		}
	}

	// Clean up archive
	fmt.Printf("Cleaning up %s directory.\n", client)
	if err := os.Remove(archivePath); err != nil {
		return fmt.Errorf("error removing archive: %w", err)
	}
	return nil
}

// extractArchive extracts a tar.gz archive
func (i *installer) extractArchive(archivePath, clientDir string) error {
	cmd := exec.Command("tar", "-xzvf", archivePath, "-C", clientDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error extracting archive: %w", err)
	}
	return nil
}

// handleGethPostExtraction handles Geth-specific post-extraction steps
func (i *installer) handleGethPostExtraction(clientDir, fileName string) error {
	extractedDir := filepath.Join(clientDir, fileName)
	mvCmd := exec.Command("mv", filepath.Join(extractedDir, "geth"), clientDir)
	if err := mvCmd.Run(); err != nil {
		return fmt.Errorf("error moving geth binary: %w", err)
	}

	// Remove extracted directory
	if err := os.RemoveAll(extractedDir); err != nil {
		return fmt.Errorf("error removing extracted directory: %w", err)
	}
	return nil
}

func (i *installer) handleJunoPostExtraction(clientDir, fileName string) error {
	fileName = strings.Replace(fileName, "v", "", 1)
	extractedDir := filepath.Join(clientDir, fileName)
	junoPath := filepath.Join(clientDir, "juno")
	version_file := filepath.Join(constants.InstallStarknetDir, "juno", ".version")
	version := strings.Replace(fileName, "juno-", "", 1)

	file, err := os.Create(version_file)
	if err != nil {
		return fmt.Errorf("Error creating file:%s", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("juno version %s", version))
	if err != nil {
		return fmt.Errorf("Error writing to file:%s", err)
	}

	fmt.Println("Data written successfully.")

	mvCmd := exec.Command("mv", extractedDir, junoPath)
	if err := mvCmd.Run(); err != nil {
		return fmt.Errorf("error moving juno binary: %w", err)
	}

	buildCmd := exec.Command("make", "juno")
	buildCmd.Dir = junoPath
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("make failed: %v", err)
	}

	return nil
}

// installJunoDependencies installs platform-specific dependencies for Juno
func (i *installer) installJunoDependencies() error {
	fmt.Printf("Installing Juno dependencies...\n")

	if runtime.GOOS == "darwin" {
		return i.installMacOSDependencies()
	} else if runtime.GOOS == "linux" {
		return i.installLinuxDependencies()
	}
	return nil
}

// installMacOSDependencies installs macOS-specific dependencies
func (i *installer) installMacOSDependencies() error {
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
	return nil
}

// installLinuxDependencies installs Linux-specific dependencies
func (i *installer) installLinuxDependencies() error {
	distro, err := getDistro()
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	switch distro {
	case "ubuntu", "debian":
		cmd = exec.Command("sudo", "apt-get", "install", "-y", "libjemalloc-dev", "libjemalloc2", "pkg-config", "libbz2-dev")
	case "fedora":
		cmd = exec.Command("sudo", "dnf", "install", "-y", "jemalloc-devel", "pkgconf-pkg-config", "bzip2-devel")
	case "arch":
		cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "jemalloc", "pkgconf", "bzip2")
	default:
		return fmt.Errorf("unsupported distro: %s", distro)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install packages: %v\n%s", err, stderr.String())
	}

	return nil
}

// setupJWTSecret creates a JWT secret file for client authentication
func setupJWTSecret() error {

	// Check if JWT already exists
	if _, err := os.Stat(constants.JwtDir); err == nil {
		return nil
	}

	// Create JWT directory
	if err := os.MkdirAll(constants.JwtDir, 0755); err != nil {
		return fmt.Errorf("failed to create JWT directory: %w", err)
	}

	// Generate random JWT secret (32 bytes)
	cmd := exec.Command("openssl", "rand", "-hex", "32")
	jwt, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate JWT secret: %w", err)
	}

	// Write JWT to file
	if err := os.WriteFile(constants.JWTPath, jwt, 0600); err != nil {
		return fmt.Errorf("failed to write JWT secret: %w", err)
	}

	return nil
}

// RemoveClient removes a client's installation
func (i *installer) RemoveClient(client types.ClientType) error {
	var clientDir string
	if client == types.ClientJuno {
		clientDir = filepath.Join(constants.InstallStarknetDir, string(types.ClientJuno))
	} else {
		clientDir = filepath.Join(i.InstallDir, string(client))
	}

	if _, err := os.Stat(clientDir); err == nil {

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
	fmt.Printf("Successfully removed %s", client)
	return nil
}

// GetClientVersion gets the installed version of a client
func (i *installer) GetClientVersion(client types.ClientType) (string, error) {
	var clientDir string
	if client == types.ClientJuno {
		clientDir = filepath.Join(constants.InstallStarknetDir, string(types.ClientJuno))
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
		path := filepath.Join(constants.InstallStarknetDir, "juno", ".version")
		version, _ := os.ReadFile(path)
		versionMatch := regexp.MustCompile(`juno version (\d+\.\d+\.\d+)`).FindStringSubmatch(string(version))
		if len(versionMatch) > 1 {
			return versionMatch[1], nil
		}
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(i.InstallDir); err != nil {
		return "", fmt.Errorf("failed to change to installation directory: %w", err)
	}
	defer os.Chdir(currentDir)

	version := versions.GetVersionNumber(string(client))
	if version == "" {
		return "", fmt.Errorf("failed to get version for %s", client)
	}

	return version, nil
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

func CompareClientVersions(client, installedVersion, latestVersion string) bool {
	return compareVersions(installedVersion, latestVersion) >= 0
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

func readFoldersWithReadDir(dirPath string) ([]types.ClientType, error) {
	clients := make([]types.ClientType, 0)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			clients = append(clients, types.ClientType(entry.Name()))
		}
	}
	return clients, nil
}

func stringTile(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(r)) + s[size:]
}
