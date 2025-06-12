package commands

import (
	"fmt"
	"starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var ShowConfigCommand = &cobra.Command{
	Use:   "show",
	Short: "show the configured Ethereum clients",
	Long: `The show command shows the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration.`,
	Run: showconfigcommand,
}

func showconfigcommand(cmd *cobra.Command, args []string) {
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)
  return
}
