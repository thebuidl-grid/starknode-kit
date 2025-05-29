package pkg

import (
	"os"
	"path/filepath"
)

var (
  InstallDir = getHomeDir()
	JWTPath = filepath.Join(InstallDir, "ethereum_clients", "jwt", "jwt.hex")
)

func getHomeDir() string{
  homeDir, err := os.UserConfigDir()
  if err != nil{
    panic(err)
  }
  return homeDir
}
