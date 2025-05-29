package commands

import (
	"buidlguidl-go/cli/cmd/options"
	"buidlguidl-go/pkg"

	"github.com/spf13/cobra"
)

var InstallCommand = &cobra.Command{
	Use:   "add",
	Short: "Add an Ethereum client to the config",
	Long: `The add command registers a new Ethereum client (such as Prysm, Lighthouse, Geth, etc.)
to the local configuration. This sets up the necessary parameters for managing and running
the client as part of your node stack.`,
	RunE: installCommand,
}

func installCommand(cmd *cobra.Command, args []string) error {
	if options.ConsensusClient != "" {
		client, err := pkg.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			return err
		}
		err = pkg.Newinstaller().InstallClient(client)
		if err != nil {
			return err
		}
	}
	if options.ExecutionClient != "" {
		client, err := pkg.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			return err
		}
		err = pkg.Newinstaller().InstallClient(client)
		if err != nil {
			return err
		}
	}

	return nil
}
