package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var StatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Display status of running clients",
	Long: `Displays the status (running, stopped, installed, version, PID) of a specific client or all clients.

Provide a client name (e.g., geth, lighthouse, juno) to see the status of a single client, or run without arguments to see all.`,
	Args: cobra.MaximumNArgs(1),
	Run:  statusCommand,
}

func statusCommand(cmd *cobra.Command, args []string) {

	fmt.Println(utils.Yellow("--- Client Status ---"))

	if len(args) == 0 {
		runningClients := utils.GetRunningClients()
		if len(runningClients) == 0 {
			fmt.Println(utils.Red("❌ No client running"))
			return
		}
		for _, clientType := range runningClients {
			displayClientStatus(types.ClientType(clientType.Name))
		}
	} else {
		clientName := args[0]
		clientType, err := utils.ResolveClientType(clientName)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid client name: %s", clientName)))
			return
		}
		displayClientStatus(clientType)
	}
}

func displayClientStatus(clientType types.ClientType) {
	clientName := string(clientType)
	fmt.Printf("Client: %s\n", utils.Blue(clientName))

	if !utils.IsInstalled(clientType) {
		fmt.Printf("  Status: %s\n", utils.Yellow("Not Installed"))
		return
	}

	version := utils.GetClientVersion(clientName)
	fmt.Printf("  Version: %s\n", utils.Green(version))

	processInfo := process.GetProcessInfo(clientName)
	if processInfo != nil {
		fmt.Printf("  Status: %s (PID: %d)\n", utils.Green("Running"), processInfo.PID)
		fmt.Printf("  Uptime: %s\n", utils.Green(processInfo.Uptime.Round(time.Second).String()))
	} else {
		fmt.Printf("  Status: %s\n", utils.Red("Stopped"))
	}
}

