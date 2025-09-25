package commands

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

// TODO use loggers and not print
var AddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add an Ethereum or Starknet client to the config",
	Long: `The add command registers a new client (such as Prysm, Lighthouse, Geth, Reth, or Juno)
to the local configuration. This sets up the necessary parameters for managing and running
the client as part of your node setup.`,
	Run: addCommand,
}

func addCommand(cmd *cobra.Command, args []string) {
	if cmd.Flags().NFlag() == 0 {
		cmd.Help()
		return
	}
	if options.ConsensusClient != "" {
		client, err := utils.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid consensus client: %v", err)))
			return
		}
		fmt.Println(utils.Cyan(fmt.Sprintf("⏳ Installing %s...", client)))
		err = options.Installer.InstallClient(client)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to install %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' installed successfully.", client)))
	}
	if options.ExecutionClient != "" {
		client, err := utils.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid execution client: %v", err)))
			return
		}
		fmt.Println(utils.Cyan(fmt.Sprintf("⏳ Installing %s...", client)))
		err = options.Installer.InstallClient(client)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to install %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' installed successfully.", client)))
	}
	if options.StarknetClient != "" {
		client, err := utils.GetStarknetClient(options.StarknetClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid Starknet client: %v", err)))
			return
		}
		fmt.Println(utils.Cyan(fmt.Sprintf("⏳ Installing %s...", client)))
		err = options.Installer.InstallClient(client)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to install %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' installed successfully.", client)))
	}
}

func init() {
	options.InitGlobalOptions(AddCommand)
}
