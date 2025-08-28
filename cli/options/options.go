package options

import (
	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg"
)

var (
	ConsensusClient string
	ExecutionClient string
	StarknetClient  string
	Installer       = pkg.NewInstaller()
)

func InitGlobalOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ConsensusClient, "consensus-client", "c", "", "Specify the consensus client")
	cmd.PersistentFlags().StringVarP(&ExecutionClient, "execution-client", "e", "", "Specify the execution client")
	cmd.PersistentFlags().StringVarP(&StarknetClient, "starknet-client", "s", "", "Specify the Starknet client")
}
