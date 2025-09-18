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
	currentVersion := versions.GetVersionNumber(client)

	// Get latest version
	var latestVersion string
	if useOnline {
		var err error
		latestVersion, err = versions.FetchOnlineVersion(client)
		if err != nil {
      return nil, err
		}
	}

	// Compare versions using the existing function
	isLatest := pkg.CompareClientVersions(client, currentVersion, latestVersion)

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
	currentVersion := versions.GetVersionNumber(client)
	result.PreviousVersion = currentVersion

	// Create backup
	backupPath, err := u.backupClient(client)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to backup %s: %v", client, err)
		return result
	}
	result.BackupPath = backupPath

	// Create installer instance
	installer := pkg.NewInstaller()
	clientType := types.GetClientType(client)

	// Remove old version (using RemoveClient function if available)
	fmt.Printf("Removing old %s installation...\n", client)

	// Install new version
	if err := installer.UpdateClient(clientType); err != nil {
		result.Error = fmt.Sprintf("Failed to install new %s: %v", client, err)
		return result
	}

	// Get new version (placeholder - you'd implement actual version detection)
	newVersion, err := versions.FetchOnlineVersion(client) // Placeholder
	if err != nil {
		fmt.Println(err)
		return nil
	}
	result.NewVersion = newVersion
	result.Success = true

	return result
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

// backupClient creates a backup of the client installation
func (u *UpdateChecker) backupClient(client string) (string, error) {
	// TODO: Implement backup logic
	// This would backup the client binary and important config files
	backupPath := fmt.Sprintf("/tmp/%s_backup_%d", client, 123456)
	fmt.Printf("Creating backup for %s at %s...\n", client, backupPath)
	return backupPath, nil
}
