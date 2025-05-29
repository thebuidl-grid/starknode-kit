package pkg

import (
	"os"
	"path"
)

var (
	InstallDir = path.Join(getHomeDir(), "starcknode-kit")

	InstallClientsDir = path.Join(InstallDir, "ethereum_clients")

	jwtDir  = path.Join(InstallDir, "ethereum_clients", "jwt")
	JWTPath = path.Join(jwtDir, "jwt.hex")
)

func getHomeDir() string {
	homeDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
