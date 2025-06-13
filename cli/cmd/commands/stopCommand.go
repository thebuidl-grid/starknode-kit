package commands

import (
	"fmt"
	"starknode-kit/pkg/process"
	"starknode-kit/pkg/utils"

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
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("❌ No config found.")
		fmt.Println("💡 Run `starknode init` to create a config file.")
		return
	}

	elClient := config.ExecutionCientSettings
	clClient := config.ConsensusCientSettings

	fmt.Println("🔍 Checking client processes...")

	elProcess := process.GetProcessInfo(string(elClient.Name))
	clProcess := process.GetProcessInfo(string(clClient.Name))

	if elProcess == nil {
		fmt.Printf("⚠️  Execution client '%s' is not running.\n", elClient.Name)
	} else {
		fmt.Printf("🛑 Stopping execution client '%s'...\n", elClient.Name)
		if err := process.StopClient(string(elClient.Name)); err != nil {
			fmt.Printf("❌ Failed to stop execution client: %v\n", err)
			return
		}
		fmt.Printf("✅ Execution client '%s' stopped successfully.\n", elClient.Name)
	}

	if clProcess == nil {
		fmt.Printf("⚠️  Consensus client '%s' is not running.\n", clClient.Name)
	} else {
		fmt.Printf("🛑 Stopping consensus client '%s'...\n", clClient.Name)
		if err := process.StopClient(string(clClient.Name)); err != nil {
			fmt.Printf("❌ Failed to stop consensus client: %v\n", err)
			return
		}
		fmt.Printf("✅ Consensus client '%s' stopped successfully.\n", clClient.Name)
	}
}

