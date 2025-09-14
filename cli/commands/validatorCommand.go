package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/validator"
)

// ANSI color codes
const (
	ColorGreen = "\033[32m"
	ColorAqua  = "\033[36m"
	ColorWhite = "\033[97m"
	ColorReset = "\033[0m"
)

var ValidatorCommand = &cobra.Command{
	Use:   "validator",
	Short: "Get validator information",
	Long:  `Displays information about the validator associated with the configured wallet.`,
	Run:   validatorCommand,
}

func validatorCommand(cmd *cobra.Command, args []string) {
	config, err := utils.LoadConfig()
	if err != nil {
		fmt.Printf("❌ Error loading config: %v", err)
		return
	}

	if !config.IsValidatorNode {
		fmt.Printf("❌ Error No validator node configured")
		return
	}

	rpcProvider, err := utils.CreateRPCProvider(config.Network)
	if err != nil {
		fmt.Printf("❌ Error creating RPC provider: %v", err)
		return
	}

	validatorInfo, err := validator.GetValidatorInfo(*rpcProvider, config.Wallet.Wallet)
	if err != nil {
		fmt.Printf("❌ Error getting validator info: %v", err)
		return
	}

	fmt.Printf("%s✅ Validator Information ✅%s\n\n", ColorGreen, ColorReset)

	fmt.Printf("%s%s%s\n", ColorAqua, "Reward Address:", ColorReset)
	fmt.Printf("%s%s%s\n\n", ColorWhite, validatorInfo.RewardAddress, ColorReset)

	fmt.Printf("%s%s%s\n", ColorAqua, "Operational Address:", ColorReset)
	fmt.Printf("%s%s%s\n\n", ColorWhite, validatorInfo.OperationalAddress, ColorReset)

	fmt.Printf("%s%s%s\n", ColorAqua, "Total Staked:", ColorReset)
	fmt.Printf("%s%.4f STRK%s\n\n", ColorWhite, validatorInfo.TotalStaked, ColorReset)

	fmt.Printf("%s%s%s\n", ColorAqua, "Unclaimed Rewards:", ColorReset)
	fmt.Printf("%s%.4f STRK%s\n", ColorWhite, validatorInfo.UnclaimedRewards, ColorReset)
}
