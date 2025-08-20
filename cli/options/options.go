package options

import (
	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg/installer"
)

var (
	ConsensusClient string
	ExecutionClient string
	StarknetClient  string
	Installer       = installer.NewInstaller()
)

func InitGlobalOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ConsensusClient, "consensus_client", "c", "", "Specify the consensus client")
	cmd.PersistentFlags().StringVarP(&ExecutionClient, "execution_client", "e", "", "Specify the execution client")
	cmd.PersistentFlags().StringVarP(&StarknetClient, "starknet_client", "s", "", "Specify the Starknet client")
}
