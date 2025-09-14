package validator

import (
	"context"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
)

func GetValidatorInfo(rpcProvider rpc.Provider, wallet types.Wallet) (types.ValidatorInfo, error) {
	accnt, err := newAccount(wallet, rpcProvider)
	if err != nil {
		return types.ValidatorInfo{}, err
	}

	contractAddress, err := utils.HexToFelt(constants.StakingContract)
	if err != nil {
		return types.ValidatorInfo{}, err
	}

	txn := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt("get_staker_info_v1"),
		Calldata:           []*felt.Felt{accnt.Address},
	}

	result, err := rpcProvider.Call(context.Background(), txn, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return types.ValidatorInfo{}, err
	}
	return types.ValidatorInfo{
		RewardAddress:      result[1].String(),
		OperationalAddress: result[2].String(),
		TotalStaked:        utils.FRIToSTRK(result[4]),
		UnclaimedRewards:   utils.FRIToSTRK(result[5]),
	}, nil
}
