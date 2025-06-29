package chroot

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// StepCopyFiles copies some files from the host into the chroot environment.
//
// Produces:
//
//	copy_files_cleanup CleanupFunc - A function to clean up the copied files
//	early.
type StepCopyFiles struct {
	files []string
}

func (s *StepCopyFiles) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	mountPath := state.Get("mount_path").(string)
	ui := state.Get("ui").(packersdk.Ui)
	wrappedCommand := state.Get("wrappedCommand").(CommandWrapper)
	stderr := new(bytes.Buffer)

	s.files = make([]string, 0, len(config.CopyFiles))
	if len(config.CopyFiles) > 0 {
		ui.Say("Copying files from host to chroot...")
		for _, path := range config.CopyFiles {
			ui.Message(path)
			chrootPath := filepath.Join(mountPath, path)
			log.Printf("Copying '%s' to '%s'", path, chrootPath)

			cmdText, err := wrappedCommand(fmt.Sprintf("cp --remove-destination %s %s", path, chrootPath))
			if err != nil {
				err := fmt.Errorf("error building copy command: %w", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}

			stderr.Reset()
			cmd := ShellCommand(cmdText)
			cmd.Stderr = stderr
			if err := cmd.Run(); err != nil {
				err := fmt.Errorf(
					"error copying file: %w\nnStderr: %s", err, stderr.String())
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}

			s.files = append(s.files, chrootPath)
		}
	}

	state.Put("copy_files_cleanup", s)
	return multistep.ActionContinue
}

func (s *StepCopyFiles) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	if err := s.CleanupFunc(state); err != nil {
		ui.Error(err.Error())
	}
}

func (s *StepCopyFiles) CleanupFunc(state multistep.StateBag) error {
	wrappedCommand := state.Get("wrappedCommand").(CommandWrapper)
	if s.files != nil {
		for _, file := range s.files {
			log.Printf("Removing: %s", file)
			localCmdText, err := wrappedCommand(fmt.Sprintf("rm -f %s", file))
			if err != nil {
				return err
			}

			localCmd := ShellCommand(localCmdText)
			if err := localCmd.Run(); err != nil {
				return err
			}
		}
	}

	s.files = nil
	return nil
}
