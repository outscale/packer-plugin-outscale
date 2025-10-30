package common

import (
	"context"
	"fmt"
	"log"

	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"

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

func (s *StepDeregisterOMI) Run(
	ctx context.Context,
	state multistep.StateBag,
) multistep.StepAction {
	// Check for force deregister
	if !s.ForceDeregister {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)

	s.Regions = append(s.Regions, s.AccessConfig.GetRegion())

	log.Printf("LOG_ s.Regions: %#+v\n", s.Regions)

	for _, region := range s.Regions {
		// get new connection for each region in which we need to deregister vms
		conn, err := s.AccessConfig.NewOSCClientByRegion(region)
		if err != nil {
			err := fmt.Errorf("error describing OMI: %w", err)
			state.Put("error", err)
			ui.Error(err.Error())

			return multistep.ActionHalt
		}

		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageNames: &[]string{s.OMIName}},
		}
		resp, err := conn.ReadImages(ctx, filterReq)
		if err != nil {
			err := fmt.Errorf("error describing OMI: %w", err)
			state.Put("error", err)
			ui.Error(err.Error())

			return multistep.ActionHalt
		}

		log.Printf("LOG_ resp.Images: %#+v\n", resp.Images)

		// Deregister image(s) by name
		for _, img := range *resp.Images {
			// We are supposing that DeleteImage does the same action as DeregisterImage
			_, err := conn.DeleteImage(ctx, oscgo.DeleteImageRequest{
				ImageId: img.ImageId,
			})
			if err != nil {
				err := fmt.Errorf("error deregistering existing OMI: %w", err)
				state.Put("error", err)
				ui.Error(err.Error())

				return multistep.ActionHalt
			}

			ui.Say(
				fmt.Sprintf(
					"Deregistered OMI %s, id: %s",
					s.OMIName,
					img.ImageId,
				),
			)

			// Delete snapshot(s) by image
			if s.ForceDeleteSnapshot {
				for _, b := range *img.BlockDeviceMappings {
					if b.Bsu.SnapshotId != nil {
						request := oscgo.DeleteSnapshotRequest{SnapshotId: *b.Bsu.SnapshotId}
						_, err := conn.DeleteSnapshot(ctx, request)
						if err != nil {
							err := fmt.Errorf("error deleting existing snapshot: %w", err)
							state.Put("error", err)
							ui.Error(err.Error())

							return multistep.ActionHalt
						}

						ui.Say(fmt.Sprintf("Deleted snapshot: %s", *b.Bsu.SnapshotId))
					}
				}
			}
		}
	}

	return multistep.ActionContinue
}

func (s *StepDeregisterOMI) Cleanup(state multistep.StateBag) {
}
