package commands

import (
	"starknode-kit/pkg"
	"fmt"

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
	err := pkg.CreateStarkNodeConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	return

}

func helpFunction(cmd *cobra.Command, args []string) {
	fmt.Println(`Usage: 
    starknode init

    Initializes the default configuration file (config.yaml).

    This command will create a basic configuration file with default
    execution and consensus port values. If the file already exists,
    it will not be overwritten.

    Example:
      starknode init

    Note: This command does not accept global flags.`)
}

func init() {
	InitCommand.SetHelpFunc(helpFunction)
	InitCommand.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return InitCommand.Help()
	})
}
