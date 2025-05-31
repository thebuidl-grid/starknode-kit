package commands

import (
	"buidlguidl-go/cli/cmd/commands"
	"buidlguidl-go/cli/cmd/options"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stacknode",
		Short: "To add",
		Long:  `Add long command here`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is root test command")
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	options.InitGlobalOptions(rootCmd)
	rootCmd.AddCommand(commands.InstallCommand)
	rootCmd.AddCommand(commands.RemoveCommand)
}
