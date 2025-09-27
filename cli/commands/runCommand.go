package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/clients"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var RunCmd = &cobra.Command{
	Use:   "run [client]",
	Short: "Run a specific local infrastructure service",
	Long: `Run a specific local infrastructure service by name.

This command starts a single client using its settings from your 'starknode.yaml' configuration file.

Supported clients:
  - geth, reth (Execution)
  - lighthouse, prysm (Consensus)
  - juno (Starknet)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !options.LoadedConfig {
			fmt.Println(utils.Red("❌ No config found."))
			fmt.Println(utils.Yellow("💡 Run `starknode-kit config new` to create a config file."))
			return
		}

		clientName := args[0]
		clientType, err := utils.ResolveClientType(clientName)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid client name: %s", clientName)))
			return
		}
		if !utils.IsInstalled(clientType) {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Client %s not installed", clientName)))
			return
		}
		if options.IsClientRunning(clientType) {
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Client already running: %v", err)))
				return
			}
			fmt.Println(utils.Cyan("⏳ Waiting for log files to be created..."))
			options.LoadLogs([]string{string(clientType)})
		}
		fmt.Println(utils.Cyan(fmt.Sprintf("🚀 Attempting to run %s...", clientName)))

		switch clientType {
		case types.ClientGeth, types.ClientReth:
			// It's an execution client
			if options.Config.ExecutionCientSettings.Name != clientType {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Configured execution client is %s, not %s.", options.Config.ExecutionCientSettings.Name, clientName)))
				return
			}
			eClient, err := clients.NewExecutionClient(options.Config.ExecutionCientSettings, options.Config.Network)
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating execution client: %v", err)))
				return
			}
			if err = eClient.Start(); err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting execution client: %v", err)))
				return
			}
			fmt.Println(utils.Green(fmt.Sprintf("✅ %s started successfully.", clientName)))

		case types.ClientLighthouse, types.ClientPrysm:
			// It's a consensus client
			if options.Config.ConsensusCientSettings.Name != clientType {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Configured consensus client is %s, not %s.", options.Config.ConsensusCientSettings.Name, clientName)))
				return
			}
			cClient, err := clients.NewConsensusClient(options.Config.ConsensusCientSettings, options.Config.Network)
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating consensus client: %v", err)))
				return
			}
			if err = cClient.Start(); err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting consensus client: %v", err)))
				return
			}
			fmt.Println(utils.Green(fmt.Sprintf("✅ %s started successfully.", clientName)))

		case types.ClientJuno:
			j, err := clients.NewJunoClient(options.Config.JunoConfig, options.Config.Network, options.Config.IsValidatorNode)
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating Juno client: %v", err)))
				return
			}
			if err = j.Start(); err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting Juno: %v", err)))
				return
			}
			fmt.Println(utils.Green("✅ Juno started successfully."))

		default:
			fmt.Println(utils.Red(fmt.Sprintf("❌ Don't know how to run client: %s", clientName)))
		}
		fmt.Println(utils.Cyan("⏳ Waiting for log files to be created..."))
		options.LoadLogs([]string{string(clientType)})

	},
}
