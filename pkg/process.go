package pkg

import (
	"fmt"
	"os"
	"syscall"
)

func StartProcess(name, command, logsPath string, args ...string) error {
	loggerName := fmt.Sprintf("%s.log", name)
	logger, err := os.OpenFile(loggerName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	cmd := execCommand(command, args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	cmd.Stdout = logger
	cmd.Stderr = logger
	err = cmd.Start()
	if err != nil {
		return err
	}
	fmt.Println(cmd.Process.Pid)
	return nil
}
