package stellar

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Install installs stellar-core via the platform package manager and wires up
// the managed binary symlink. Creates the required directory structure under
// <UserConfigDir>/starknode-kit/stellar_clients/stellar-core/.
func Install() error {
	dir := installDir()
	for _, sub := range []string{"logs", "database"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0755); err != nil {
			return fmt.Errorf("create %s dir: %w", sub, err)
		}
	}

	switch runtime.GOOS {
	case "linux":
		return installLinux(dir)
	case "darwin":
		return installMacOS(dir)
	default:
		return fmt.Errorf("unsupported OS %q — install stellar-core manually", runtime.GOOS)
	}
}

// InitDB runs stellar-core new-db to initialise the database before first start.
// Must be called once after Install and after placing a valid stellar-core.cfg.
func InitDB(cfg Config) error {
	if BinaryPath() == "" {
		return fmt.Errorf("stellar-core is not installed")
	}
	cmd := exec.Command(BinaryPath(), "new-db", "--conf", cfg.configPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("stellar-core new-db: %w", err)
	}
	return nil
}

// Remove deletes the managed stellar-core directory entirely
func Remove() error {
	dir := installDir()
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("remove stellar-core: %w", err)
	}
	fmt.Println("stellar-core removed successfully")
	return nil
}

// InstalledVersion returns the version string reported by the installed binary
func InstalledVersion() string {
	bin := BinaryPath()
	if bin == "" {
		return ""
	}
	out, err := exec.Command(bin, "version").Output()
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(out))
}

// --- platform-specific helpers ---

func installLinux(dir string) error {
	distro, err := linuxDistro()
	if err != nil {
		return err
	}
	switch distro {
	case "ubuntu", "debian":
		return installApt(dir)
	default:
		return fmt.Errorf("unsupported distro %q — install stellar-core manually", distro)
	}
}

// installApt adds the official Stellar apt repo and installs stellar-core
func installApt(dir string) error {
	steps := [][]string{
		{"sudo", "apt-get", "install", "-y", "apt-transport-https", "curl"},
		{"sh", "-c", `curl -sSL https://apt.stellar.org/SDF.asc | sudo tee /etc/apt/trusted.gpg.d/SDF.asc > /dev/null`},
		{"sh", "-c", `echo "deb https://apt.stellar.org focal stable" | sudo tee /etc/apt/sources.list.d/SDF.list`},
		{"sudo", "apt-get", "update"},
		{"sudo", "apt-get", "install", "-y", "stellar-core"},
	}
	for _, args := range steps {
		if err := runPrinted(args[0], args[1:]...); err != nil {
			return err
		}
	}
	return symlinkBinary("/usr/bin/stellar-core", dir)
}

func installMacOS(dir string) error {
	if err := runPrinted("brew", "install", "stellar-core"); err != nil {
		return err
	}
	out, err := exec.Command("which", "stellar-core").Output()
	if err != nil {
		return fmt.Errorf("stellar-core not found after brew install: %w", err)
	}
	return symlinkBinary(string(bytes.TrimSpace(out)), dir)
}

func symlinkBinary(src, dir string) error {
	dst := filepath.Join(dir, "stellar-core")
	_ = os.Remove(dst)
	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("symlink %s -> %s: %w", dst, src, err)
	}
	fmt.Printf("stellar-core linked: %s -> %s\n", dst, src)
	return nil
}

func runPrinted(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func linuxDistro() (string, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.ReadFrom(f)
	for _, line := range bytes.Split(buf.Bytes(), []byte("\n")) {
		if bytes.HasPrefix(line, []byte("ID=")) {
			return string(bytes.Trim(bytes.TrimPrefix(line, []byte("ID=")), `"`)), nil
		}
	}
	return "", fmt.Errorf("could not determine Linux distro from /etc/os-release")
}
