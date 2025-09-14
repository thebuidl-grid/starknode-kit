package utils

import (
	"fmt"

	"github.com/thebuidl-grid/starknode-kit/pkg/types"
)

// ResolveClientType checks if a client name is valid across all client types.
func ResolveClientType(clientName string) (types.ClientType, error) {
	if client, err := GetExecutionClient(clientName); err == nil {
		return client, nil
	}
	if client, err := GetConsensusClient(clientName); err == nil {
		return client, nil
	}
	if client, err := GetStarknetClient(clientName); err == nil {
		return client, nil
	}
	return "", fmt.Errorf("unknown client: %s", clientName)
}
