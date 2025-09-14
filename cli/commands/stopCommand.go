package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var StopCommand = &cobra.Command{
	Use:   "stop [client]",
	Short: "Stop running clients",
	Long: `Stops a specific running client or all clients if the --all flag is provided.

Provide a client name (e.g., geth, lighthouse, juno) to stop a single client.`,
	Args: cobra.MaximumNArgs(1),
	Run:  stopCommand,
}

func stopClient(clientName string) {
	processInfo := process.GetProcessInfo(clientName)
	if processInfo == nil {
		fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client '%s' is not running.", clientName)))
		return
	}

	fmt.Printf("🛑 Stopping client '%s' (PID %d)...", processInfo.Name, processInfo.PID)
	err := process.StopClient(processInfo.PID)
	if err != nil {
		if err.Error() == "os: process already finished" {
			fmt.Printf("ℹ️  Client '%s' is already stopped.\n", processInfo.Name)
		} else {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to stop client '%s': %v", processInfo.Name, err)))
		}
	} else {
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client '%s' stopped successfully.", processInfo.Name)))
	}
}

func stopAllClients() {
	fmt.Println(utils.Cyan("🔍 Stopping all running clients..."))

	runningClients := utils.GetRunningClients()
	if len(runningClients) == 0 {
		fmt.Println(utils.Green("✅ No clients are currently running."))
		return
	}

	for _, client := range runningClients {
		stopClient(client.Name)
	}
}

func stopCommand(cmd *cobra.Command, args []string) {
	if options.Config == nil {
		fmt.Println(utils.Red("❌ No config found."))
		fmt.Println(utils.Yellow("💡 Run `starknode-kit config new` to create a config file."))
		return
	}

	all, _ := cmd.Flags().GetBool("all")

	if all {
		stopAllClients()
		return
	}

	if len(args) > 0 {
		clientName := args[0]
		_, err := utils.ResolveClientType(clientName)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid client name: %s", clientName)))
			return
		}
		stopClient(clientName)
		return
	}

	fmt.Println(utils.Yellow("Please specify a client to stop or use the --all flag."))
	cmd.Help()
}

func init() {
	StopCommand.Flags().Bool("all", false, "Stop all running clients")
}

