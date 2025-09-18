package options

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var (
	ConsensusClient string
	ExecutionClient string
	StarknetClient  string
	Installer       = pkg.NewInstaller()
	Config          types.StarkNodeKitConfig
	LoadedConfig    = false
)

func InitGlobalOptions(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&ConsensusClient, "consensus-client", "c", "", "Specify the consensus client")
	cmd.PersistentFlags().StringVarP(&ExecutionClient, "execution-client", "e", "", "Specify the execution client")
	cmd.PersistentFlags().StringVarP(&StarknetClient, "starknet-client", "s", "", "Specify the Starknet client")
}

func Wait() {
	fmt.Println(utils.Cyan("Showing logs. Press Ctrl+C to exit."))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
