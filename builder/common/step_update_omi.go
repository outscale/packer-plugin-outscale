package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

type StepUpdateOMIAttributes struct {
	AccountIds         []string
	SnapshotAccountIds []string
	RawRegion          string
	GlobalPermission   bool
	Ctx                interpolate.Context
}

func (s *StepUpdateOMIAttributes) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui, ok := state.Get("ui").(packersdk.Ui)
	if !ok {
		return multistep.ActionContinue
	}
	config, ok := state.Get("accessConfig").(*AccessConfig)
	if !ok {
		return multistep.ActionContinue
	}
	omis := state.Get("omis").(map[string]string)
	snapshots := state.Get("snapshots").(map[string][]string)

	// Determine if there is any work to do.
	valid := false
	valid = valid || (s.AccountIds != nil && len(s.AccountIds) > 0)
	valid = valid || (s.SnapshotAccountIds != nil && len(s.SnapshotAccountIds) > 0)
	valid = valid || (s.GlobalPermission != false)

	if !valid {
		return multistep.ActionContinue
	}

	s.Ctx.Data = extractBuildInfo(s.RawRegion, state)

	updateSnapshoptRequest := oscgo.UpdateSnapshotRequest{
		PermissionsToCreateVolume: oscgo.PermissionsOnResourceCreation{
			Additions: &oscgo.PermissionsOnResource{
				AccountIds: &s.AccountIds,
			},
		},
	}

	updateImageRequest := oscgo.UpdateImageRequest{
		PermissionsToLaunch: oscgo.PermissionsOnResourceCreation{
			Additions: &oscgo.PermissionsOnResource{
				AccountIds:       &s.AccountIds,
				GlobalPermission: &s.GlobalPermission,
			},
		},
	}

	// Updating image attributes
	for region, omi := range omis {
		ui.Say(fmt.Sprintf("Updating attributes on OMI (%s)...", omi))
		regionconn := config.NewOSCClientByRegion(region)

		ui.Message(fmt.Sprintf("Updating: %s", omi))
		updateImageRequest.ImageId = omi
		_, _, err := regionconn.Api.ImageApi.UpdateImage(regionconn.Auth).UpdateImageRequest(updateImageRequest).Execute()
		if err != nil {
			err := fmt.Errorf("Error updating OMI: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	// Updating snapshot attributes
	for region, region_snapshots := range snapshots {
		for _, snapshot := range region_snapshots {
			ui.Say(fmt.Sprintf("Updating attributes on snapshot (%s)...", snapshot))
			regionconn := config.NewOSCClientByRegion(region)

			ui.Message(fmt.Sprintf("Updating: %s", snapshot))
			updateSnapshoptRequest.SnapshotId = snapshot
			_, _, err := regionconn.Api.SnapshotApi.UpdateSnapshot(regionconn.Auth).UpdateSnapshotRequest(updateSnapshoptRequest).Execute()
			if err != nil {
				err := fmt.Errorf("Error updating snapshot: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}

		}
	}

	return multistep.ActionContinue
}

func (s *StepUpdateOMIAttributes) Cleanup(state multistep.StateBag) {
	// No cleanup...
}
