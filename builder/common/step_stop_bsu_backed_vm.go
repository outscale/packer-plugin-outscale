package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/packer-plugin-outscale/builder/common/retry"
)

type StepStopBSUBackedVm struct {
	Skip          bool
	DisableStopVm bool
}

func (s *StepStopBSUBackedVm) Run(
	ctx context.Context,
	state multistep.StateBag,
) multistep.StepAction {
	oscconn := state.Get("osc").(*OscClient)
	vm := state.Get("vm").(oscgo.Vm)
	ui := state.Get("ui").(packersdk.Ui)

	// Skip when it is a spot vm
	if s.Skip {
		return multistep.ActionContinue
	}

	var err error

	if !s.DisableStopVm {
		// Stop the vm so we can create an AMI from it
		ui.Say("Stopping the source vm...")

		// Amazon EC2 API follows an eventual consistency model.

		// This means that if you run a command to modify or describe a resource
		// that you just created, its ID might not have propagated throughout
		// the system, and you will get an error responding that the resource
		// does not exist.

		// Work around this by retrying a few times, up to about 5 minutes.
		err := retry.Run(10, 60, 6, func(i uint) (bool, error) {
			ui.Message(fmt.Sprintf("Stopping vm, attempt %d", i+1))

			_, err = oscconn.StopVms(ctx, oscgo.StopVmsRequest{
				VmIds: []string{vm.VmId},
			})
			if err == nil {
				// success
				return true, nil
			}

			// TODO: manager errors

			// errored, but not in expected way. Don't want to retry
			return true, err
		})
		if err != nil {
			err := fmt.Errorf("error stopping vm: %w", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

	} else {
		ui.Say("Automatic vm stop disabled. Please stop vm manually.")
	}

	// Wait for the vm to actually stop
	ui.Say("Waiting for the vm to stop...")
	switch vm.VmInitiatedShutdownBehavior {
	case StopShutdownBehavior:
		err = waitUntilOscVmStopped(oscconn, vm.VmId)
	case TerminateShutdownBehavior:
		err = waitUntilOscVmDeleted(oscconn, vm.VmId)
	default:
		err := errors.New("wrong value for the shutdown behavior")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt

	}

	if err != nil {
		err := fmt.Errorf("error waiting for vm to stop: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepStopBSUBackedVm) Cleanup(multistep.StateBag) {
	// No cleanup...
}
