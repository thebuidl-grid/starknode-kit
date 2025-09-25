package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/clients"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Run the configured Ethereum clients",
	Long: `The run command starts the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration. This executes the clients using the
defined settings and manages them as part of your node stack.`,
	Run: startCommand,
}

func startCommand(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ No config found."))
		fmt.Println(utils.Yellow("💡 Run `starknode-kit config new` to create a config file."))
		return
	}
	el := options.Config.ExecutionCientSettings
	cl := options.Config.ConsensusCientSettings
	elClientName := string(el.Name)
	clClientName := string(cl.Name)

	// Validate clients are installed
	elClientType, err := utils.GetExecutionClient(elClientName)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid execution client in config: %%v", err)))
		return
	}
	clClientType, err := utils.GetConsensusClient(clClientName)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid consensus client in config: %%v", err)))
		return
	}
	if !utils.IsInstalled(elClientType) {
		fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client '%s' is not installed.", elClientType)))
		fmt.Printf("Please run: starknode-kit add -e %s\n", elClientType)
		return
	}
	if !utils.IsInstalled(clClientType) {
		fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client '%s' is not installed.", clClientType)))
		fmt.Printf("Please run: starknode-kit add -c %s\n", clClientType)
		return
	}

	fmt.Println(utils.Cyan("🚀 Starting consensus and execution clients in the background..."))
	cClient, err := clients.NewConsensusClient(cl, options.Config.Network)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating consensus client: %%v", err)))
		return
	}
	eClient, err := clients.NewExecutionClient(el, options.Config.Network)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating execution client: %%v", err)))
		return
	}

	if err = cClient.Start(); err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting consensus client: %%v", err)))
		return
	}
	if err = eClient.Start(); err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting execution client: %%v", err)))
		return
	}
	fmt.Println(utils.Green("✅ Clients started successfully in the background."))
}
