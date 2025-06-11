package commands

import (
	"fmt"
	"starknode-kit/pkg"
	"starknode-kit/pkg/clients"

	"github.com/spf13/cobra"
)

var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Run the configured Ethereum clients",
	Long: `The run command starts the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration. This executes the clients using the
defined settings and manages them as part of your node stack.`,
	Run: startCommand,
}

func startCommand(cmd *cobra.Command, args []string) {
	config, err := pkg.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	el := config.ExecutionCientSettings
	cl := config.ConsensusCientSettings
	elClient, err := pkg.GetExecutionClient(string(el.Name))
	if err != nil {
		fmt.Println("Supported execution clients are:")
		fmt.Println(" - geth")
		fmt.Println(" - reth")
		return
	}
	clClient, err := pkg.GetConsensusClient(string(cl.Name))
	if err != nil {
		fmt.Println("Supported consensus clients are:")
		fmt.Println(" - lighhouse")
		fmt.Println(" - prysm")
		return
	}

	err = pkg.IsInstalled(elClient)
	if err != nil {
		fmt.Printf("Client \"%s\" is not installed.\n", elClient)
		fmt.Printf("Please run: starknode add -e %s\n", elClient)
		return
	}
	err = pkg.IsInstalled(clClient)

	if err != nil {
		fmt.Printf("Client \"%s\" is not installed.\n", clClient)
		fmt.Printf("Please run: starknode add -c %s\n", clClient)
		return
	}
	cClient, err := clients.NewConsensusClient(cl)
	if err != nil {
		fmt.Println(err)
		return
	}
	eClient, err := clients.NewExecutionClient(el)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = cClient.Start(); err != nil {
		fmt.Println(err)
		return
	}
	if err = eClient.Start(); err != nil {
		fmt.Println(err)
		return
	}
	return
}
