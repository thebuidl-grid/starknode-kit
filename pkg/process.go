package pkg

import (
	"fmt"
	"io"
	"os"
	"path"
	"syscall"
	"time"
)

func IsRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err != nil
}
func StartProcess(name, command string, logPath io.Writer, args ...string) error {

	cmd := execCommand(command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	cmd.Stdout = logPath
	cmd.Stderr = logPath
	err := cmd.Start()
	if err != nil {
		return err
	}
	if err = writeToPIDFile(cmd.Process.Pid, name); err != nil {
		return fmt.Errorf("Failed to write PID file: %v\n", err)
	}
	return nil
}

func writeToPIDFile(pid int, name string) error {
	pidFile := path.Join(InstallDir, ".process", fmt.Sprintf("%s.pid", name))
	pidWrite := fmt.Sprintf("%d\n", pid)
	err := os.WriteFile(pidFile, []byte(pidWrite), 0644)
	if err != nil {
		return err
	}
	return nil

}

func StopProcess(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 3)
	if running := IsRunning(pid); running {
		err = p.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}
	return nil
}
