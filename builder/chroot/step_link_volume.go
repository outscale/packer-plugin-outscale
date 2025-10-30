package chroot

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepLinkVolume attaches the previously created volume to an
// available device location.
//
// Produces:
//
//	device string - The location where the volume was attached.
//	attach_cleanup CleanupFunc
type StepLinkVolume struct {
	attached bool
	volumeId string
}

func (s *StepLinkVolume) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	oscconn := state.Get("osc").(*osccommon.OscClient)
	device := state.Get("device").(string)
	vm := state.Get("vm").(oscgo.Vm)
	ui := state.Get("ui").(packersdk.Ui)
	volumeId := state.Get("volume_id").(string)

	// For the API call, it expects "sd" prefixed devices.
	// linkVolume := strings.Replace(device, "/xvd", "/sd", 1)
	linkVolume := device

	ui.Say(fmt.Sprintf("Attaching the root volume to %s", linkVolume))
	opts := oscgo.LinkVolumeRequest{
		DeviceName: linkVolume,
		VmId:       vm.VmId,
		VolumeId:   volumeId,
	}
	_, err := oscconn.LinkVolume(ctx, opts)
	if err != nil {
		err := fmt.Errorf("error attaching volume: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Mark that we attached it so we can detach it later
	s.attached = true
	s.volumeId = volumeId

	// Wait for the volume to become attached
	err = osccommon.WaitUntilOscVolumeIsLinked(oscconn, s.volumeId)
	if err != nil {
		err := fmt.Errorf("error waiting for volume: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("attach_cleanup", s)
	return multistep.ActionContinue
}

func (s *StepLinkVolume) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	if err := s.CleanupFunc(state); err != nil {
		ui.Error(err.Error())
	}
}

func (s *StepLinkVolume) CleanupFunc(state multistep.StateBag) error {
	if !s.attached {
		return nil
	}

	oscconn := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Detaching BSU volume...")
	opts := oscgo.UnlinkVolumeRequest{
		VolumeId: s.volumeId,
	}
	_, err := oscconn.UnlinkVolume(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("error detaching BSU volume: %w", err)
	}

	s.attached = false

	// Wait for the volume to detach
	err = osccommon.WaitUntilOscVolumeIsUnlinked(oscconn, s.volumeId)
	if err != nil {
		return fmt.Errorf("error waiting for volume: %w", err)
	}

	return nil
}
