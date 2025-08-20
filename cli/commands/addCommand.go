package commands

import (
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"

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
	options.Installer.InstallClient(types.ClientStarkValidator)
}

func init() {
	options.InitGlobalOptions(InstallCommand)
}
