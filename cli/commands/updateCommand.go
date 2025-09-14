package commands

import (
	"fmt"
	"strings"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/updater"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

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
	  starknode-kit update                    # Check all clients for updates
	  starknode-kit update geth               # Update specific client
	  starknode-kit update --check-only       # Only check, don't install
	  starknode-kit update -y                 # Auto-confirm all updates`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUpdate,
}

func init() {
	UpdateCommand.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for updates, don't install")
	UpdateCommand.Flags().BoolVarP(&autoConfirm, "yes", "y", false, "Automatically confirm all updates without prompting")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	installDir := constants.InstallDir
	updateChecker := updater.NewUpdateChecker(installDir)

	// Determine which clients to check
	eth_clients, err := options.Installer.GetInsalledClients(constants.InstallClientsDir)
	if err != nil {
		return err
	}
	stark_clients, err := options.Installer.GetInsalledClients(constants.InstallStarknetDir)
	if err != nil {
		fmt.Println(utils.Yellow("No starknet client installed"))
	}

	if useOnline {
		fmt.Println(utils.Cyan("⏳ Fetching latest versions from GitHub..."))
	}

	clients := append(eth_clients, stark_clients...)

	// Check for updates
	var updatesAvailable []updater.UpdateInfo
	for _, client := range clients {
		updateInfo, err := updateChecker.CheckClientForUpdate(string(client), true)
		if err != nil {
			fmt.Println(utils.Yellow(fmt.Sprintf("⚠️  Warning: Could not check %s: %v", client, err)))
			continue
		}

		if updateInfo.UpdateRequired {
			updatesAvailable = append(updatesAvailable, *updateInfo)
		}
	}

	// Display results
	if len(updatesAvailable) == 0 {
		fmt.Println(utils.Green("✅ All checked clients are up to date!"))
		return nil
	}

	fmt.Printf("\n📦 Found %d update(s) available:\n\n", len(updatesAvailable))

	// Display update information
	for _, update := range updatesAvailable {
		clientType := getClientTypeEmoji(update.ClientType)
		fmt.Printf("%s %s (%s client)\n", clientType, utils.Bold(update.Client), update.ClientType)
		fmt.Printf("   Current: %s → Latest: %s\n\n", utils.Red(update.CurrentVersion), utils.Green(update.LatestVersion))
	}

	// If check-only mode, exit here
	if checkOnly {
		fmt.Println(utils.Yellow("👀 Check-only mode enabled. No updates will be installed."))
		return nil
	}

	// Confirm updates
	if !autoConfirm {
		fmt.Print(utils.Cyan("❓ Do you want to proceed with the updates? [y/N]: "))
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println(utils.Red("❌ Update cancelled."))
			return nil
		}
	}

	// Perform updates
	fmt.Println(utils.Cyan("\n🚀 Starting updates..."))

	var successful, failed int
	for _, update := range updatesAvailable {
		fmt.Printf("\n⬆️  Updating %s...\n", update.Client)

		result := updateChecker.UpdateClient(update.Client)

		if result.Success {
			successful++
			fmt.Println(utils.Green(fmt.Sprintf("✅ %s updated successfully: %s → %s",
				update.Client, result.PreviousVersion, result.NewVersion)))
		} else {
			failed++
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to update %s: %s", update.Client, result.Error)))
		}
	}

	// Summary
	fmt.Printf("\n📊 Update Summary:\n")
	fmt.Printf("   ✅ Successful: %d\n", successful)
	fmt.Printf("   ❌ Failed: %d\n", failed)

	if failed > 0 {
		fmt.Println(utils.Yellow("\n⚠️  Some updates failed. Check the error messages above."))
		return fmt.Errorf("update process completed with %d failure(s)", failed)
	}

	fmt.Println(utils.Green("\n🎉 All updates completed successfully!"))
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