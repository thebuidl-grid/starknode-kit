package versions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
)

var (
	StarkNodeVersion string
)

// ClientReleaseUrls maps client names to their GitHub release URLs
var ClientReleaseUrls = map[string]string{
	"geth":                "https://github.com/ethereum/go-ethereum/releases",
	"reth":                "https://github.com/paradigmxyz/reth/releases",
	"lighthouse":          "https://github.com/sigp/lighthouse/releases",
	"prysm":               "https://github.com/prysmaticlabs/prysm/releases",
	"juno":                "https://github.com/NethermindEth/juno/releases",
	"starknet-staking-v2": "https://github.com/NethermindEth/starknet-staking-v2/releases",
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
	case "starknet-staking-v2":
		return FetchLatestStarknetValidatorVersion()
	default:
		return "", fmt.Errorf("unsupported client: %s", client)
	}
}

func GetVersionNumber(client string) string {

	var argument string

	switch client {
	case "juno":
		path := filepath.Join(constants.InstallStarknetDir, "juno", ".version")
		version, _ := os.ReadFile(path)
		versionMatch := regexp.MustCompile(`juno version (\d+\.\d+\.\d+)`).FindStringSubmatch(string(version))
		if len(versionMatch) > 1 {
			return versionMatch[1]
		}
		return ""
	case "reth", "lighthouse", "geth", "starknet-staking-v2":
		argument = "--version"
	case "prysm":
		argument = "beacon-chain --version"
	default:
		fmt.Printf("Unknown client: %s\n", client)
		return ""
	}

	var clientCommand string
	switch runtime.GOOS {
	case "darwin", "linux":
		if client == "prysm" {
			clientCommand = filepath.Join(constants.InstallClientsDir, client, fmt.Sprintf("%s.sh", client))
		} else {
			clientCommand = filepath.Join(constants.InstallClientsDir, client, client)
		}
	case "windows":
		fmt.Println("getVersionNumber() for windows is not yet implemented")
		os.Exit(1)
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		return ""
	}

	if client == "starknet-staking-v2" {
		clientCommand = filepath.Join(constants.InstallStarknetDir, client, "validator")
	}

	cmdParts := strings.Split(argument, " ")
	cmd := exec.Command(clientCommand, cmdParts...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command for %s: %v\n", client, err)
		return ""
	}

	versionOutput := strings.TrimSpace(string(output))
	var versionMatch []string

	switch client {
	case "reth":
		versionMatch = regexp.MustCompile(`reth-ethereum-cli Version: (\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	case "lighthouse":
		versionMatch = regexp.MustCompile(`Lighthouse v(\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	case "geth":
		versionMatch = regexp.MustCompile(`geth version (\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	case "prysm":
		versionMatch = regexp.MustCompile(`beacon-chain-v(\d+\.\d+\.\d+)-`).FindStringSubmatch(versionOutput)
	case "starknet-staking-v2":
		versionMatch = regexp.MustCompile(`validator version (\d+\.\d+\.\d+)`).FindStringSubmatch(versionOutput)
	}

	if len(versionMatch) > 1 {
		return versionMatch[1]
	}

	fmt.Printf("Unable to parse version number for %s\n", client)
	return ""
}
