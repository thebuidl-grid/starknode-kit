package validator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	starknet "github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

func newAccount(wallet types.Wallet, rpcProvider rpc.Provider) (*account.Account, error) {
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
		&rpcProvider,
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

func execute(ctx context.Context, user_account *account.Account, txn []rpc.InvokeFunctionCall) (rpc.AddInvokeTransactionResponse, error) {

	resp, err := user_account.BuildAndSendInvokeTxn(ctx, txn, nil)
	if err != nil {
		return rpc.AddInvokeTransactionResponse{}, err
	}
	transactionUrl := fmt.Sprintf("https://sepolia.voyager.online/tx/%s", utils.FormatTransactionHash(resp.Hash))
	fmt.Println("Transaction successfull, view here: ", transactionUrl)
	return resp, nil
}
