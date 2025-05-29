package pkg

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

func isRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err != nil
}
func StartProcess(command string, logPath io.Writer, args ...string) error {

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
	fmt.Println(cmd.Process.Pid)
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
	if running := isRunning(pid); running {
		err = p.Signal(syscall.SIGTERM)
		if err != nil {
			return err
		}
	}
	return nil
}
