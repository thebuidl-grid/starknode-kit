package commands

import (
	"buidlguidl-go/cli/cmd/options"
	"buidlguidl-go/pkg"
	"fmt"

	"github.com/spf13/cobra"
)

var RemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Removes a specified resource",
	Long: `The remove command allows you to delete a specified resource 
from your application. This command is typically used when cleaning up 
or deprovisioning resources.`,
	Run: removeCommand,
}

func removeCommand(cmd *cobra.Command, args []string) {
	if options.ConsensusClient != "" {
		client, err := pkg.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = installer.RemoveClient(client)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if options.ExecutionClient != "" {
		client, err := pkg.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = installer.RemoveClient(client)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	return
}
