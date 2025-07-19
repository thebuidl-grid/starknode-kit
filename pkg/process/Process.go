package process

import (
	"io"
	"os/exec"
	"syscall"

	t "github.com/thebuidl-grid/starknode-kit/pkg/types"
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

func StopClient(pid int) error {
	return stopProcess(pid)
}

func GetProcessInfo(p string) *t.ProcessInfo {
	return getProcessInfo(p)
}
