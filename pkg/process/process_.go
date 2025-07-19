package process

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
)

func stopProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 3)
	if running := IsProcessRunning(pid); running {
		err = p.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}
	return nil
}

func getProcessInfo(processName string) *t.ProcessInfo {
	// Read /proc to find the process
	procDirs, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return nil
	}

	for _, procDir := range procDirs {
		// Read cmdline to get the command
		cmdlineFile := filepath.Join(procDir, "cmdline")
		cmdlineBytes, err := os.ReadFile(cmdlineFile)
		if err != nil {
			continue
		}

		cmdline := string(cmdlineBytes)
		if strings.Contains(cmdline, processName) {
			// Extract PID from directory name
			pidStr := filepath.Base(procDir)
			pid, err := strconv.Atoi(pidStr)
			if err != nil {
				continue
			}

			// Get process status
			statusFile := filepath.Join(procDir, "stat")
			statusBytes, err := os.ReadFile(statusFile)
			if err != nil {
				continue
			}

			// Parse stat file for uptime
			statFields := strings.Fields(string(statusBytes))
			if len(statFields) > 21 {
				startTimeJiffies, err := strconv.ParseUint(statFields[21], 10, 64)
				if err == nil {
					// Calculate uptime (simplified)
					uptimeSeconds := time.Now().Unix() - int64(startTimeJiffies/100)
					uptime := time.Duration(uptimeSeconds) * time.Second

					return &t.ProcessInfo{
						PID:    pid,
						Name:   processName,
						Status: "running",
						Uptime: uptime,
					}
				}
			}
		}
	}

	return nil
}
