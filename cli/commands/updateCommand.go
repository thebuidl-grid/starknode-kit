package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/thebuidl-grid/starknode-kit/pkg"
	"github.com/thebuidl-grid/starknode-kit/pkg/updater"

	"github.com/spf13/cobra"
)

var (
	checkOnly   bool
	useOnline   bool
	clientName  string
	autoConfirm bool
)

var UpdateCommand = &cobra.Command{
	Use:   "update [client]",
	Short: "Check for and install client updates",
	Long: `Check if newer versions are available for Ethereum clients and optionally install them.

Supported clients:
  - Execution clients: geth, reth
  - Consensus clients: lighthouse, prysm  
  - Starknet clients: juno

Examples:
  github.com/thebuidl-grid/starknode-kit update                    # Check all clients for updates
  github.com/thebuidl-grid/starknode-kit update geth               # Update specific client
  github.com/thebuidl-grid/starknode-kit update --check-only       # Only check, don't install
  github.com/thebuidl-grid/starknode-kit update --online           # Fetch latest versions online
  github.com/thebuidl-grid/starknode-kit update --auto-confirm     # Auto-confirm all updates`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdate,
}

func init() {
	UpdateCommand.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for updates, don't install")
	UpdateCommand.Flags().BoolVar(&useOnline, "online", false, "Fetch latest versions online instead of using static versions")
	UpdateCommand.Flags().BoolVar(&autoConfirm, "auto-confirm", false, "Automatically confirm all updates without prompting")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	installDir := pkg.InstallDir
	updateChecker := updater.NewUpdateChecker(installDir)

	// Determine which clients to check
	eth_clients, err := installer.GetInsalledClients(pkg.InstallClientsDir)
	if err != nil {
		return err
	}
	stark_clients, err := installer.GetInsalledClients(pkg.InstallStarknetDir)
	if err != nil {
		return err
	}

	if useOnline {
		fmt.Println("⏳ Fetching latest versions from GitHub...")
		time.Sleep(1 * time.Second) // Give visual feedback
	}

	clients := append(eth_clients, stark_clients...)

	// Check for updates
	var updatesAvailable []updater.UpdateInfo
	for _, client := range clients {
		updateInfo, err := updateChecker.CheckClientForUpdate(string(client), useOnline)
		if err != nil {
			fmt.Printf("⚠️  Warning: Could not check %s: %v\n", client, err)
			continue
		}

		if updateInfo.UpdateRequired {
			updatesAvailable = append(updatesAvailable, *updateInfo)
		}
	}

	// Display results
	if len(updatesAvailable) == 0 {
		fmt.Println("✅ All checked clients are up to date!")
		return nil
	}

	fmt.Printf("\n📦 Found %d update(s) available:\n\n", len(updatesAvailable))

	// Display update information
	for _, update := range updatesAvailable {
		clientType := getClientTypeEmoji(update.ClientType)
		fmt.Printf("%s %s (%s client)\n", clientType, update.Client, update.ClientType)
		fmt.Printf("   Current: %s → Latest: %s\n\n", update.CurrentVersion, update.LatestVersion)
	}

	// If check-only mode, exit here
	if checkOnly {
		fmt.Println("👀 Check-only mode enabled. No updates will be installed.")
		return nil
	}

	// Confirm updates
	if !autoConfirm {
		fmt.Print("❓ Do you want to proceed with the updates? [y/N]: ")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ Update cancelled.")
			return nil
		}
	}

	// Perform updates
	fmt.Println("\n🚀 Starting updates...")

	var successful, failed int
	for _, update := range updatesAvailable {
		fmt.Printf("\n⬆️  Updating %s...\n", update.Client)

		result := updateChecker.UpdateClient(update.Client)

		if result.Success {
			successful++
			fmt.Printf("✅ %s updated successfully: %s → %s\n",
				update.Client, result.PreviousVersion, result.NewVersion)
		} else {
			failed++
			fmt.Printf("❌ Failed to update %s: %s\n", update.Client, result.Error)
		}
	}

	// Summary
	fmt.Printf("\n📊 Update Summary:\n")
	fmt.Printf("   ✅ Successful: %d\n", successful)
	fmt.Printf("   ❌ Failed: %d\n", failed)

	if failed > 0 {
		fmt.Println("\n⚠️  Some updates failed. Check the error messages above.")
		return fmt.Errorf("update process completed with %d failure(s)", failed)
	}

	fmt.Println("\n🎉 All updates completed successfully!")
	return nil
}

func getClientTypeEmoji(clientType string) string {
	switch clientType {
	case "execution":
		return "⚡"
	case "consensus":
		return "🏛️"
	case "starknet":
		return "🌟"
	default:
		return "🔧"
	}
}

// TODO should not crash if folder does not exist return a message if staknet folder does not exist
