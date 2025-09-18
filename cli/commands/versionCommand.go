package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/versions"
)

var VersionCommand = &cobra.Command{
	Use:   "version [client]",
	Short: "Show version of starknode-kit or a specific client",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if versions.StarkNodeVersion == "" {
				versions.StarkNodeVersion = "dev" // fallback for local go run
			}
			fmt.Printf("starknode-kit version %s\n", utils.Green(versions.StarkNodeVersion))
			return
		}

		clientName := args[0]
		clientType, err := utils.ResolveClientType(clientName)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid client name: %s", clientName)))
			return
		}

		if !utils.IsInstalled(clientType) {
			fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client %s is not installed.", clientName)))
			return
		}

    version := versions.GetVersionNumber(clientName)

		fmt.Printf("%s version: %s\n", clientName, utils.Green(version))
	},
}
