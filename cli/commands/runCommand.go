package commands

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/pkg/clients"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

// RunJunoCmd represents the run juno command
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run local Starknet infrastructure services",
	Long: `Run local Starknet infrastructure services using Starknode Kit.

This command serves as a parent for specific components like Juno (a Starknet full node).
You can use subcommands to run individual services such as a Juno node with custom configuration.`,
}

var runJunoCmd = &cobra.Command{
	Use:   "juno",
	Short: "Run a local Juno Starknet node",
	Long: `Run a local Juno Starknet node with configurable options.
Juno is a Go-based Starknet node implementation by Nethermind that provides
full JSON-RPC support for Starknet networks.

Juno requires an Ethereum node connection to verify L1 state. You can specify
an Ethereum node URL using the --eth-node flag.

Example:
  starknode-kit run juno`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := utils.LoadConfig()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		j, err := clients.NewJunoClient(config.JunoConfig, config.Network)
		err = j.Start()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Juno started")
	},
}

func init() {
	RunCmd.AddCommand(runJunoCmd)
}
