package common

import (
	"context"
	"fmt"
	"log"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepDeregisterOMI struct {
	AccessConfig        *AccessConfig
	ForceDeregister     bool
	ForceDeleteSnapshot bool
	OMIName             string
	Regions             []string
}

func (s *StepDeregisterOMI) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	// Check for force deregister
	if !s.ForceDeregister {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)

	s.Regions = append(s.Regions, s.AccessConfig.GetRegion())

	log.Printf("LOG_ s.Regions: %#+v\n", s.Regions)

	for _, region := range s.Regions {
		// get new connection for each region in which we need to deregister vms
		conn := s.AccessConfig.NewOSCClientByRegion(region)

		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageNames: &[]string{s.OMIName}},
		}
		resp, _, err := conn.Api.ImageApi.ReadImages(conn.Auth).ReadImagesRequest(filterReq).Execute()
		if err != nil {
			err := fmt.Errorf("Error describing OMI: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())

			return multistep.ActionHalt
		}

		log.Printf("LOG_ resp.Images: %#+v\n", resp.Images)

		// Deregister image(s) by name
		for i := range resp.GetImages() {
			//We are supposing that DeleteImage does the same action as DeregisterImage
			_, _, err := conn.Api.ImageApi.DeleteImage(conn.Auth).DeleteImageRequest(oscgo.DeleteImageRequest{
				ImageId: resp.GetImages()[i].GetImageId(),
			}).Execute()
			if err != nil {
				err := fmt.Errorf("Error deregistering existing OMI: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())

				return multistep.ActionHalt
			}

			ui.Say(fmt.Sprintf("Deregistered OMI %s, id: %s", s.OMIName, resp.GetImages()[i].GetImageId()))

			// Delete snapshot(s) by image
			if s.ForceDeleteSnapshot {
				for _, b := range resp.GetImages()[i].GetBlockDeviceMappings() {
					if b.Bsu.SnapshotId != nil {
						request := oscgo.DeleteSnapshotRequest{SnapshotId: *b.GetBsu().SnapshotId}
						_, _, err := conn.Api.SnapshotApi.DeleteSnapshot(conn.Auth).DeleteSnapshotRequest(request).Execute()
						if err != nil {
							err := fmt.Errorf("Error deleting existing snapshot: %s", err)
							state.Put("error", err)
							ui.Error(err.Error())

							return multistep.ActionHalt
						}

						ui.Say(fmt.Sprintf("Deleted snapshot: %s", *b.GetBsu().SnapshotId))
					}
				}
			}
		}
	}

	return multistep.ActionContinue
}

func (s *StepDeregisterOMI) Cleanup(state multistep.StateBag) {
}
