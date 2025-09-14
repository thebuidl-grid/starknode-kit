package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/clients"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

var RunCmd = &cobra.Command{
	Use:   "run [client]",
	Short: "Run a local infrastructure service",
	Long: `Run a local infrastructure service, such as a Starknet node.

Currently supported clients:
  - juno`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if options.Config == nil {
			fmt.Println(utils.Red("❌ No config found."))
			fmt.Println(utils.Yellow("💡 Run `starknode-kit config new` to create a config file."))
			return
		}

		clientName := args[0]

		switch clientName {
		case "juno":
			fmt.Println(utils.Cyan("🚀 Starting Juno node..."))
			j, err := clients.NewJunoClient(options.Config.JunoConfig, options.Config.Network, options.Config.IsValidatorNode)
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating Juno client: %v", err)))
				return
			}
			err = j.Start()
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting Juno: %v", err)))
				return
			}
			fmt.Println(utils.Green("✅ Juno started successfully."))

			if options.Config.IsValidatorNode {
				fmt.Println(utils.Cyan("🚀 Starting Validator client..."))
				validatorNode, err := clients.NewValidatorClient(options.Config.ValidatorConfig)
				if err != nil {
					fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating validator client: %v", err)))
					return
				}
				err = validatorNode.Start()
				if err != nil {
					fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting validator client: %v", err)))
					return
				}
				fmt.Println(utils.Green("✅ Validator client started successfully."))
			}

		default:
			fmt.Println(utils.Red(fmt.Sprintf("❌ Unknown client: %s", clientName)))
			fmt.Println(utils.Yellow("Currently, only 'juno' is supported by the run command."))
		}
	},
}

func init() {
	// No subcommands to add anymore
}