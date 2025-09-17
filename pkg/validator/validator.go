package validator

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	starkutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

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
	return types.ValidatorInfo{
		RewardAddress:      result[1].String(),
		OperationalAddress: result[2].String(),
		TotalStaked:        starkutils.FRIToSTRK(result[4]),
		UnclaimedRewards:   starkutils.FRIToSTRK(result[5]),
	}, nil
}

func GetValidatorBalance(rpcProvider *rpc.Provider, wallet types.Wallet) (float64, error) {
	accnt, err := newAccount(wallet, rpcProvider)
	if err != nil {
		return 0, err
	}
	balance, err := utils.CheckBalance(rpcProvider, accnt.Address)
	if err != nil {
		return 0, err
	}
	balanceStark := starkutils.FRIToSTRK(balance)
	return balanceStark, nil
}

// StakeStark performs the staking process for a given wallet configuration and network.
// It checks for sufficient balance, builds the necessary transactions (approve, stake, set_commission),
// sends them as a multi-call, and waits for the transaction to be confirmed.
func StakeStark(network string, rpcProvider *rpc.Provider, wallet types.WalletConfig) error {
	accnt, err := newAccount(wallet.Wallet, rpcProvider)
	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	// Check for sufficient balance before proceeding
	balance, err := utils.CheckBalance(rpcProvider, accnt.Address)
	if err != nil {
		return fmt.Errorf("failed to check balance: %w", err)
	}

	err = approveStakes(network, accnt, rpcProvider, balance)
	if err != nil {
		return err
	}

	// Build transactions
	stackingAddr, err := starkutils.HexToFelt(constants.StakingContract)
	if err != nil {
		return fmt.Errorf("failed to convert staking contract address to felt: %w", err)
	}

	rewardAddress, err := starkutils.HexToFelt(wallet.RewardAddress)
	if err != nil {
		return fmt.Errorf("failed to convert reward address to felt: %w", err)
	}

	stakeTxn := rpc.FunctionCall{
		ContractAddress:    stackingAddr,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("stake"),
		Calldata:           []*felt.Felt{rewardAddress, accnt.Address, constants.Stakes[network][0]},
	}

	commissionInt, err := strconv.Atoi(wallet.StakeCommision)
	if err != nil {
		return fmt.Errorf("failed to convert stake commission to integer: %w", err)
	}

	// Commission is in basis points (e.g., 10% = 1000)
	commissionFelt := starkutils.BigIntToFelt(
		big.NewInt(int64(uint16(commissionInt * 100))),
	)

	setCommissionTxn := rpc.FunctionCall{
		ContractAddress:    stackingAddr,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("set_commission"),
		Calldata:           []*felt.Felt{commissionFelt},
	}

	invokeTxn, err := estimateFee(accnt, []rpc.FunctionCall{stakeTxn, setCommissionTxn}, balance)
	if err != nil {
		return err
	}

	return executeTxn(accnt, invokeTxn)
}

func approveStakes(network string, accnt *account.Account, rpcProvider *rpc.Provider, balance *felt.Felt) error {
	stakingAddr, err := starkutils.HexToFelt(constants.StakingContract)
	if err != nil {
		return err
	}
	starkAddr, err := starkutils.HexToFelt(constants.StrkTokenAddress)
	if err != nil {
		return err
	}
	allowance, err := getAllowance(rpcProvider, accnt.Address, stakingAddr)
	if err != nil {
		return err
	}
	allowanceFelt := starkutils.FeltToBigInt(allowance)
	if allowanceFelt.Cmp(starkutils.FeltToBigInt(constants.Stakes[network][0])) >= 0 {
		return nil
	}
	approveTxn := rpc.FunctionCall{
		ContractAddress:    starkAddr,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("approve"),
		Calldata:           []*felt.Felt{stakingAddr, constants.Stakes[network][0], constants.Stakes[network][1]},
	}
	invokeTxn, err := estimateFee(accnt, []rpc.FunctionCall{approveTxn}, balance)
	if err != nil {
		return err
	}

	return executeTxn(accnt, invokeTxn)
}