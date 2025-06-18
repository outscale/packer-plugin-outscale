package chroot

import (
	"context"
	"fmt"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	sl "github.com/hashicorp/packer-plugin-sdk/shell-local"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

func RunLocalCommands(commands []string, wrappedCommand CommandWrapper, ictx interpolate.Context, ui packersdk.Ui) error {
	ctx := context.TODO()
	for _, rawCmd := range commands {
		intCmd, err := interpolate.Render(rawCmd, &ictx)
		if err != nil {
			return fmt.Errorf("error interpolating: %w", err)
		}

		command, err := wrappedCommand(intCmd)
		if err != nil {
			return fmt.Errorf("error wrapping command: %w", err)
		}

		ui.Say(fmt.Sprintf("executing command: %s", command))
		comm := &sl.Communicator{
			ExecuteCommand: []string{"sh", "-c", command},
		}
		cmd := &packersdk.RemoteCmd{Command: command}
		if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
			return fmt.Errorf("error executing command: %w", err)
		}
		if cmd.ExitStatus() != 0 {
			return fmt.Errorf(
				"received non-zero exit code %d from command: %s",
				cmd.ExitStatus(),
				command)
		}
	}
	return nil
}
