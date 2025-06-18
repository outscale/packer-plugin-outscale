package bsusurrogate

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepRegisterOMI creates the OMI.
type StepRegisterOMI struct {
	RootDevice    RootBlockDevice
	OMIDevices    []oscgo.BlockDeviceMappingImage
	LaunchDevices []oscgo.BlockDeviceMappingVmCreation
	image         *oscgo.Image
	RawRegion     string
	ProductCodes  []string
	BootModes     []oscgo.BootMode
}

func (s *StepRegisterOMI) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	oscconn := state.Get("osc").(*(osccommon.OscClient))
	snapshotIds := state.Get("snapshot_ids").(map[string]string)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Registering the OMI...")

	blockDevices := s.combineDevices(snapshotIds)
	architecture := "x86_64"
	registerOpts := oscgo.CreateImageRequest{
		ImageName:           &config.OMIName,
		Architecture:        &architecture,
		RootDeviceName:      &s.RootDevice.DeviceName,
		BlockDeviceMappings: &blockDevices,
	}

	if config.OMIDescription != "" {
		registerOpts.Description = &config.OMIDescription
	}
	if config.ProductCodes != nil {
		registerOpts.ProductCodes = &config.ProductCodes
	}
	if len(config.OMIBootModes) > 0 {
		registerOpts.SetBootModes(config.GetBootModes())
	}
	registerResp, _, err := oscconn.Api.ImageApi.CreateImage(oscconn.Auth).CreateImageRequest(registerOpts).Execute()
	if err != nil {
		state.Put("error", fmt.Errorf("error registering OMI: %w", err))
		ui.Error(state.Get("error").(error).Error())
		return multistep.ActionHalt
	}

	// Set the OMI ID in the state
	ui.Say(fmt.Sprintf("OMI: %s", *registerResp.GetImage().ImageId))
	omis := make(map[string]string)
	omis[s.RawRegion] = *registerResp.GetImage().ImageId
	state.Put("omis", omis)

	// Wait for the image to become ready
	ui.Say("Waiting for OMI to become ready...")
	if err := osccommon.WaitUntilOscImageAvailable(oscconn, *registerResp.GetImage().ImageId); err != nil {
		err := fmt.Errorf("error waiting for OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	filterReq := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{ImageIds: &[]string{*registerResp.GetImage().ImageId}},
	}
	imagesResp, _, err := oscconn.Api.ImageApi.ReadImages(oscconn.Auth).ReadImagesRequest(filterReq).Execute()
	if err != nil {
		err := fmt.Errorf("error searching for OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	s.image = &imagesResp.GetImages()[0]

	snapshots := make(map[string][]string)
	block := imagesResp.GetImages()[0].BlockDeviceMappings
	for _, blockDeviceMapping := range *block {
		if blockDeviceMapping.Bsu.GetSnapshotId() != "" {
			snapshots[s.RawRegion] = append(snapshots[s.RawRegion], blockDeviceMapping.Bsu.GetSnapshotId())
		}
	}
	state.Put("snapshots", snapshots)

	return multistep.ActionContinue
}

func (s *StepRegisterOMI) Cleanup(state multistep.StateBag) {
	if s.image == nil {
		return
	}

	_, cancelled := state.GetOk(multistep.StateCancelled)
	_, halted := state.GetOk(multistep.StateHalted)
	if !cancelled && !halted {
		return
	}

	oscconn := state.Get("osc").(*(osccommon.OscClient))
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Deregistering the OMI because cancellation or error...")
	deregisterOpts := oscgo.DeleteImageRequest{ImageId: *s.image.ImageId}
	_, _, err := oscconn.Api.ImageApi.DeleteImage(oscconn.Auth).DeleteImageRequest(deregisterOpts).Execute()
	if err != nil {
		ui.Error(fmt.Sprintf("error deregistering OMI, may still be around: %s", err.Error()))
		return
	}
}

func (s *StepRegisterOMI) combineDevices(snapshotIDs map[string]string) []oscgo.BlockDeviceMappingImage {
	devices := map[string]oscgo.BlockDeviceMappingImage{}

	for _, device := range s.OMIDevices {
		devices[device.GetDeviceName()] = device
	}

	// Devices in launch_block_device_mappings override any with
	// the same name in ami_block_device_mappings, except for the
	// one designated as the root device in omi_root_device
	for _, device := range s.LaunchDevices {
		snapshotID, ok := snapshotIDs[device.GetDeviceName()]
		if ok && snapshotID != "" {
			device.Bsu.SnapshotId = &snapshotID
		}
		if device.GetDeviceName() == s.RootDevice.SourceDeviceName {
			device.DeviceName = &s.RootDevice.DeviceName

			if _, ok := device.Bsu.GetVolumeTypeOk(); ok {
				device.Bsu.VolumeType = &s.RootDevice.VolumeType
				if device.Bsu.GetVolumeType() != "io1" {
					device.Bsu.Iops = nil
				}
			}

		}
		devices[device.GetDeviceName()] = copyToDeviceMappingImage(device)
	}

	blockDevices := []oscgo.BlockDeviceMappingImage{}
	for _, device := range devices {
		blockDevices = append(blockDevices, device)
	}
	return blockDevices
}

func copyToDeviceMappingImage(device oscgo.BlockDeviceMappingVmCreation) oscgo.BlockDeviceMappingImage {
	log.Printf("Copy device mapping image ")
	deviceImage := oscgo.BlockDeviceMappingImage{
		DeviceName: device.DeviceName,
		Bsu: &oscgo.BsuToCreate{
			DeleteOnVmDeletion: device.Bsu.DeleteOnVmDeletion,
			Iops:               device.Bsu.Iops,
			SnapshotId:         device.Bsu.SnapshotId,
			VolumeSize:         device.Bsu.VolumeSize,
			VolumeType:         device.Bsu.VolumeType,
		},
	}
	return deviceImage
}
