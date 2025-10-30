package bsu

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

type stepCreateOMI struct {
	image        *oscgo.Image
	RawRegion    string
	ProductCodes []string
	BootModes    []oscgo.BootMode
}

func (s *stepCreateOMI) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	oscconn := state.Get("osc").(*osccommon.OscClient)
	vm := state.Get("vm").(oscgo.Vm)
	ui := state.Get("ui").(packersdk.Ui)

	// Create the image
	omiName := config.OMIName

	ui.Say(fmt.Sprintf("Creating OMI %s from vm %s", omiName, vm.VmId))
	blockDeviceMapping := config.BlockDevices.BuildOscOMIDevices()
	createOpts := oscgo.CreateImageRequest{
		ImageName: &omiName,
	}
	if len(blockDeviceMapping) == 0 {
		createOpts.VmId = &vm.VmId
	} else {
		createOpts.BlockDeviceMappings = &blockDeviceMapping
		if rootDName := config.RootDeviceName; rootDName != "" {
			createOpts.RootDeviceName = &rootDName
		} else {
			err := errors.New("error: MissingParameter: You must provide 'RootDeviceName' when creating omi with 'omi_block_device_mappings'")
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}
	if prodCode := config.ProductCodes; prodCode != nil {
		createOpts.ProductCodes = &prodCode
	}
	if len(config.OMIBootModes) > 0 {
		createOpts.BootModes = config.GetBootModes()
	}

	if description := config.OMIDescription; description != "" {
		createOpts.Description = &description
	}
	if prodCode := config.ProductCodes; prodCode != nil {
		createOpts.ProductCodes = &prodCode
	}

	resp, err := oscconn.CreateImage(ctx, createOpts)
	if err != nil {
		err := fmt.Errorf("error creating OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	image := *resp.Image

	// Set the OMI ID in the state
	ui.Message(fmt.Sprintf("OMI: %s", image.ImageId))
	omis := make(map[string]string)
	omis[s.RawRegion] = image.ImageId
	state.Put("omis", omis)

	// Wait for the image to become ready
	ui.Say("Waiting for OMI to become ready...")
	if err := osccommon.WaitUntilOscImageAvailable(oscconn, image.ImageId); err != nil {
		log.Printf("Error waiting for OMI: %s", err)
		req := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{image.ImageId}},
		}
		imagesResp, err := oscconn.ReadImages(ctx, req)
		if err != nil {
			log.Printf("Unable to determine reason waiting for OMI failed: %s", err)
			err = errors.New("unknown error waiting for OMI")
		} else {
			stateReason := (*imagesResp.Images)[0].StateComment
			err = fmt.Errorf("error waiting for OMI. Reason: %s", *stateReason.StateMessage)
		}

		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	req := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{image.ImageId}},
	}
	imagesResp, err := oscconn.ReadImages(ctx, req)
	if err != nil {
		err := fmt.Errorf("error searching for OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	if len(*imagesResp.Images) <= 0 {
		err := fmt.Errorf("error while reading the image': %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	s.image = &(*imagesResp.Images)[0]
	if s.image == nil {
		err := fmt.Errorf("error while reading an empty image id': %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	snapshots := make(map[string][]string)
	blockMapping := (*imagesResp.Images)[0].BlockDeviceMappings
	for _, blockDeviceMapping := range *blockMapping {
		if blockDeviceMapping.Bsu.SnapshotId != nil {
			snapshots[s.RawRegion] = append(
				snapshots[s.RawRegion],
				*blockDeviceMapping.Bsu.SnapshotId,
			)
		}
	}
	state.Put("snapshots", snapshots)

	return multistep.ActionContinue
}

func (s *stepCreateOMI) Cleanup(state multistep.StateBag) {
	if s.image == nil {
		return
	}

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !cancelled && !halted {
		return
	}

	oscconn := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Deregistering the OMI because cancellation or error...")
	deleteOpts := oscgo.DeleteImageRequest{ImageId: s.image.ImageId}
	_, err := oscconn.DeleteImage(context.Background(), deleteOpts)
	if err != nil {
		ui.Error(fmt.Sprintf("error Deleting OMI, may still be around:': %v", err))
		return
	}
}
