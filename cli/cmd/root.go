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
		Use:   "starknode",
		Short: "Tool for setting up and managing Ethereum and  Starknet nodes",
		Long: `starknode is a CLI tool designed to simplify the setup and management 
of Ethereum and  Starknet nodes. It helps developers quickly configure, 
launch, monitor, and maintain full nodes or validator setups for both networks.

This tool aims to streamline the experience for node operators, 
developers, and testers working with decentralized infrastructure.`,
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
	rootCmd.AddCommand(commands.SetCommand)
	rootCmd.AddCommand(commands.InitCommand)
	rootCmd.AddCommand(commands.RemoveCommand)
}
