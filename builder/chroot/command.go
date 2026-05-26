package chroot

import (
	"context"
	"os/exec"
)

// CommandWrapper is a type that given a command, will possibly modify that
// command in-flight. This might return an error.
type CommandWrapper func(string) (string, error)

// ShellCommand takes a command string and returns an *exec.Cmd to execute
// it within the context of a shell (/bin/sh).
func ShellCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "/bin/sh", "-c", command) //nolint:gosec
}
