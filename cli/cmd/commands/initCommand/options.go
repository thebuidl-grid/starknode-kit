package initcommand

import "starknode-kit/pkg/types"

var (
	elClientOptions = []types.ClientType{
		types.ClientGeth, types.ClientReth,
	}
	clClientOptions = []types.ClientType{
		types.ClientPrysm, types.ClientLighthouse,
	}

	supportedNetorks = []string{
		"mainnet", "sepolia",
	}

	clClientPort = []int{5052, 9000}
	elClientPort = []int{30303}
)

// Helper functions
func clientTypesToStrings(clients []types.ClientType) []string {
	result := make([]string, len(clients))
	for i, client := range clients {
		result[i] = client.String()
	}
	return result
}
