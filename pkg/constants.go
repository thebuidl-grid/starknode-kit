package pkg

import (
	"fmt"
	"os"
	"path"
)

var (
	InstallDir = path.Join(getHomeDir(), "starknode-kit")

	InstallClientsDir = path.Join(InstallDir, "ethereum_clients")

	JwtDir      = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath     = path.Join(JwtDir, "jwt.hex")
	ConfigDir   = path.Join(InstallDir, "config")
	ConfigPath  = fmt.Sprintf("%s/starknode.yaml", ConfigDir)
	EnvFIlePath = fmt.Sprintf("%s/.starknode.env", ConfigDir)
)

func getHomeDir() string {
	homeDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
