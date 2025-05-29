package pkg

import (
	"fmt"
	"io"
	"syscall"
)

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
