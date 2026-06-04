package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/outscale/goutils/sdk/ptr"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/packer-plugin-outscale/builder/common/retry"
)

type StepCreateTags struct {
	Tags         TagMap
	SnapshotTags TagMap
	Ctx          interpolate.Context
}

func (s *StepCreateTags) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	if !s.Tags.IsSet() && !s.SnapshotTags.IsSet() {
		return multistep.ActionContinue
	}

	ui := state.Get("ui").(packersdk.Ui)
	config := state.Get("accessConfig").(*AccessConfig)
	omis := state.Get("omis").(map[string]string)

	// Adds tags to OMIs and snapshots
	for _, ami := range omis {
		ui.Say(fmt.Sprintf("Adding tags to OMI (%s)...", ami))

		regionconn, err := config.NewOSCClient()
		if err != nil {
			err = fmt.Errorf("error retrieving details for OMI (%s): %w", ami, err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Retrieve image list for given OMI
		resourceIds := []string{ami}
		imageResp, err := regionconn.ReadImages(ctx, oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{
				ImageIds: &resourceIds,
			},
		})
		if err != nil {
			err = fmt.Errorf("error retrieving details for OMI (%s): %w", ami, err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		images := ptr.From(imageResp.Images)
		if len(images) == 0 {
			err = fmt.Errorf("error retrieving details for OMI (%s), no images found", ami)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		image := images[0]
		snapshotIds := []string{}

		// Add only those with a Snapshot ID, i.e. not Ephemeral
		for _, device := range ptr.From(image.BlockDeviceMappings) {
			snapshotID := ptr.From(ptr.From(device.Bsu).SnapshotId)
			if snapshotID != "" {
				ui.Say("Tagging snapshot: " + snapshotID)
				resourceIds = append(resourceIds, snapshotID)
				snapshotIds = append(snapshotIds, snapshotID)
			}
		}

		// Convert tags to oapi.Tag format
		ui.Say("Creating OMI tags")
		omiTags, err := s.Tags.OSCTags(s.Ctx, config.RawRegion, state)
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		omiTags.Report(ui)

		ui.Say("Creating snapshot tags")
		snapshotTags, err := s.SnapshotTags.OSCTags(s.Ctx, config.RawRegion, state)
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
		snapshotTags.Report(ui)
		// Retry creating tags for about 2.5 minutes
		err = retry.Run(0.2, 30, 11, func(_ uint) (bool, error) {
			// Tag images and snapshots
			request := oscgo.CreateTagsRequest{
				ResourceIds: resourceIds,
				Tags:        omiTags,
			}

			_, err = regionconn.CreateTags(ctx, request)
			if err != nil {
				return false, err
			}

			requestSnap := oscgo.CreateTagsRequest{
				ResourceIds: snapshotIds,
				Tags:        snapshotTags,
			}
			// Override tags on snapshots
			if len(snapshotTags) > 0 {
				_, err = regionconn.CreateTags(ctx, requestSnap)
			}
			if err == nil {
				return true, nil
			}
			return true, err
		})
		if err != nil {
			err = fmt.Errorf("error adding tags to Resources (%#v): %w", resourceIds, err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *StepCreateTags) Cleanup(state multistep.StateBag) {
	// No cleanup...
}
