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
	install, _ := cmd.Flags().GetBool("install")

	// Only deploy account if validator flag is set
	var deployedWallet *types.Wallet
	if validator {
		wallet, err := utils.DeployAccount()
		if err != nil {
			fmt.Printf("Error deploying account: %v\n", err)
			return
		}
		deployedWallet = wallet
	}

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
		if validator {
			utils.DeployAccount()
		}
		defaultConfig.JunoConfig = defaultJunoConfig
	}

	if validator && deployedWallet != nil {
		walletConfig := types.WalletConfig{
			Name: "default",
			Wallet: types.Wallet{
				Address:    "${STARKNET_WALLET}",
				ClassHash:  "${STARKNET_CLASS_HASH}",
				Deployed:   true,
				Legacy:     false,
				PrivateKey: "${STARKNET_PRIVATE_KEY}",
				PublicKey:  "${STARKNET_PUBLIC_KEY}",
				Salt:       "${STARKNET_SALT}",
			},
		}

		// Set the main Wallet field and add to Wallets array
		defaultConfig.Wallet = walletConfig
		defaultConfig.Wallets = []types.WalletConfig{walletConfig}

		// Set up validator config with environment variables
		defaultConfig.ValidatorConfig = types.ValidatorConfig{
			ProviderConfig: struct {
				JunoRPC string `json:"http" yaml:"juno_rpc_http"`
				JunoWS  string `json:"ws" yaml:"juno_rpc_ws"`
			}{
				JunoRPC: "http://localhost:6060",
				JunoWS:  "ws://localhost:6060",
			},
			SignerConfig: struct {
				OperationalAddress string `json:"operational_address"`
				WalletPrivateKey   string `json:"privateKey"`
			}{
				OperationalAddress: "${STARKNET_WALLET}",
				WalletPrivateKey:   "${STARKNET_PRIVATE_KEY}",
			},
		}
	}
	defaultConfig.ConsensusCientSettings = defaultConsensusClientSettings
	defaultConfig.ExecutionCientSettings = defaultExecutionCientSettings

	err := utils.CreateStarkNodeConfig(defaultConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	if install {
		clients := []types.ClientType{defaultConsensusClientSettings.Name, defaultExecutionCientSettings.Name}
		for _, i := range clients {
			err := options.Installer.InstallClient(i)
			if err != nil {
				errMessage := fmt.Sprintf("Could not install client %s\nError: %v", i, err.Error())
				fmt.Println(errMessage)
				return
			}
		}
		if starknet_node {
			err := options.Installer.InstallClient(types.ClientJuno)
			if err != nil {
				errMessage := fmt.Sprintf("Could not install client Juno\nError: %v", err.Error())
				fmt.Println(errMessage)
				return
			}
		}
	}
}

func init() {
	options.InitGlobalOptions(newConfigCommand)
	newConfigCommand.Flags().String("network", "sepolia", "Select the network to connect to (e.g., 'mainnet', 'sepolia')")
	newConfigCommand.Flags().Bool("starknet-node", false, "Install a Starknet node")
	newConfigCommand.Flags().Bool("validator", false, "Configure a validator node (deploys account and sets up wallet config)")
	newConfigCommand.Flags().BoolP("install", "i", true, "Install clients automatically after setup")
}
