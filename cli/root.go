package cli

import (
	"os"

	"github.com/thebuidl-grid/starknode-kit/cli/commands"
	configcommand "github.com/thebuidl-grid/starknode-kit/cli/commands/configCommand"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "starknode",
		Short: "Tool for setting up and managing Ethereum and  Starknet nodes",
		Long: `starknode-kit is a CLI tool designed to simplify the setup and management 
of Ethereum and  Starknet nodes. It helps developers quickly configure, 
launch, monitor, and maintain full nodes or validator setups for both networks.

This tool aims to streamline the experience for node operators, 
developers, and testers working with decentralized infrastructure.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Load config globally for all commands
			cfg, err := utils.LoadConfig()
			if err == nil {
				options.Config = cfg
				options.LoadedConfig = true
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Errors are already printed with colors by the commands
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(commands.VersionCommand)
	rootCmd.AddCommand(commands.MonitorCmd)
	rootCmd.AddCommand(commands.StopCommand)
	rootCmd.AddCommand(commands.AddCommand)
	rootCmd.AddCommand(commands.StartCommand)
	rootCmd.AddCommand(commands.RemoveCommand)
	rootCmd.AddCommand(commands.RunCmd)
	rootCmd.AddCommand(commands.UpdateCommand)
	rootCmd.AddCommand(commands.ValidatorCommand)
	rootCmd.AddCommand(configcommand.ConfigCommand)
}

