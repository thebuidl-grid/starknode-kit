package process

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
)

func getSystemBootTime() (time.Time, error) {
	content, err := os.ReadFile("/proc/stat")
	if err != nil {
		return time.Time{}, err
	}

	lines := strings.SplitSeq(string(content), "\n")
	for line := range lines {
		if strings.HasPrefix(line, "btime") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				bootTimeUnix, err := strconv.ParseInt(fields[1], 10, 64)
				if err != nil {
					return time.Time{}, err
				}
				return time.Unix(bootTimeUnix, 0), nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("btime not found in /proc/stat")
}

func stopProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Signal(syscall.SIGTERM)
	if err == nil {
		return nil
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
					bootTime, err := getSystemBootTime()
					if err != nil {
						return nil
					}

					startTimeSeconds := float64(startTimeJiffies) / 100.0

					processStartTime := bootTime.Add(time.Duration(startTimeSeconds) * time.Second)

					// Calculate uptime
					uptime := time.Since(processStartTime)

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
