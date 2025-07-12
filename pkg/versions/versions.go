package versions

const (
	StarkNodeVersion        = "0.0.1"
	LatestGethVersion       = "1.15.10"
	LatestRethVersion       = "1.3.4"
	LatestLighthouseVersion = "7.0.1"
	LatestPrysmVersion      = "v5.1.4" // Updated to match current Prysm releases
	LatestJunoVersion       = "v0.14.6"
)

// ClientReleaseUrls maps client names to their GitHub release URLs
var ClientReleaseUrls = map[string]string{
	"geth":       "https://github.com/ethereum/go-ethereum/releases",
	"reth":       "https://github.com/paradigmxyz/reth/releases",
	"lighthouse": "https://github.com/sigp/lighthouse/releases",
	"prysm":      "https://github.com/prysmaticlabs/prysm/releases",
	"juno":       "https://github.com/NethermindEth/juno/releases",
}

// GetClientLatestVersion returns the static latest version for a client
func GetClientLatestVersion(client string) string {
	switch client {
	case "geth":
		return LatestGethVersion
	case "reth":
		return LatestRethVersion
	case "lighthouse":
		return LatestLighthouseVersion
	case "prysm":
		return LatestPrysmVersion
	case "juno":
		return LatestJunoVersion
	default:
		return "unknown"
	}
}
