package chroot

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// StepFlock provisions the instance within a chroot.
//
// Produces:
//
//	flock_cleanup Cleanup - To perform early cleanup
type StepFlock struct {
	fh *os.File
}

func (s *StepFlock) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)

	lockfile := "/var/lock/packer-chroot/lock"
	if err := os.MkdirAll(filepath.Dir(lockfile), 0755); err != nil {
		err := fmt.Errorf("error creating lock: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("Obtaining lock: %s", lockfile)
	f, err := os.Create(lockfile)
	if err != nil {
		err := fmt.Errorf("error creating lock: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// LOCK!
	if err := lockFile(f); err != nil {
		err := fmt.Errorf("error obtaining lock: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Set the file handle, we can't close it because we need to hold
	// the lock.
	s.fh = f

	state.Put("flock_cleanup", s)
	return multistep.ActionContinue
}

func (s *StepFlock) Cleanup(state multistep.StateBag) {
	_ = s.CleanupFunc(state)
}

func (s *StepFlock) CleanupFunc(state multistep.StateBag) error {
	if s.fh == nil {
		return nil
	}

	log.Printf("Unlocking: %s", s.fh.Name())
	if err := unlockFile(s.fh); err != nil {
		return err
	}

	s.fh = nil
	return nil
}
