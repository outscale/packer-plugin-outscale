package testacc

import (
	"fmt"
	"os"
	"os/exec"
)

func CheckWithLogs(buildCommand *exec.Cmd, logfile string) error {
	if buildCommand.ProcessState == nil || buildCommand.ProcessState.ExitCode() == 0 {
		return nil
	}

	logs, err := os.ReadFile(logfile)
	if err != nil {
		return fmt.Errorf("bad exit code. logfile: %s (failed to read logs: %w)", logfile, err)
	}

	return fmt.Errorf("bad exit code. logfile: %s\n%s", logfile, string(logs))
}
