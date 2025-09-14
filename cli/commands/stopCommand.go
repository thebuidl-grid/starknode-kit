package commands

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var StopCommand = &cobra.Command{
	Use:   "stop",
	Short: "stop the configured Ethereum clients",
	Long: `The stop command stops the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration.`,
	Run: stopCommand,
}

func stopCommand(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ No config found."))
		fmt.Println(utils.Yellow("💡 Run `starknode-kit config new` to create a config file."))
		return
	}

	fmt.Println(utils.Cyan("🔍 Checking for running clients..."))

	runningClients := utils.GetRunningClients()
	if len(runningClients) == 0 {
		fmt.Println(utils.Green("✅ No clients are currently running."))
		return
	}

	for _, client := range runningClients {

		fmt.Printf("🛑 Stopping client '%s' (PID %d)...\n", client.Name, client.PID)
		err := process.StopClient(client.PID)
		if err != nil {
			// Special case for already-finished process
			if err.Error() == "os: process already finished" {
				fmt.Printf("ℹ️  Client '%s' is already stopped.\n", client.Name)
			} else {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to stop client '%s': %v", client.Name, err)))
				continue
			}
		} else {
			fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' stopped successfully.", client.Name)))
		}
	}
}
