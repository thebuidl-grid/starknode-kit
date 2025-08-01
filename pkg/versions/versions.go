package versions

import "fmt"

const (
	StarkNodeVersion        = "0.0.1"
	LatestGethVersion       = "1.15.10"
	LatestRethVersion       = "1.3.4"
	LatestLighthouseVersion = "7.0.1"
	LatestPrysmVersion      = "v5.1.4" // Updated to match current Prysm releases
	LatestJunoVersion       = "0.14.6"
)

// ClientReleaseUrls maps client names to their GitHub release URLs
var ClientReleaseUrls = map[string]string{
	"geth":       "https://github.com/ethereum/go-ethereum/releases",
	"reth":       "https://github.com/paradigmxyz/reth/releases",
	"lighthouse": "https://github.com/sigp/lighthouse/releases",
	"prysm":      "https://github.com/prysmaticlabs/prysm/releases",
	"juno":       "https://github.com/NethermindEth/juno/releases",
}

func FetchOnlineVersion(client string) (string, error) {
	switch client {
	case "geth":
		return FetchLatestGethVersion()
	case "reth":
		return FetchLatestRethVersion()
	case "lighthouse":
		return FetchLatestLighthouseVersion()
	case "prysm":
		return FetchLatestPrysmVersion()
	case "juno":
		return FetchLatestJunoVersion()
	default:
		return "", fmt.Errorf("unsupported client: %s", client)
	}
}
