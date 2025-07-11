package initcommand

import "starknode-kit/pkg/types"

var (
	elClientOptions = []types.ClientType{
		types.ClientGeth, types.ClientReth,
	}
	clClientOptions = []types.ClientType{
		types.ClientPrysm, types.ClientLighthouse,
	}
)
