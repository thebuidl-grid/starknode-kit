package commands

import (
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"

	"github.com/spf13/cobra"
)

var (
	AddCmd = &cobra.Command{
		Use:   "add [client]",
		Short: "Install a client",
		Long: `Install a client to be used with StarkNode Kit.

		It is recommended to use this command to install clients, as it will
		install the correct versions and configure them for you.

		You can also use this command to install a specific version of a client.

		Example:
		starknode-kit add geth@v1.13.12
		`, // TODO: add version support
		Run: installCommand,
	}
)

func installCommand(cmd *cobra.Command, args []string) {
	options.Installer.InstallClient(types.ClientStarkValidator)
}

func init() {
	options.InitGlobalOptions(AddCmd)
}
