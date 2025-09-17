package validator

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	starkutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

// newAccount creates a new Starknet account instance.
func newAccount(wallet types.Wallet, rpcProvider *rpc.Provider) (*account.Account, error) {
	userWalletAddress, err := starkutils.HexToFelt(wallet.Address)
	if err != nil {
		return nil, err
	}

	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(wallet.PrivateKey, 0)
	if !ok {
		return nil, fmt.Errorf("failed to convert private key to big.Int")
	}
	ks.Put(wallet.PublicKey, privKeyBI)

	return account.NewAccount(
		rpcProvider,
		userWalletAddress,
		wallet.PublicKey,
		ks,
		account.CairoV2,
	)
}

// approveStakes checks allowance and approves the staking contract to spend STRK if necessary.
func approveStakes(network string, accnt *account.Account, rpcProvider *rpc.Provider, balance *felt.Felt) error {
	stakingAddr, err := starkutils.HexToFelt(constants.StakingContract)
	if err != nil {
		return err
	}

	allowance, err := getAllowance(rpcProvider, accnt.Address, stakingAddr)
	if err != nil {
		return err
	}

	if starkutils.FeltToBigInt(allowance).Cmp(starkutils.FeltToBigInt(constants.Stakes[network][0])) >= 0 {
		return nil
	}

	starkAddr, err := starkutils.HexToFelt(constants.StrkTokenAddress)
	if err != nil {
		return err
	}

	approveTxn := rpc.FunctionCall{
		ContractAddress:    starkAddr,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("approve"),
		Calldata:           []*felt.Felt{stakingAddr, constants.Stakes[network][0], constants.Stakes[network][1]},
	}

	invokeTxn, estimatedFee, err := estimateFee(accnt, []rpc.FunctionCall{approveTxn})
	if err != nil {
		return err
	}

	stakeAmount := constants.Stakes[network][0]
	requiredAmount := new(felt.Felt).Add(stakeAmount, estimatedFee)
	if balance.Cmp(requiredAmount) < 0 {
		needed := starkutils.FRIToSTRK(requiredAmount)
		have := starkutils.FRIToSTRK(balance)
		return fmt.Errorf("insufficient balance to approve. Have: %.6f STRK, Need: %.6f STRK", have, needed)
	}

	return executeTxn(accnt, invokeTxn)
}

// stakeAndSetCommission builds and executes the stake and set_commission transactions.
func stakeAndSetCommission(network string, accnt *account.Account, wallet types.WalletConfig, balance *felt.Felt) error {
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

	commissionFelt := starkutils.BigIntToFelt(big.NewInt(int64(uint16(commissionInt * 100))))

	setCommissionTxn := rpc.FunctionCall{
		ContractAddress:    stackingAddr,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("set_commission"),
		Calldata:           []*felt.Felt{commissionFelt},
	}

	invokeTxn, estimatedFee, err := estimateFee(accnt, []rpc.FunctionCall{stakeTxn, setCommissionTxn})
	if err != nil {
		return err
	}

	stakeAmount := constants.Stakes[network][0]
	requiredAmount := new(felt.Felt).Add(stakeAmount, estimatedFee)
	if balance.Cmp(requiredAmount) < 0 {
		needed := starkutils.FRIToSTRK(requiredAmount)
		have := starkutils.FRIToSTRK(balance)
		return fmt.Errorf("insufficient balance to stake. Have: %.6f STRK, Need: %.6f STRK", have, needed)
	}

	return executeTxn(accnt, invokeTxn)
}

// getAllowance retrieves the allowance for a spender from an owner.
func getAllowance(rpcProvider *rpc.Provider, owner, spender *felt.Felt) (*felt.Felt, error) {
	contractAddress, err := starkutils.HexToFelt(constants.StrkTokenAddress)
	if err != nil {
		return nil, err
	}

	txn := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: starkutils.GetSelectorFromNameFelt("allowance"),
		Calldata:           []*felt.Felt{owner, spender},
	}

	result, err := rpcProvider.Call(context.Background(), txn, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

// estimateFee estimates the fee for a set of transactions.
func estimateFee(accnt *account.Account, calls []rpc.FunctionCall) (*rpc.BroadcastInvokeTxnV3, *felt.Felt, error) {
	fmt.Println("Estimating transaction fees...")
	invokeTxn, feesEstimate, err := utils.EstimateGasFee(accnt, calls)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to estimate gas fee: %w", err)
	}

	estimatedFee := feesEstimate[0].OverallFee
	estimatedFeeStark := starkutils.FRIToSTRK(estimatedFee)
	fmt.Printf("Estimated fee: %.6f STRK. Adjusting for a higher success rate...\n", estimatedFeeStark)

	invokeTxn.ResourceBounds = starkutils.FeeEstToResBoundsMap(feesEstimate[0], 1.5)

	if err := accnt.SignInvokeTransaction(context.Background(), invokeTxn); err != nil {
		return nil, nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return invokeTxn, estimatedFee, nil
}

// executeTxn sends a transaction and waits for its confirmation.
func executeTxn(accnt *account.Account, invokeTxn *rpc.BroadcastInvokeTxnV3) error {
	fmt.Println("Sending transactions...")
	resp, err := accnt.SendTransaction(context.Background(), invokeTxn)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Transaction successfully submitted! Transaction hash: %s\n", utils.FormatTransactionHash(resp.Hash))
	fmt.Println("Waiting for transaction confirmation...")

	receipt, err := accnt.WaitForTransactionReceipt(context.Background(), resp.Hash, 15*time.Second)
	if err != nil {
		transactionURL := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", utils.FormatTransactionHash(resp.Hash))
		fmt.Printf("Transaction failed or timed out. View details here: %s\n", transactionURL)
		return fmt.Errorf("error waiting for transaction receipt: %w", err)
	}

	transactionURL := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", utils.FormatTransactionHash(receipt.Hash))
	fmt.Printf("Transaction successful! View details here: %s\n", transactionURL)
	return nil
}
