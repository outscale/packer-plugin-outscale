package chroot

import (
	"context"
	"errors"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// StepCheckRootDevice makes sure the root device on the OMI is BSU-backed.
type StepCheckRootDevice struct{}

func (s *StepCheckRootDevice) Run(
	_ context.Context,
	state multistep.StateBag,
) multistep.StepAction {
	image := state.Get("source_image").(oscgo.Image)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Checking the root device on source OMI...")

	// It must be BSU-backed otherwise the build won't work
	if image.RootDeviceType != "ebs" {
		err := errors.New("the root device of the source OMI must be BSU-backed. ")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepCheckRootDevice) Cleanup(multistep.StateBag) {}
