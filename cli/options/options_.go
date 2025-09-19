package options

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/thebuidl-grid/starknode-kit/pkg/constants"
	"github.com/thebuidl-grid/starknode-kit/pkg/types"
	"github.com/thebuidl-grid/starknode-kit/pkg/utils"
)

func getLatestLogFile(clientName string) (string, error) {
	var logDir string
	clientType, err := utils.ResolveClientType(clientName)
	if err != nil {
		return "", fmt.Errorf("could not resolve client type for %s: %w", clientName, err)
	}

	baseDir := constants.InstallClientsDir
	starknetClients := []types.ClientType{types.ClientJuno, types.ClientStarkValidator}
	if slices.Contains(starknetClients, clientType) {
		baseDir = constants.InstallStarknetDir
	}

	logDir = filepath.Join(baseDir, clientName, "logs")

	files, err := os.ReadDir(logDir)
	if err != nil {
		return "", fmt.Errorf("could not read log directory %s: %w", logDir, err)
	}

	var newestFile os.FileInfo
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileInfo, err := f.Info()
		if err != nil {
			continue
		}
		if newestFile == nil || fileInfo.ModTime().After(newestFile.ModTime()) {
			newestFile = fileInfo
		}
	}

	if newestFile == nil {
		return "", fmt.Errorf("no log files found in %s", logDir)
	}

	return filepath.Join(logDir, newestFile.Name()), nil
}
