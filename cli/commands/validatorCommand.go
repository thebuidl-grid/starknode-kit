package commands

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/clients"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/validator"
	"github.com/thebuidl-grid/starknode-kit/pkg/versions"
)

var rpcProvider *rpc.Provider

var ValidatorCommand = &cobra.Command{
	Use:   "validator",
	Short: "Manage validator",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().NFlag() == 0 {
			cmd.Help()
			return
		}
		versionFlag, _ := cmd.Flags().GetBool("version")
		rpc_url, _ := cmd.Flags().GetString("rpc")
		if versionFlag {
			if !utils.IsInstalled(types.ClientStarkValidator) {
				fmt.Println(utils.Yellow(fmt.Sprintf("🤔 Client %s is not installed.", types.ClientStarkValidator)))
				return
			}
			version := versions.GetVersionNumber(string(types.ClientStarkValidator))

			fmt.Printf("%s version: %s\n", clientName, utils.Green(version))
			return
		}
		if rpc_url != "" {
			var err error
			url, err := url.ParseRequestURI(rpc_url)
			if err != nil {
				fmt.Printf("invalid URL format for rpc_url: '%s'", utils.Red(rpc_url))
				return
			}
			switch url.Scheme {
			case "https", "http":
				options.Config.ValidatorConfig.ProviderConfig.JunoRPC = rpc_url
				fmt.Printf("Successfully set Provider HTTP url to '%s'\n", utils.Green(rpc_url))
			case "wss", "ws":
				options.Config.ValidatorConfig.ProviderConfig.JunoWS = rpc_url
				fmt.Printf("Successfully set Provider WS url to '%s'\n", utils.Green(rpc_url))
			default:
				err = fmt.Errorf("Inalid scheme '%s'", utils.Red(url.Scheme))
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := utils.UpdateStarkNodeConfig(options.Config); err != nil {
				fmt.Println(utils.Red(fmt.Sprintf("❌ Failed to save config: %v", err)))
			}
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
		}
		if !options.Config.IsValidatorNode {
			fmt.Println(utils.Red("❌ This is not a validator node. Check your configuration."))
			os.Exit(1)
		}
		var err error
		rpcProvider, err = utils.CreateRPCProvider(options.Config.Network)
		if err != nil {
			fmt.Printf(utils.Red("❌ Error creating RPC provider: %v\n"), err)
			os.Exit(1)
		}
	},
}

var validatorInfoCommand = &cobra.Command{
	Use:   "info",
	Short: "Get validator information",
	Long:  `Displays information about the validator associated with the configured wallet.`,
	Run:   validatorInfoCommandRun,
}

var validatorStopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stop the validator client",
	Long:  `Stops the running starknet validator client process.`,
	Run:   validatorStopCommandRun,
}

var validatorStartCommand = &cobra.Command{
	Use:   "start",
	Short: "start the validator client",
	Long:  `Starts the running starknet validator client process.`,
	Run:   validatorStartCommandRun,
}

var validatorStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Check validator client status",
	Run:   validatorStatusCommandRun,
}

func validatorInfoCommandRun(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}

	validatorInfo, err := validator.GetValidatorInfo(rpcProvider, options.Config.Wallet.Wallet)
	if err != nil {
		fmt.Printf(utils.Red("❌ Error getting validator info: %v\n"), err)
		return
	}

	fmt.Printf("%s\n\n", utils.Green("✅ Validator Information ✅"))

	utils.PrintKV("Reward Address", validatorInfo.RewardAddress)
	utils.PrintKV("Operational Address", validatorInfo.OperationalAddress)
	utils.PrintKV("Total Staked", fmt.Sprintf("%.4f STRK", validatorInfo.TotalStaked))
	utils.PrintKV("Unclaimed Rewards", fmt.Sprintf("%.4f STRK", validatorInfo.UnclaimedRewards))
}

func validatorStopCommandRun(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}
	processInfo := process.GetProcessInfo(string(types.ClientStarkValidator))
	if processInfo == nil {
		fmt.Println(utils.Yellow("Validator client is not running."))
		return
	}
	err := process.StopClient(processInfo.PID)
	if err != nil {
		fmt.Printf(utils.Red("Could not stop validator process: %v\n"), err)
		return
	}
	fmt.Println(utils.Green("✅ Validator client stopped successfully."))
}

func validatorStartCommandRun(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}

	fmt.Println(utils.Cyan("🚀 Starting Validator client..."))
	validatorNode, err := clients.NewValidatorClient(options.Config.ValidatorConfig)
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error creating validator client: %v", err)))
		return
	}
	err = validatorNode.Start()
	if err != nil {
		fmt.Println(utils.Red(fmt.Sprintf("❌ Error starting validator client: %v", err)))
		return
	}
	fmt.Println(utils.Cyan("✅ Validator started"))
	fmt.Println(utils.Cyan("⏳ Waiting for log files to be created..."))
	options.LoadLogs([]string{string(types.ClientStarkValidator)})
}

var validatorBalanceCommand = &cobra.Command{
	Use:   "balance",
	Short: "Get validator balance",
	Long:  `Get the STRK balance of the validator wallet.`,
	Run:   validatorBalanceCommandRun,
}

func validatorBalanceCommandRun(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}

	balance, err := validator.GetValidatorBalance(rpcProvider, options.Config.Wallet.Wallet)
	if err != nil {
		fmt.Printf(utils.Red("❌ Error getting validator balance: %v"), err)
		return
	}

	fmt.Printf("%s %.4f STRK\n", utils.Green("✅ Validator Balance:"), balance)
}

func validatorStatusCommandRun(cmd *cobra.Command, args []string) {
	clientName := string(types.ClientStarkValidator)
	processInfo := process.GetProcessInfo(clientName)
	if processInfo != nil {
		fmt.Printf("Client: %s\n", utils.Blue(processInfo.Name))
		fmt.Printf("  Status: %s (PID: %d)\n", utils.Green("Running"), processInfo.PID)
		fmt.Printf("  Uptime: %s\n", utils.Green(processInfo.Uptime.Round(time.Second).String()))
	} else {
		fmt.Printf("  Status: %s\n", utils.Red("Stopped"))
	}

	junoMetrics := utils.GetJunoMetrics(options.Config.Network)
	fmt.Printf("\nJuno Node Status:\n")
	if junoMetrics.IsSyncing {
		fmt.Printf("  Sync Status: %s\n", utils.Yellow("Syncing"))
		fmt.Printf("  Sync Percent: %s\n", utils.Yellow(fmt.Sprintf("%.2f%%", junoMetrics.SyncPercent)))
	} else {
		fmt.Printf("  Sync Status: %s\n", utils.Green("Synced"))
	}
	fmt.Printf("  Current Block: %s\n", utils.Green(fmt.Sprintf("%d", junoMetrics.CurrentBlock)))
	fmt.Printf("  Network: %s\n", utils.Green(junoMetrics.NetworkName))
}

func init() {
	ValidatorCommand.Flags().BoolP("version", "v", false, "Get validator version")
	ValidatorCommand.Flags().String("rpc", "", "Set juno RPC endpoint")
	ValidatorCommand.AddCommand(validatorInfoCommand)
	ValidatorCommand.AddCommand(validatorStatusCommand)
	ValidatorCommand.AddCommand(validatorStopCommand)
	ValidatorCommand.AddCommand(validatorStartCommand)
	ValidatorCommand.AddCommand(validatorBalanceCommand)
}
