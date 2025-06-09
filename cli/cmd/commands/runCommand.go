package commands

import (
	"starknode-kit/pkg"
	"starknode-kit/pkg/clients"
	"fmt"

	"github.com/spf13/cobra"
)

var RunCommand = &cobra.Command{
	Use:   "run",
	Short: "Run the configured Ethereum clients",
	Long: `The run command starts the Ethereum clients (e.g., Prysm, Lighthouse, Geth, etc.)
that have been added to your local configuration. This executes the clients using the
defined settings and manages them as part of your node stack.`,
	Run: runcommand,
}

func runcommand(cmd *cobra.Command, args []string) {
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
	switch clClient {
	case pkg.ClientLighthouse:
		if err = clients.StartLightHouse(cl.Port...); err != nil {
			fmt.Println(err)
			return
		}
	case pkg.ClientPrysm:
		if err = clients.StartPrsym(cl.Port...); err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Printf("Client \"%s\" is not installed.\n", clClient)
		fmt.Printf("Please run: starknode add -c %s\n", clClient)
		return
	}
	switch elClient {
	case pkg.ClientGeth:
		if err = clients.StartGeth(el.ExecutionType, el.Port); err != nil {
			fmt.Println(err)
			return
		}
	case pkg.ClientReth:
		if err = clients.StartReth(el.ExecutionType, el.Port); err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Printf("Client \"%s\" is not installed.\n", elClient)
		fmt.Printf("Please run: starknode add -e %s\n", elClient)
		return
	}

}
