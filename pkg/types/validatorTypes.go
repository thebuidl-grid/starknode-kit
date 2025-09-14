package types

type (
	ValidatorInfo struct {
		RewardAddress      string
		OperationalAddress string
		TotalStaked        float64
		UnclaimedRewards   float64
	}
)
