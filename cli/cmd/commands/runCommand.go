package commands

import "github.com/spf13/cobra"

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Run the configured Ethereum clients",
	Long: `The run command starts the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration. This executes the clients using the
defined settings and manages them as part of your node stack.`,
	Run: runcommand,
}

func runcommand(cmd *cobra.Command, args []string) {}
