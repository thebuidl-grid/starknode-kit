package commands

import (
	"fmt"
	"starknode-kit/cli/cmd/options"
	"starknode-kit/pkg"

	"github.com/spf13/cobra"
)

// TODO use loggers and not print
var InstallCommand = &cobra.Command{
	Use:   "add",
	Short: "Add an Ethereum or Starknet client to the config",
	Long: `The add command registers a new client (such as Prysm, Lighthouse, Geth, Reth, or Juno)
to the local configuration. This sets up the necessary parameters for managing and running
the client as part of your node setup.`,
	Run: installCommand,
}

func installCommand(cmd *cobra.Command, args []string) {
	if options.ConsensusClient != "" {
		client, err := pkg.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = installer.InstallClient(client)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if options.ExecutionClient != "" {
		client, err := pkg.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = installer.InstallClient(client)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if options.StarknetClient != "" {
		client, err := pkg.GetStarknetClient(options.StarknetClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = installer.InstallClient(client)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}
