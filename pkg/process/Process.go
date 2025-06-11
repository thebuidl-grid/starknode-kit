package process

import (
	"fmt"
	"io"
	"os/exec"
	"syscall"

	"starknode-kit/pkg"
	"starknode-kit/pkg/monitoring"
	"starknode-kit/pkg/types"
)

func IsProcessRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err != nil
}

func StartClient(name, command string, logPath io.Writer, args ...string) error {

	cmd := exec.Command(command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	cmd.Stdout = logPath
	cmd.Stderr = logPath
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}


func StopClient(clientName string) error {
	processInfo := getProcessInfo(clientName)
	if processInfo == nil {
		return fmt.Errorf("client %s is not running", clientName)
	}

	return stopProcess(processInfo.PID)
}

func GetRunningClients() []types.ClientStatus {
	var clients []types.ClientStatus

	// Check for Geth
	if gethInfo := getProcessInfo("geth"); gethInfo != nil {
		status := types.ClientStatus{
			Name:       "Geth",
			Status:     gethInfo.Status,
			PID:        gethInfo.PID,
			Uptime:     gethInfo.Uptime,
			Version:    pkg.GetClientVersion("geth"),
			SyncStatus: monitoring.GetGethSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Reth
	if rethInfo := getProcessInfo("reth"); rethInfo != nil {
		status := types.ClientStatus{
			Name:       "Reth",
			Status:     rethInfo.Status,
			PID:        rethInfo.PID,
			Uptime:     rethInfo.Uptime,
			Version:    pkg.GetClientVersion("reth"),
			SyncStatus: monitoring.GetRethSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Lighthouse
	if lighthouseInfo := getProcessInfo("lighthouse"); lighthouseInfo != nil {
		status := types.ClientStatus{
			Name:       "Lighthouse",
			Status:     lighthouseInfo.Status,
			PID:        lighthouseInfo.PID,
			Uptime:     lighthouseInfo.Uptime,
			Version:    pkg.GetClientVersion("lighthouse"),
			SyncStatus: monitoring.GetLighthouseSyncStatus(),
		}
		clients = append(clients, status)
	}

	// Check for Prysm
	if prysmInfo := getProcessInfo("prysm"); prysmInfo != nil {
		status := types.ClientStatus{
			Name:       "Prysm",
			Status:     prysmInfo.Status,
			PID:        prysmInfo.PID,
			Uptime:     prysmInfo.Uptime,
			Version:    pkg.GetClientVersion("prysm"),
			SyncStatus: monitoring.GetPrysmSyncStatus(),
		}
		clients = append(clients, status)
	}

	return clients
}
