package configcommand

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var ConfigCommand = &cobra.Command{
	Use:   "config",
	Short: "Manage Starknet node configuration",
	Long: `Create, show, and update your Starknet node configuration.
This command allows you to interact with your 'starknode.yaml' file.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// show subcommand with flags
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}

		all, _ := cmd.Flags().GetBool("all")
		el, _ := cmd.Flags().GetBool("el")
		cl, _ := cmd.Flags().GetBool("cl")

		if all {
			showConfig("all")
		} else if el {
			showConfig("el")
		} else if cl {
			showConfig("cl")
		}
	},
}

// set subcommands
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var setNetworkCmd = &cobra.Command{
	Use:   "network [network]",
	Short: "Set network (e.g., 'mainnet', 'sepolia')",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setNetwork(args[0])
	},
}

func showConfig(part string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}

	var title string
	switch part {
	case "all":
		title = "Full Configuration"
	case "el":
		title = "Execution Client"
	case "cl":
		title = "Consensus Client"
	}
	fmt.Printf("%s\n", utils.Cyan(utils.Bold(fmt.Sprintf("---"+" %s ---", title))))

	if part == "all" {
		utils.PrintSection("General")
		utils.PrintKV("Network", options.Config.Network)
		utils.PrintKV("Validator Mode", options.Config.IsValidatorNode)
	}

	if part == "all" || part == "el" {
		if part == "all" {
			utils.PrintSection("Execution Client")
		}
		utils.PrintKV("Client", options.Config.ExecutionCientSettings.Name)
		utils.PrintKV("Type", options.Config.ExecutionCientSettings.ExecutionType)
		utils.PrintKV("Ports", options.Config.ExecutionCientSettings.Port)
	}

	if part == "all" || part == "cl" {
		if part == "all" {
			utils.PrintSection("Consensus Client")
		}
		utils.PrintKV("Client", options.Config.ConsensusCientSettings.Name)
		utils.PrintKV("Ports", options.Config.ConsensusCientSettings.Port)
		utils.PrintKV("Checkpoint", options.Config.ConsensusCientSettings.ConsensusCheckpoint)
	}

	if part == "all" && options.Config.IsValidatorNode {
		utils.PrintSection("Juno Node")
		utils.PrintKV("Port", options.Config.JunoConfig.Port)
		utils.PrintKV("Eth Node", options.Config.JunoConfig.EthNode)
		utils.PrintKV("Environment", options.Config.JunoConfig.Environment)

		utils.PrintSection("Wallet")
		utils.PrintKV("Name", options.Config.Wallet.Name)
		utils.PrintKV("Reward Address", options.Config.Wallet.RewardAddress)
	}
}

func setNetwork(network string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}

	if err := utils.SetNetwork(&options.Config, network); err != nil {
		fmt.Println(utils.Red(err.Error()))
		return
	}

	if err := utils.UpdateStarkNodeConfig(options.Config); err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to save config: %v", err)))
		return
	}

	fmt.Printf("%s\n", utils.Green(fmt.Sprintf("Network set to %s", network)))
}

func init() {
	// Add flags to showCmd
	showCmd.Flags().Bool("all", false, "Show all client settings")
	showCmd.Flags().Bool("el", false, "Show execution client settings")
	showCmd.Flags().Bool("cl", false, "Show consensus client settings")

	// Add set commands
	setCmd.AddCommand(setNetworkCmd)
	setCmd.AddCommand(setCLCmd)
	setCmd.AddCommand(setELCmd)
	setCmd.AddCommand(setStarknetCmd)

	// Add top-level commands to config
	ConfigCommand.AddCommand(showCmd)
	ConfigCommand.AddCommand(setCmd)
	ConfigCommand.AddCommand(newConfigCommand)
}
