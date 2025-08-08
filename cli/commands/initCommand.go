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
}

func init() {
	options.InitGlobalOptions(ConfigCommand)
	ConfigCommand.Flags().String("network", "sepolia", "Select network")
	ConfigCommand.Flags().Bool("validator", false, "Configure validotor node")
	ConfigCommand.Flags().Bool("install", true, "Install clients after setup")
}
