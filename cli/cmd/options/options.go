package options

import "github.com/spf13/cobra"

var (
	ConsensusClient string
	ExecutionClient string
)

func InitGlobalOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ConsensusClient, "consensus_client", "cl", "", "Specify the consensus client")
	cmd.PersistentFlags().StringVarP(&ExecutionClient, "execution_client", "el", "", "Specify the execution client")
}
