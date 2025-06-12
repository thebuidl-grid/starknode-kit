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
		fmt.Println("No config found")
		fmt.Println("Run `starknode init` to create config file")
		return
	}
	elClient := config.ExecutionCientSettings
	clClient := config.ConsensusCientSettings

	elprocess := process.GetProcessInfo(string(elClient.Name))
	clprocess := process.GetProcessInfo(string(clClient.Name))

	if elprocess == nil {
		fmt.Println(fmt.Sprintf("client %s is not running", elClient.Name))
		return
	}

	if clprocess == nil {
		fmt.Println(fmt.Sprintf("client %s is not running", clClient.Name))
		return
	}

	if err := process.StopClient(string(elClient.Name)); err != nil {
		fmt.Println(err)
		return
	}
	if err := process.StopClient(string(clClient.Name)); err != nil {
		fmt.Println(err)
		return
	}
}
