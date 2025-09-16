package configcommand

import (
	"errors"
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg"
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
	defaultExecutionClientSettings = types.ClientConfig{
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
		fmt.Println(utils.Cyan("🚀 Deploying new wallet for validator..."))
		wallet, err := utils.DeployAccount(network)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Error deploying account: %v", err)))
			return
		}
		deployedWallet = wallet
		fmt.Println(utils.Green("✅ Wallet deployed successfully!"))
	}

	if network != "mainnet" && network != "sepolia" {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid Network: %s", network)))
		return
	}
	defaultConfig.Network = network

	if options.ConsensusClient != "" {
		client, err := utils.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid consensus client: %v", err)))
			return
		}
		defaultConsensusClientSettings.Name = client
	}
	if options.ExecutionClient != "" {
		client, err := utils.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Invalid execution client: %v", err)))
			return
		}
		defaultExecutionClientSettings.Name = client
	}

	if starknet_node {
		defaultConfig.IsValidatorNode = validator
		defaultConfig.JunoConfig = defaultJunoConfig
	}

	// Set up validator configuration if validator flag is set
	if validator && deployedWallet != nil {
		var rewardAddr string
		fmt.Print(utils.Cyan("❓ Enter your reward Address here: "))
		fmt.Scan(&rewardAddr)

		var stakeCommission int
		for {
			fmt.Print(utils.Cyan("❓ Enter your staking commission (1-100): "))
			_, err := fmt.Scan(&stakeCommission)
			if err != nil {
				fmt.Println(utils.Red("❌ Invalid input. Please enter a number."))
				var discard string
				fmt.Scanln(&discard)
				continue
			}
			if stakeCommission >= 1 && stakeCommission <= 100 {
				break
			} else {
				fmt.Println(utils.Red("❌ Commission must be between 1 and 100."))
			}
		}

		// Populate WalletConfig with environment variable syntax
		walletConfig := types.WalletConfig{
			Name:           "default",
			RewardAddress:  rewardAddr,
			StakeCommision: fmt.Sprintf("%d", stakeCommission),
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
	defaultConfig.ExecutionCientSettings = defaultExecutionClientSettings

	fmt.Println(utils.Cyan("📝 Creating starknode.yaml configuration file..."))
	err := utils.CreateStarkNodeConfig(defaultConfig)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating config file: %v", err)))
	}
	fmt.Println(utils.Green("✅ Configuration file 'starknode.yaml' created successfully."))

	if validator {
		fmt.Println(utils.Cyan("💰 Staking STARK for validator..."))
		loadConfig, err := utils.LoadConfig()
		if err != nil {
			fmt.Println(utils.Red("❌ Could not load config"))
			return
		}
		err = utils.StakeStark(network, loadConfig.Wallet)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Error staking STARK: %v", err)))
			return
		}
		fmt.Println(utils.Green("✅ Staking successful!"))
	}

	if install {
		fmt.Println(utils.Cyan("🚀 Installing clients..."))
		clients := []types.ClientType{defaultConsensusClientSettings.Name, defaultExecutionClientSettings.Name}
		for _, i := range clients {
			err := options.Installer.InstallClient(i)
			if errors.Is(err, pkg.ErrClientIsInstalled) {
				fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client %s is already installed. Skipping.", i)))
				continue
			}
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Could not install client %s: %v", i, err.Error())))
				return
			}
			fmt.Println(utils.Green(fmt.Sprintf("✅ Client %s installed successfully.", i)))
		}
		if starknet_node {
			err := options.Installer.InstallClient(types.ClientJuno)
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Could not install client Juno: %v", err.Error())))
				return
			}
			fmt.Println(utils.Green("✅ Client Juno installed successfully."))
			if validator {
				err := options.Installer.InstallClient(types.ClientStarkValidator)
				if err != nil {
					fmt.Println(utils.Red(fmt.Sprintf("❌ Could not install client Starknet Validator: %v", err.Error())))
					return
				}
				fmt.Println(utils.Green("✅ Client Starknet Validator installed successfully."))
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
