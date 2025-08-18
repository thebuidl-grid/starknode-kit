package configcommand

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	defaultConfig                  = new(types.StarkNodeKitConfig)
	defaultConsensusClientSettings = types.ClientConfig{
		Name:                types.ClientPrysm,
		Port:                []int{5052, 9000},
		ConsensusCheckpoint: "https://mainnet-checkpoint-sync.stakely.io/",
	}
	defaultExecutionCientSettings = types.ClientConfig{
		Name:          types.ClientGeth,
		Port:          []int{30303},
		ExecutionType: "full",
	}
	defaultJunoConfig = types.JunoConfig{
		Port:    6060,
		EthNode: "wss://eth.drpc.org",
		Environment: []string{
			"JUNO_HTTP_PORT=6060",
			"JUNO_HTTP_HOST=0.0.0.0",
		},
	}
)

var (
	newConfigCommand = &cobra.Command{
		Use:   "new",
		Short: "Create a new Starknet node configuration",
		Long: `Creates a default configuration file for a Starknet node. 
This command helps you get started with a new setup by generating a 'starknode.yaml' file with sensible defaults. 
You can customize the configuration by using the available flags.`,
		Run: runNewConfigCommand,
	}
)

func runNewConfigCommand(cmd *cobra.Command, args []string) {
	network, _ := cmd.Flags().GetString("network")
	starknet_node, _ := cmd.Flags().GetBool("starknet-node")
	validator, _ := cmd.Flags().GetBool("validator")
	//	install, _ := cmd.Flags().GetBool("install")

	if network != "mainnet" && network != "sepolia" {
		errMessage := fmt.Sprintf("Invalid Network: %s", network)
		fmt.Println(errMessage)
		return
	}
	defaultConfig.Network = network

	if options.ConsensusClient != "" {
		client, err := utils.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		defaultConsensusClientSettings.Name = client
	}
	if options.ExecutionClient != "" {
		client, err := utils.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(err)
			return
		}
		defaultExecutionCientSettings.Name = client
	}

	if starknet_node {
		defaultJunoConfig.IsValidatorNode = validator
		defaultConfig.JunoConfig = defaultJunoConfig
	}
	defaultConfig.ConsensusCientSettings = defaultConsensusClientSettings
	defaultConfig.ExecutionCientSettings = defaultExecutionCientSettings

	err := utils.CreateStarkNodeConfig(defaultConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	options.InitGlobalOptions(newConfigCommand)
	newConfigCommand.Flags().String("network", "sepolia", "Select the network to connect to (e.g., 'mainnet', 'sepolia')")
	newConfigCommand.Flags().Bool("starknet-node", false, "Install a Starknet node")
	newConfigCommand.Flags().Bool("validator", false, "Configure a validator node")
	newConfigCommand.Flags().BoolP("install", "i", true, "Install clients automatically after setup")
}
