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
		Short: "Create a default configuration file",
		Run:   initCommand,
	}
)

func initCommand(cmd *cobra.Command, args []string) {
	err := utils.CreateStarkNodeConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	return

}

func init() {
	options.InitGlobalOptions(InitCommand)
	InitCommand.Flags().Bool("install", false, "Install clients")
	InitCommand.Flags().String("network", "", "Install clients")
	InitCommand.Flags().Bool("validator_node", false, "Is validator node")
}

