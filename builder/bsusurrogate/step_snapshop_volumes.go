package bsusurrogate

import (
	"context"
	"fmt"
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepSnapshotVolumes creates snapshots of the created volumes.
//
// Produces:
//
//	snapshot_ids map[string]string - IDs of the created snapshots
type StepSnapshotVolumes struct {
	LaunchDevices []oscgo.BlockDeviceMappingVmCreation
	snapshotIds   map[string]string
}

func (s *StepSnapshotVolumes) snapshotVolume(ctx context.Context, deviceName string, state multistep.StateBag) error {
	oscconn := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)
	vm := state.Get("vm").(oscgo.Vm)

	var volumeId string
	for _, volume := range *vm.BlockDeviceMappings {
		if volume.GetDeviceName() == deviceName {
			volumeId = *volume.GetBsu().VolumeId
		}
	}
	if volumeId == "" {
		return fmt.Errorf("volume ID for device %s not found", deviceName)
	}

	ui.Say(fmt.Sprintf("Creating snapshot of EBS Volume %s...", volumeId))
	description := fmt.Sprintf("Packer: %s", time.Now().String())

	request := oscgo.CreateSnapshotRequest{
		Description: &description,
		VolumeId:    &volumeId,
	}
	createSnapResp, _, err := oscconn.Api.SnapshotApi.CreateSnapshot(oscconn.Auth).CreateSnapshotRequest(request).Execute()
	if err != nil {
		return err
	}

	// Set the snapshot ID so we can delete it later
	s.snapshotIds[deviceName] = *createSnapResp.Snapshot.SnapshotId

	// Wait for snapshot to be created
	err = osccommon.WaitUntilOscSnapshotCompleted(oscconn, *createSnapResp.Snapshot.SnapshotId)
	return err
}

func (s *StepSnapshotVolumes) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)

	s.snapshotIds = map[string]string{}

	var wg sync.WaitGroup
	var errs *multierror.Error
	for _, device := range s.LaunchDevices {
		wg.Add(1)
		go func(device oscgo.BlockDeviceMappingVmCreation) {
			defer wg.Done()
			if err := s.snapshotVolume(ctx, device.GetDeviceName(), state); err != nil {
				errs = multierror.Append(errs, err)
			}
		}(device)
	}

	wg.Wait()

	if errs != nil {
		state.Put("error", errs)
		ui.Error(errs.Error())
		return multistep.ActionHalt
	}

	state.Put("snapshot_ids", s.snapshotIds)
	return multistep.ActionContinue
}

func (s *StepSnapshotVolumes) Cleanup(state multistep.StateBag) {
	if len(s.snapshotIds) == 0 {
		return
	}

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)

	if cancelled || halted {
		oscconn := state.Get("osc").(*osccommon.OscClient)
		ui := state.Get("ui").(packersdk.Ui)
		ui.Say("Removing snapshots since we cancelled or halted...")
		for _, snapshotID := range s.snapshotIds {
			request := oscgo.DeleteSnapshotRequest{SnapshotId: snapshotID}
			_, _, err := oscconn.Api.SnapshotApi.DeleteSnapshot(oscconn.Auth).DeleteSnapshotRequest(request).Execute()
			if err != nil {
				ui.Error(fmt.Sprintf("Error: %s", err))
			}
		}
	}
}
