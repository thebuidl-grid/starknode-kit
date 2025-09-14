package commands

import (
	"fmt"
	"os"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/cli/options"
	"github.com/thebuidl-grid/starknode-kit/pkg/clients"
	"github.com/thebuidl-grid/starknode-kit/pkg/process"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/validator"
)

var rpcProvider *rpc.Provider

var ValidatorCommand = &cobra.Command{
	Use:   "validator",
	Short: "Manage validator",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Parent() != nil && cmd.Parent().PersistentPreRun != nil {
			cmd.Parent().PersistentPreRun(cmd.Parent(), args)
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

func validatorInfoCommandRun(cmd *cobra.Command, args []string) {
	if !options.LoadedConfig {
		fmt.Println(utils.Red("❌ Config not found. Please run `starknode-kit config new`"))
		return
	}
	if !options.Config.IsValidatorNode {
		fmt.Println(utils.Red("❌ This is not a validator node. Check your configuration."))
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
	if !options.Config.IsValidatorNode {
		fmt.Println(utils.Red("❌ This is not a validator node. Check your configuration."))
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
	if !options.Config.IsValidatorNode {
		fmt.Println(utils.Red("❌ This is not a validator node. Check your configuration."))
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
	if !options.Config.IsValidatorNode {
		fmt.Println(utils.Red("❌ This is not a validator node. Check your configuration."))
		return
	}

	balance, err := validator.GetValidatorBalance(rpcProvider, options.Config.Wallet.Wallet)
	if err != nil {
		fmt.Printf(utils.Red("❌ Error getting validator balance: %v"), err)
		return
	}

	fmt.Printf("%s %.4f STRK", utils.Green("✅ Validator Balance:"), balance)
}

func init() {
	ValidatorCommand.AddCommand(validatorInfoCommand)
	ValidatorCommand.AddCommand(validatorStopCommand)
	ValidatorCommand.AddCommand(validatorStartCommand)
	ValidatorCommand.AddCommand(validatorBalanceCommand)
}
