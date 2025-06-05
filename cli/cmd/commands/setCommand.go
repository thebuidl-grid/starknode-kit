package commands

import (
	"buidlguidl-go/cli/cmd/options"
	"buidlguidl-go/pkg"
	"fmt"

	"github.com/spf13/cobra"
)

var SetCommand = &cobra.Command{
	Use:   "set",
	Short: "Set config values",
	Long:  "Set allows you to modify configuration values used by the application. You can set individual keys or multiple values at once.",
	Run:   setCommand,
}

func setCommand(cmd *cobra.Command, args []string) {
	cfg, err := pkg.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	if options.ConsensusClient != "" {
		client, err := pkg.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		cfg.ConsensusCientSettings.Name = pkg.ClientType(client)
	}
	if options.ExecutionClient != "" {
		client, err := pkg.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		cfg.ExecutionCientSettings.Name = pkg.ClientType(client)
	}

	err = pkg.UpdateStackNodeConfig(cfg)
	if err != nil {
		fmt.Println(err)
		return

	}
}
