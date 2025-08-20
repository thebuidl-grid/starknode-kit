package constants

import (
	"fmt"
	"os"
	"path"

	"github.com/common-nighthawk/go-figure"
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
)

func getHomeDir() string {
	homeDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}