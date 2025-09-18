package validator

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	starkutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

// GetValidatorInfo retrieves validator information from the staking contract.
func GetValidatorInfo(rpcProvider *rpc.Provider, wallet types.Wallet) (types.ValidatorInfo, error) {
	address, err := starkutils.HexToFelt(wallet.Address)
	if err != nil {
		return types.ValidatorInfo{}, err
	}

	contractAddress, err := starkutils.HexToFelt(constants.StakingContract)
	if err != nil {
		return types.ValidatorInfo{}, err
	}

	txn := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("get_staker_info_v1"),
		Calldata:           []*felt.Felt{address},
	}

	result, err := rpcProvider.Call(context.Background(), txn, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return types.ValidatorInfo{}, err
	}

	if len(result) == 1 {
		return types.ValidatorInfo{}, fmt.Errorf("address not a validator")
	}

	return types.ValidatorInfo{
		RewardAddress:      result[1].String(),
		OperationalAddress: result[2].String(),
		TotalStaked:        starkutils.FRIToSTRK(result[4]),
		UnclaimedRewards:   starkutils.FRIToSTRK(result[5]),
	}, nil
}

// GetValidatorBalance retrieves the STRK balance for a given wallet.
func GetValidatorBalance(rpcProvider *rpc.Provider, wallet types.Wallet) (float64, error) {
	accnt, err := newAccount(wallet, rpcProvider)
	if err != nil {
		return 0, err
	}
	balance, err := utils.CheckBalance(rpcProvider, accnt.Address)
	if err != nil {
		return 0, err
	}
	return starkutils.FRIToSTRK(balance), nil
}

// StakeStark performs the staking process for a given wallet configuration and network.
func StakeStark(network string, rpcProvider *rpc.Provider, wallet types.WalletConfig) error {
	accnt, err := newAccount(wallet.Wallet, rpcProvider)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	_, err = GetValidatorInfo(rpcProvider, wallet.Wallet)
	if err == nil {
		fmt.Println(utils.Yellow("Address already a staker"))
		return nil
	}

	balance, err := utils.CheckBalance(rpcProvider, accnt.Address)
	if err != nil {
		return fmt.Errorf("failed to check balance: %w", err)
	}

	if err := approveStakes(network, accnt, rpcProvider, balance); err != nil {
		return err
	}

	return stakeAndSetCommission(network, accnt, wallet, balance)
}
