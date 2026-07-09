package testacc

import (
	"fmt"
	"os/exec"
)

func CheckWithLogs(buildCommand *exec.Cmd, logfile string) error {
	if buildCommand.ProcessState == nil || buildCommand.ProcessState.ExitCode() == 0 {
		return nil
	}

	return fmt.Errorf("bad exit code. logfile: %s", logfile)
}
