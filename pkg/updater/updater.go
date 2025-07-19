package updater

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/pkg"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/versions"
)

type UpdateInfo struct {
	Client         string `json:"client"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	UpdateRequired bool   `json:"updateRequired"`
	ClientType     string `json:"clientType"` // "execution", "consensus", "starknet"
}

type UpdateChecker struct {
	installDir string
}

type UpdateResult struct {
	Success         bool   `json:"success"`
	Client          string `json:"client"`
	PreviousVersion string `json:"previousVersion"`
	NewVersion      string `json:"newVersion"`
	BackupPath      string `json:"backupPath,omitempty"`
	Error           string `json:"error,omitempty"`
}

// NewUpdateChecker creates a new update checker instance
func NewUpdateChecker(installDir string) *UpdateChecker {
	return &UpdateChecker{
		installDir: installDir,
	}
}

// CheckAllClientsForUpdates checks all installed clients for updates
func (u *UpdateChecker) CheckAllClientsForUpdates(useOnline bool) ([]UpdateInfo, error) {
	var updates []UpdateInfo

	// Check execution clients
	for _, client := range []string{"geth", "reth"} {
		if updateInfo, err := u.CheckClientForUpdate(client, useOnline); err == nil && updateInfo != nil {
			updates = append(updates, *updateInfo)
		}
	}

	// Check consensus clients
	for _, client := range []string{"lighthouse", "prysm"} {
		if updateInfo, err := u.CheckClientForUpdate(client, useOnline); err == nil && updateInfo != nil {
			updates = append(updates, *updateInfo)
		}
	}

	// Check Starknet clients
	for _, client := range []string{"juno"} {
		if updateInfo, err := u.CheckClientForUpdate(client, useOnline); err == nil && updateInfo != nil {
			updates = append(updates, *updateInfo)
		}
	}

	return updates, nil
}

// CheckClientForUpdate checks if a specific client has an update available
func (u *UpdateChecker) CheckClientForUpdate(client string, useOnline bool) (*UpdateInfo, error) {
	// Get current installed version (this would need to be implemented to check actual installations)
	currentVersion := "1.0.0" // Placeholder - you'd implement actual version detection

	// Get latest version
	var latestVersion string
	if useOnline {
		var err error
		latestVersion, err = u.fetchOnlineVersion(client)
		if err != nil {
			// Fallback to static version
			latestVersion = u.getStaticVersion(client)
			fmt.Printf("Warning: Failed to fetch online version for %s, using static version: %s\n", client, latestVersion)
		}
	} else {
		latestVersion = u.getStaticVersion(client)
	}

	// Compare versions using the existing function
	isLatest, _ := pkg.CompareClientVersions(client, currentVersion)

	updateInfo := &UpdateInfo{
		Client:         client,
		CurrentVersion: currentVersion,
		LatestVersion:  latestVersion,
		UpdateRequired: !isLatest,
		ClientType:     getClientTypeString(client),
	}

	return updateInfo, nil
}

// UpdateClient performs the update for a specific client
func (u *UpdateChecker) UpdateClient(client string) *UpdateResult {
	result := &UpdateResult{
		Client: client,
	}

	// Get current version for backup info (placeholder implementation)
	currentVersion := "1.0.0" // You'd implement actual version detection
	result.PreviousVersion = currentVersion

	// Stop client if running
	if err := u.stopClientIfRunning(client); err != nil {
		result.Error = fmt.Sprintf("Failed to stop %s: %v", client, err)
		return result
	}

	// Create backup
	backupPath, err := u.backupClient(client)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to backup %s: %v", client, err)
		return result
	}
	result.BackupPath = backupPath

	// Create installer instance
	installer := pkg.NewInstaller(u.installDir)
	clientType := getClientType(client)

	// Remove old version (using RemoveClient function if available)
	fmt.Printf("Removing old %s installation...\n", client)

	// Install new version
	if err := installer.InstallClient(clientType); err != nil {
		result.Error = fmt.Sprintf("Failed to install new %s: %v", client, err)
		// Try to restore backup
		u.restoreBackup(client, backupPath)
		return result
	}

	// Get new version (placeholder - you'd implement actual version detection)
	newVersion := u.getStaticVersion(client) // Placeholder
	result.NewVersion = newVersion
	result.Success = true

	return result
}

// fetchOnlineVersion fetches the latest version online for a client
func (u *UpdateChecker) fetchOnlineVersion(client string) (string, error) {
	switch client {
	case "geth":
		return versions.FetchLatestGethVersion()
	case "reth":
		return versions.FetchLatestRethVersion()
	case "lighthouse":
		return versions.FetchLatestLighthouseVersion()
	case "prysm":
		return versions.FetchLatestPrysmVersion()
	case "juno":
		return versions.FetchLatestJunoVersion()
	default:
		return "", fmt.Errorf("unsupported client: %s", client)
	}
}

// getStaticVersion returns the hardcoded version for a client
func (u *UpdateChecker) getStaticVersion(client string) string {
	switch client {
	case "geth":
		return versions.LatestGethVersion
	case "reth":
		return versions.LatestRethVersion
	case "lighthouse":
		return versions.LatestLighthouseVersion
	case "prysm":
		return "latest"
	case "juno":
		return versions.LatestJunoVersion
	default:
		return "unknown"
	}
}

// getClientType converts client name to ClientType
func getClientType(client string) types.ClientType {
	switch client {
	case "geth":
		return types.ClientGeth
	case "reth":
		return types.ClientReth
	case "lighthouse":
		return types.ClientLighthouse
	case "prysm":
		return types.ClientPrysm
	case "juno":
		return types.ClientJuno
	default:
		return ""
	}
}

// getClientTypeString returns the client type category as string
func getClientTypeString(client string) string {
	switch client {
	case "geth", "reth":
		return "execution"
	case "lighthouse", "prysm":
		return "consensus"
	case "juno":
		return "starknet"
	default:
		return "unknown"
	}
}

// stopClientIfRunning stops the client if it's currently running
func (u *UpdateChecker) stopClientIfRunning(client string) error {
	// TODO: Implement client stopping logic
	// This would integrate with your process management system
	fmt.Printf("Stopping %s if running...\n", client)
	return nil
}

// backupClient creates a backup of the client installation
func (u *UpdateChecker) backupClient(client string) (string, error) {
	// TODO: Implement backup logic
	// This would backup the client binary and important config files
	backupPath := fmt.Sprintf("/tmp/%s_backup_%d", client, 123456)
	fmt.Printf("Creating backup for %s at %s...\n", client, backupPath)
	return backupPath, nil
}

// restoreBackup restores a client from backup
func (u *UpdateChecker) restoreBackup(client, backupPath string) error {
	// TODO: Implement restore logic
	fmt.Printf("Restoring %s from backup %s...\n", client, backupPath)
	return nil
}
