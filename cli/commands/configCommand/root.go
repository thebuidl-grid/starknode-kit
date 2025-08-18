package configcommand

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
	"gopkg.in/yaml.v3"
)

var ConfigCommand = &cobra.Command{
	Use:   "config",
	Short: "Create and show the configured Ethereum clients",
	Long: `The show command shows the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration.`,
	Run: configCommand,
}

func configCommand(cmd *cobra.Command, args []string) {
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("No config found")
		fmt.Println("Run `starknode-kit init` to create config file")
		return
	}

	var configBytes []byte
	network, _ := cmd.Flags().GetString("network")
	all, _ := cmd.Flags().GetBool("all")
	el, _ := cmd.Flags().GetBool("el")
	cl, _ := cmd.Flags().GetBool("cl")

	if network != "" {
		err = utils.SetNetwork(&config, network)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = utils.UpdateStarkNodeConfig(config); err != nil {
			fmt.Println(err)
			return
		}
	}

	if all {

		configBytes, err = yaml.Marshal(config)

	} else if el {

		configBytes, err = yaml.Marshal(config.ExecutionCientSettings)

	} else if cl {
		configBytes, err = yaml.Marshal(config.ExecutionCientSettings)

	}
	if err != nil {
		fmt.Println(err)
		return
	}
	if all || el || cl {
		fmt.Println("=== Configuration ===")
		fmt.Println()
		fmt.Println(string(configBytes))
		return
	}
}

func init() {
	ConfigCommand.Flags().StringP("network", "n", "", "set netowork")
	ConfigCommand.Flags().Bool("all", false, "Show all client settings")
	ConfigCommand.Flags().Bool("el", false, "Show execution client settings")
	ConfigCommand.Flags().Bool("cl", false, "Show consensus client settings")
	ConfigCommand.AddCommand(setCLCmd)
	ConfigCommand.AddCommand(setELCmd)
	ConfigCommand.AddCommand(newConfigCommand)
}
