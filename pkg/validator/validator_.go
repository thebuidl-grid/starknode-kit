package validator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	starknet "github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

func newAccount(wallet types.Wallet, rpcProvider *rpc.Provider) (*account.Account, error) {
	userWalletAddress, err := starknet.HexToFelt(wallet.Address)
	if err != nil {
		return nil, err
	}
	ks := account.NewMemKeystore()
	privKeyBI, ok := new(big.Int).SetString(wallet.PrivateKey, 0)
	if !ok {
		return nil, fmt.Errorf("Fail to convert privKey to bitInt")
	}
	ks.Put(wallet.PublicKey, privKeyBI)

	userAccount, err := account.NewAccount(
		rpcProvider,
		userWalletAddress,
		wallet.PublicKey,
		ks,
		account.CairoV2,
	)

	if err != nil {
		return nil, err
	}
	return userAccount, nil
}

func getAllowance(rpcProvider *rpc.Provider, owner *felt.Felt, spender *felt.Felt) (*felt.Felt, error) {
	contractAddress, err := starknet.HexToFelt(constants.StrkTokenAddress)
	if err != nil {
		return nil, err
	}

	txn := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: starknet.GetSelectorFromNameFelt("allowance"),
		Calldata:           []*felt.Felt{owner, spender},
	}

	result, err := rpcProvider.Call(context.Background(), txn, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

func estimateFee(accnt *account.Account, calls []rpc.FunctionCall, balance *felt.Felt) (*rpc.BroadcastInvokeTxnV3, error) {
	fmt.Println("Estimating transaction fees...")
	invokeTxn, feesEstimate, err := utils.EstimateGasFee(accnt, calls)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas fee: %w", err)
	}
	estimatedFee := feesEstimate[0].OverallFee
	estimatedFeeStark := starknet.FRIToSTRK(estimatedFee)
	fmt.Printf("Estimated fee: %.6f STRK. Adjusting for a higher success rate...\n", estimatedFeeStark)

	if balance.Cmp(estimatedFee) < 0 {
		needed := starknet.FRIToSTRK(estimatedFee)
		have := starknet.FRIToSTRK(balance)
		return nil, fmt.Errorf("insufficient balance for transaction. Have: %.6f STRK, Need: %.6f STRK", have, needed)
	}

	invokeTxn.ResourceBounds = starknet.FeeEstToResBoundsMap(feesEstimate[0], 1.5)

	// Re-sign the transaction as resource bounds have changed
	err = accnt.SignInvokeTransaction(context.Background(), invokeTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return invokeTxn, nil
}

func executeTxn(accnt *account.Account, invokeTxn *rpc.BroadcastInvokeTxnV3) error {
	err = accnt.SignInvokeTransaction(context.Background(), invokeTxn)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send the transaction
	fmt.Println("Sending transactions...")
	resp, err := accnt.SendTransaction(context.Background(), invokeTxn)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	fmt.Printf("Transaction successfully submitted! Transaction hash: %s\n", utils.FormatTransactionHash(resp.Hash))
	fmt.Println("⏳ Waiting for transaction confirmation...")

	// Wait for receipt
	receipt, err := accnt.WaitForTransactionReceipt(context.Background(), resp.Hash, 15*time.Second)
	if err != nil {
		transactionUrl := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", utils.FormatTransactionHash(resp.Hash))
		fmt.Printf("Transaction failed or timed out. View details here: %s\n", transactionUrl)
		return fmt.Errorf("error waiting for transaction receipt: %w", err)
	}

	transactionUrl := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", utils.FormatTransactionHash(receipt.Hash))
	fmt.Printf("Transaction successful! View details here: %s\n", transactionUrl)
	return nil
}