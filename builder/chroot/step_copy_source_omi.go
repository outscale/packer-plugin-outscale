package chroot

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/outscale/goutils/sdk/ptr"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepCopySourceOMI copies a source OMI into the current account when the
// original image is owned by a different account. This is necessary to be able
// to use the snapshot id associated with the BSUs of the OMI
type StepCopySourceOMI struct {
	Region string

	copiedImageID string
}

func (s *StepCopySourceOMI) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)
	sourceImage := state.Get("source_image").(oscgo.Image)

	resp, err := client.ReadAccounts(ctx, oscgo.ReadAccountsRequest{})
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if len(*resp.Accounts) == 0 {
		err := errors.New("error: no account found")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	accountId := *(*resp.Accounts)[0].AccountId

	if sourceImage.AccountId == accountId {
		return multistep.ActionContinue
	}

	ui.Say("Copying source OMI into the current account...")
	name := fmt.Sprintf("%s-packer-copy-%s", *sourceImage.ImageName, time.Now())
	copyReq := oscgo.CreateImageRequest{
		ImageName:        &name,
		SourceImageId:    &sourceImage.ImageId,
		SourceRegionName: &s.Region,
	}

	copyResp, err := client.CreateImage(ctx, copyReq)
	if err != nil {
		err := fmt.Errorf("error copying source OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	s.copiedImageID = copyResp.Image.ImageId
	ui.Say("Copied source OMI: " + s.copiedImageID)

	if err := osccommon.WaitUntilOscImageAvailable(ctx, client, s.copiedImageID); err != nil {
		err := fmt.Errorf("error waiting for copied OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	imageResp, err := client.ReadImages(ctx, oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{s.copiedImageID}},
	})
	if err != nil {
		err := fmt.Errorf("error reading copied OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	copiedImage := (*imageResp.Images)[0]
	state.Put("source_image", copiedImage)

	return multistep.ActionContinue
}

func (s *StepCopySourceOMI) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	if err := s.CleanupFunc(state); err != nil {
		ui.Error(err.Error())
	}
}

func (s *StepCopySourceOMI) CleanupFunc(state multistep.StateBag) error {
	if s.copiedImageID == "" {
		return nil
	}

	client := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	resp, err := client.ReadImages(context.Background(), oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{s.copiedImageID}},
	})
	if err != nil {
		ui.Error(fmt.Sprintf("error reading copied source OMI %s: %v", s.copiedImageID, err))
	}

	ui.Say("Deleting temporary copied source OMI...")
	_, err = client.DeleteImage(context.Background(), oscgo.DeleteImageRequest{ImageId: s.copiedImageID})
	if err != nil {
		ui.Error(fmt.Sprintf("error deleting temporary copied source OMI %s: %v", s.copiedImageID, err))
	}

	ui.Say("Deleting snapshots of copied source OMI...")
	for _, block := range ptr.From((*resp.Images)[0].BlockDeviceMappings) {
		if block.Bsu != nil {
			_, err = client.DeleteSnapshot(context.Background(), oscgo.DeleteSnapshotRequest{SnapshotId: *block.Bsu.SnapshotId})
			if err != nil {
				ui.Error(fmt.Sprintf("error deleting copied source snapshot %s: %v", *block.Bsu.SnapshotId, err))
			}
		}
	}

	s.copiedImageID = ""
	return nil
}
