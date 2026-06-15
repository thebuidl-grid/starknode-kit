package configcommand

import (
	"errors"
	"fmt"

	envsubt "github.com/emperorsixpacks/envsubst"
	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/validator"
	"gopkg.in/yaml.v3"
)

var (
	defaultConsensusClientSettings = types.ClientConfig{
		Name: types.ClientPrysm,
		Port: []int{5052, 9000},
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

// handleValidatorWalletSetup handles the setup of a wallet for a validator.
// It either uses an existing wallet or deploys a new one and prompts for configuration details.
func handleValidatorWalletSetup(network string) (*types.WalletConfig, error) {
	if options.LoadedConfig && options.Config.Wallet.Wallet.Address != "" {
		fmt.Println(utils.Green(fmt.Sprintf("✅ Using already created wallet %s", options.Config.Wallet.Wallet.Address)))
		return &options.Config.Wallet, nil
	}

	fmt.Println(utils.Cyan("🚀 Deploying new wallet for validator..."))
	_, err := utils.DeployAccount(network)
	if err != nil {
		return nil, fmt.Errorf("error deploying account: %w", err)
	}
	fmt.Println(utils.Green("✅ Wallet deployed successfully!"))

	var rewardAddr string
	fmt.Print(utils.Cyan("❓ Enter your reward Address here: "))
	if _, err := fmt.Scan(&rewardAddr); err != nil {
		return nil, errors.New("could not read reward address")
	}

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
		}
		fmt.Println(utils.Red("❌ Commission must be between 1 and 100."))
	}

	walletConfig := &types.WalletConfig{
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
	return walletConfig, nil
}

// generateConfig creates a new StarkNodeKitConfig based on user flags.
func generateConfig(network string, starknetNode bool, validator bool, walletConfig *types.WalletConfig) (*types.StarkNodeKitConfig, error) {
	if network != "mainnet" && network != "sepolia" {
		return nil, fmt.Errorf("invalid Network: %s", network)
	}

	config := new(types.StarkNodeKitConfig)
	config.Network = network

	// Consensus Client
	consensusClientSettings := defaultConsensusClientSettings
	if options.ConsensusClient != "" {
		client, err := utils.GetConsensusClient(options.ConsensusClient)
		if err != nil {
			return nil, fmt.Errorf("invalid consensus client: %w", err)
		}
		consensusClientSettings.Name = client
	}
	config.ConsensusCientSettings = consensusClientSettings

	// Execution Client
	executionClientSettings := defaultExecutionClientSettings
	if options.ExecutionClient != "" {
		client, err := utils.GetExecutionClient(options.ExecutionClient)
		if err != nil {
			return nil, fmt.Errorf("invalid execution client: %w", err)
		}
		executionClientSettings.Name = client
	}
	config.ExecutionCientSettings = executionClientSettings

	if starknetNode {
		config.IsValidatorNode = validator
		config.JunoConfig = defaultJunoConfig
	}

	if validator {
		if walletConfig != nil {
			config.Wallet = *walletConfig
		}
		config.ValidatorConfig = types.ValidatorConfig{
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

	return config, nil
}

// stakeForValidator handles staking STARK for a validator node.
func stakeForValidator(network string) error {
	fmt.Println(utils.Cyan("💰 Staking STARK for validator..."))

	var wallet types.WalletConfig
	walletsBytes, err := yaml.Marshal(options.Config.Wallet)
	if err != nil {
		return fmt.Errorf("could not marshal wallet config for staking: %w", err)
	}
	if err := envsubt.Unmarshal(walletsBytes, &wallet); err != nil {
		return fmt.Errorf("could not substitute env vars in wallet config for staking: %w", err)
	}
	rpcProvider, err := utils.CreateRPCProvider(network)
	if err != nil {
		return fmt.Errorf("❌ Error creating RPC provider: %v\n", err)
	}

	if err := validator.StakeStark(network, rpcProvider, wallet); err != nil {
		return fmt.Errorf("error staking STARK: %w", err)
	}

	fmt.Println(utils.Green("✅ Staking successful!"))
	return nil
}

// installClients installs the necessary clients based on the configuration.
func installClients(starknetNode, validator bool, consensusClient, executionClient types.ClientType) {
	fmt.Println(utils.Cyan("🚀 Installing clients..."))
	clients := []types.ClientType{consensusClient, executionClient}

	for _, client := range clients {
		err := options.Installer.InstallClient(client)
		if errors.Is(err, pkg.ErrClientIsInstalled) {
			fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client %s is already installed. Skipping.", client)))
			continue
		}
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Could not install client %s: %v", client, err)))
			return
		}
		fmt.Println(utils.Green(fmt.Sprintf("✅ Client %s installed successfully.", client)))
	}

	if starknetNode {
		err := options.Installer.InstallClient(types.ClientJuno)
		if errors.Is(err, pkg.ErrClientIsInstalled) {
			fmt.Println(utils.Yellow("🤔 Client Juno is already installed. Skipping."))
		} else if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Could not install client Juno: %v", err)))
			return
		}
		fmt.Println(utils.Green("✅ Client Juno installed successfully."))

		if validator {
			err := options.Installer.InstallClient(types.ClientStarkValidator)
			if err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Could not install client Starknet Validator: %v", err)))
				return
			}
			fmt.Println(utils.Green("✅ Client Starknet Validator installed successfully."))
		}
	}
}

func runNewConfigCommand(cmd *cobra.Command, args []string) {
	network, _ := cmd.Flags().GetString("network")
	starknetNode, _ := cmd.Flags().GetBool("starknet-node")
	validator, _ := cmd.Flags().GetBool("validator")
	install, _ := cmd.Flags().GetBool("install")

	var walletConfig *types.WalletConfig
	var err error

	if validator {
		walletConfig, err = handleValidatorWalletSetup(network)
		if err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Error setting up validator wallet: %v", err)))
			return
		}
	}

	config, err := generateConfig(network, starknetNode, validator, walletConfig)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error generating config: %v", err)))
		return
	}

	fmt.Println(utils.Cyan("📝 Creating starknode.yaml configuration file..."))
	err = utils.CreateStarkNodeConfig(config)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating config file: %v", err)))
	}
	fmt.Println(utils.Green("✅ Configuration file 'starknode.yaml' created successfully."))

	options.Config, err = utils.LoadConfig()
	if err != nil {
		fmt.Println(utils.Red("❌ Could not load config"))
		return
	}

	if validator {
		if err := stakeForValidator(network); err != nil {
			fmt.Println(utils.Red(fmt.Sprintf("❌ Error staking for validator: %v", err)))
			return
		}
	}

	if install {
		installClients(starknetNode, validator, config.ConsensusCientSettings.Name, config.ExecutionCientSettings.Name)
	}
}

func init() {
	options.InitGlobalOptions(newConfigCommand)
	newConfigCommand.Flags().String("network", "sepolia", "Select the network to connect to (e.g., 'mainnet', 'sepolia')")
	newConfigCommand.Flags().Bool("starknet-node", false, "Install a Starknet node")
	newConfigCommand.Flags().Bool("validator", false, "Configure a validator node (deploys account and sets up wallet config)")
	newConfigCommand.Flags().BoolP("install", "i", true, "Install clients automatically after setup")
}
