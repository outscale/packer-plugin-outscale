package chroot

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/osc/common"
)

// StepSnapshot creates a snapshot of the created volume.
//
// Produces:
//
//	snapshot_id string - ID of the created snapshot
type StepSnapshot struct {
	snapshotId string
	RawRegion  string
}

func (s *StepSnapshot) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	oscconn := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)
	volumeId := state.Get("volume_id").(string)

	ui.Say("Creating snapshot...")
	description := fmt.Sprintf("Packer: %s", time.Now().String())

	request := oscgo.CreateSnapshotRequest{
		Description: &description,
		VolumeId:    &volumeId,
	}
	createSnapResp, _, err := oscconn.Api.SnapshotApi.CreateSnapshot(oscconn.Auth).CreateSnapshotRequest(request).Execute()
	if err != nil {
		err := fmt.Errorf("Error creating snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Set the snapshot ID so we can delete it later
	s.snapshotId = createSnapResp.Snapshot.GetSnapshotId()
	ui.Message(fmt.Sprintf("Snapshot ID: %s", s.snapshotId))

	// Wait for the snapshot to be ready
	err = osccommon.WaitUntilOscSnapshotDone(oscconn, s.snapshotId)
	if err != nil {
		err := fmt.Errorf("Error waiting for snapshot: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("snapshot_id", s.snapshotId)

	snapshots := map[string][]string{
		s.RawRegion: {s.snapshotId},
	}
	state.Put("snapshots", snapshots)

	return multistep.ActionContinue
}

func (s *StepSnapshot) Cleanup(state multistep.StateBag) {
	if s.snapshotId == "" {
		return
	}

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)

	if cancelled || halted {
		oscconn := state.Get("osc").(*osccommon.OscClient)
		ui := state.Get("ui").(packersdk.Ui)
		ui.Say("Removing snapshot since we cancelled or halted...")
		request := oscgo.DeleteSnapshotRequest{SnapshotId: s.snapshotId}
		_, _, err := oscconn.Api.SnapshotApi.DeleteSnapshot(oscconn.Auth).DeleteSnapshotRequest(request).Execute()
		if err != nil {
			ui.Error(fmt.Sprintf("Error: %s", err))
		}
	}
}
