package chroot

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepCreateOMI creates the OMI.
type StepCreateOMI struct {
	RootVolumeSize int64
	RawRegion      string
	ProductCodes   []string
}

func (s *StepCreateOMI) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	osconn := state.Get("osc").(*(osccommon.OscClient))
	snapshotId := state.Get("snapshot_id").(string)
	ui := state.Get("ui").(packersdk.Ui)

	//regionconn := config.NewOSCClientByRegion(region)
	ui.Say("Creating the OMI...")

	var (
		registerOpts   oscgo.CreateImageRequest
		mappings       []oscgo.BlockDeviceMappingImage
		image          oscgo.Image
		rootDeviceName string
	)

	if config.FromScratch {
		mappings = config.OMIBlockDevices.BuildOscOMIDevices()
		rootDeviceName = config.RootDeviceName
	} else {
		image = state.Get("source_image").(oscgo.Image)
		mappings = image.GetBlockDeviceMappings()
		rootDeviceName = image.GetRootDeviceName()
	}

	newMappings := make([]oscgo.BlockDeviceMappingImage, len(mappings))
	for i, device := range mappings {
		newDevice := device

		//FIX: Temporary fix
		gibSize := newDevice.Bsu.GetVolumeSize() / (1024 * 1024 * 1024)
		newDevice.Bsu.VolumeSize = &gibSize

		if newDevice.GetDeviceName() == rootDeviceName {
			if newDevice.Bsu != nil {
				newDevice.Bsu.SnapshotId = &snapshotId
			} else {
				newDevice.Bsu = &oscgo.BsuToCreate{SnapshotId: &snapshotId}
			}

			if config.FromScratch || int32(s.RootVolumeSize) > newDevice.Bsu.GetVolumeSize() {
				*newDevice.Bsu.VolumeSize = int32(s.RootVolumeSize)
			}
		}

		newMappings[i] = newDevice
	}

	if config.FromScratch {
		architecture := "x86_64"
		registerOpts = oscgo.CreateImageRequest{
			ImageName:           &config.OMIName,
			Architecture:        &architecture,
			RootDeviceName:      &rootDeviceName,
			BlockDeviceMappings: &newMappings,
		}
	} else {
		registerOpts = buildRegisterOpts(config, image, newMappings)
	}

	if config.OMIDescription != "" {
		registerOpts.Description = &config.OMIDescription
	}
	if config.ProductCodes != nil {
		registerOpts.ProductCodes = &config.ProductCodes
	}

	registerResp, _, err := osconn.Api.ImageApi.CreateImage(osconn.Auth).CreateImageRequest(registerOpts).Execute()
	if err != nil {
		state.Put("error", fmt.Errorf("Error registering OMI: %s", err))
		ui.Error(state.Get("error").(error).Error())
		return multistep.ActionHalt
	}

	imageID := registerResp.Image.GetImageId()

	// Set the OMI ID in the state
	ui.Say(fmt.Sprintf("OMI: %s", imageID))
	omis := make(map[string]string)
	omis[s.RawRegion] = imageID
	state.Put("omis", omis)

	ui.Say("Waiting for OMI to become ready...")
	if err := osccommon.WaitUntilOscImageAvailable(osconn, imageID); err != nil {
		err := fmt.Errorf("Error waiting for OMI: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepCreateOMI) Cleanup(state multistep.StateBag) {}

func buildRegisterOpts(config *Config, image oscgo.Image, mappings []oscgo.BlockDeviceMappingImage) oscgo.CreateImageRequest {
	registerOpts := oscgo.CreateImageRequest{
		ImageName:           &config.OMIName,
		Architecture:        image.Architecture,
		RootDeviceName:      image.RootDeviceName,
		BlockDeviceMappings: &mappings,
	}
	return registerOpts
}
