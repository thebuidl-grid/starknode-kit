package process

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

func IsRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err != nil
}
func StartProcess(name, command string, logPath io.Writer, args ...string) error {

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
	if err = writeToPIDFile(cmd.Process.Pid, name); err != nil {
		return fmt.Errorf("Failed to write PID file: %v\n", err)
	}
	return nil
}

func loadProcesses() (Process, error) {
	var processes Process
	processPath := path.Join(InstallDir, ".process")
	if _, err := os.Stat(processPath); os.IsNotExist(err) {
		return Process{}, fmt.Errorf("No running process")
	}
	prb, err := os.ReadFile(processPath)
	if err != nil {
		return Process{}, err
	}

	if err = yaml.Unmarshal(prb, processes); err != nil {
		return Process{}, err
	}

	return processes, nil

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
