package commands

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var RemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes a specified resource",
	Long: `The remove command allows you to delete a specified resource 
from your application. This command is typically used when cleaning up 
or deprovisioning resources.`,
	Run: removeCommand,
}

func removeCommand(cmd *cobra.Command, args []string) {
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
		fmt.Println(utils.Cyan(fmt.Sprintf("⏳ Removing %s...", client)))
		err = options.Installer.RemoveClient(client)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to remove %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' removed successfully.", client)))
	}
	if options.ExecutionClient != "" {
		client, err := utils.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid execution client: %v", err)))
			return
		}
		fmt.Println(utils.Cyan(fmt.Sprintf("⏳ Removing %s...", client)))
		err = options.Installer.RemoveClient(client)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to remove %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' removed successfully.", client)))
	}
	if options.StarknetClient != "" {
		client, err := utils.GetStarknetClient(options.StarknetClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid Starknet client: %v", err)))
			return
		}
		fmt.Println(utils.Cyan(fmt.Sprintf("⏳ Removing %s...", client)))
		err = options.Installer.RemoveClient(client)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to remove %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' removed successfully.", client)))
	}
}

func init() {
	options.InitGlobalOptions(RemoveCommand)
}
