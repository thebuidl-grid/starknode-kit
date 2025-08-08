package commands

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	InitCommand = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Starknet node configuration",
		Long: `Creates a default configuration file for a Starknet node. 
This command helps you get started with a new setup by generating a 'starknode.yaml' file with sensible defaults. 
You can customize the configuration by using the available flags.`,
		Run: initCommand,
	}
)

func initCommand(cmd *cobra.Command, args []string) {
	err := utils.CreateStarkNodeConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	options.InitGlobalOptions(InitCommand)
	InitCommand.Flags().String("network", "sepolia", "Select the network to connect to (e.g., 'mainnet', 'sepolia')")
	InitCommand.Flags().Bool("starknet_node", false, "Install a Starknet node")
	InitCommand.Flags().Bool("validator", false, "Configure a validator node")
	InitCommand.Flags().Bool("install", true, "Install clients automatically after setup")
}
