package constants

import (
	"fmt"
	"os"
	"path"

	"github.com/NethermindEth/juno/core/felt"
	starkutils "github.com/NethermindEth/starknet.go/utils"
	"github.com/common-nighthawk/go-figure"
)

const (
	mainnetStake = "20000000000000000000000"
	testnetStake = "1000000000000000000"

	PredeployedClassHash = "0x61dac032f228abef9c6626f995015233097ae253a7f72d68552db02f2971b8f"
	StrkTokenAddress     = "0x04718f5a0fc34cc1af16a1cdee98ffb20c31f5cd61d6ab07201858f4287c938d"
	StakingContract      = "0x03745ab04a431fc02871a139be6b93d9260b0ff3e779ad9c8b377183b23109f1"
)

var (
	InstallDir = path.Join(getHomeDir(), "starknode-kit")

	InstallClientsDir  = path.Join(InstallDir, "ethereum_clients")
	InstallStarknetDir = path.Join(InstallDir, "starknet_clients")

	JwtDir      = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath     = path.Join(JwtDir, "jwt.hex")
	ConfigDir   = path.Join(InstallDir, "config")
	ConfigPath  = fmt.Sprintf("%s/starknode.yaml", ConfigDir)
	EnvFIlePath = fmt.Sprintf("%s/.starknode.env", ConfigDir)
	Banner      = figure.NewColorFigure("Starknode kit", "slant", "green", true)

	RPCURL = map[string]string{
		"mainnet": "https://starknet-mainnet.public.blastapi.io/rpc/v0_9",
		"sepolia": "https://starknet-sepolia.public.blastapi.io/rpc/v0_9",
	}

	mainnetBig, _ = starkutils.HexToU256Felt(starkutils.StrToHex(mainnetStake))
	sepoliaBig, _ = starkutils.HexToU256Felt(starkutils.StrToHex(testnetStake))

	Stakes = map[string][]*felt.Felt{
		"mainnet": mainnetBig,
		"sepolia": sepoliaBig,
	}
)

func getHomeDir() string {
	homeDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}

