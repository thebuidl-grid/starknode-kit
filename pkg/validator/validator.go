package validator

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	starkutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

func GetValidatorInfo(rpcProvider *rpc.Provider, wallet types.Wallet) (types.ValidatorInfo, error) {
	accnt, err := newAccount(wallet, rpcProvider)
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
		Calldata:           []*felt.Felt{accnt.Address},
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
