package chroot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/outscale/goutils/sdk/ptr"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepCreateVolume creates a new volume from the snapshot of the root
// device of the OMI.
//
// Produces:
//
//	volume_id string - The ID of the created volume
type StepCreateVolume struct {
	volumeId       string
	RootVolumeSize int64
	RootVolumeType string
	RootVolumeTags osccommon.TagMap
	RawRegion      string
	Ctx            interpolate.Context
}

func (s *StepCreateVolume) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	oscconn := state.Get("osc").(*osccommon.OscClient)
	vm := state.Get("vm").(oscgo.Vm)
	ui := state.Get("ui").(packersdk.Ui)

	var err error

	volTags, err := s.RootVolumeTags.OSCTags(s.Ctx, s.RawRegion, state)
	if err != nil {
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	var createVolume *oscgo.CreateVolumeRequest
	if config.FromScratch {
		rootVolumeType := osccommon.VolumeTypeGp2
		if s.RootVolumeType == "io1" {
			err := errors.New("cannot use io1 volume when building from scratch")
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		} else if s.RootVolumeType != "" {
			rootVolumeType = s.RootVolumeType
		}
		size := int(s.RootVolumeSize)
		createVolume = &oscgo.CreateVolumeRequest{
			SubregionName: vm.Placement.SubregionName,
			Size:          &size,
			VolumeType:    new(oscgo.VolumeType(rootVolumeType)),
		}
	} else {
		// Determine the root device snapshot
		image := state.Get("source_image").(oscgo.Image)
		rootDeviceName := ptr.From(image.RootDeviceName)
		log.Printf("Searching for root device of the image (%s)", rootDeviceName)
		var rootDevice *oscgo.BlockDeviceMappingImage
		for _, device := range ptr.From(image.BlockDeviceMappings) {
			if ptr.From(device.DeviceName) == rootDeviceName {
				rootDevice = &device
				break
			}
		}

		ui.Say("Creating the root volume...")
		createVolume, err = s.buildCreateVolumeInput(vm.Placement.SubregionName, rootDevice)
		if err != nil {
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	log.Printf("Create args: %+v", createVolume)

	createVolumeResp, err := oscconn.CreateVolume(ctx, *createVolume)
	if err != nil {
		err := fmt.Errorf("error creating root volume: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Set the volume ID so we remember to delete it later
	s.volumeId = createVolumeResp.Volume.VolumeId
	log.Printf("Volume ID: %s", s.volumeId)

	// Create tags for volume
	if len(volTags) > 0 {
		if err := osccommon.CreateOSCTags(ctx, oscconn, s.volumeId, ui, volTags); err != nil {
			err := fmt.Errorf("error creating tags for volume: %w", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	// Wait for the volume to become ready
	err = osccommon.WaitUntilOscVolumeAvailable(ctx, oscconn, s.volumeId)
	if err != nil {
		err := fmt.Errorf("error waiting for volume: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("volume_id", s.volumeId)
	return multistep.ActionContinue
}

func (s *StepCreateVolume) Cleanup(state multistep.StateBag) {
	if s.volumeId == "" {
		return
	}

	oscconn := state.Get("osc").(*osccommon.OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	ui.Say("Deleting the created BSU volume...")
	_, err := oscconn.DeleteVolume(
		context.Background(),
		oscgo.DeleteVolumeRequest{VolumeId: s.volumeId},
	)
	if err != nil {
		ui.Error(fmt.Sprintf("Error deleting BSU volume: %s", err))
	}
}

func (s *StepCreateVolume) buildCreateVolumeInput(
	subregionName string,
	rootDevice *oscgo.BlockDeviceMappingImage,
) (*oscgo.CreateVolumeRequest, error) {
	if rootDevice == nil {
		return nil, errors.New("couldn't find root device")
	}

	createVolumeInput := &oscgo.CreateVolumeRequest{
		SubregionName: subregionName,
		Size:          rootDevice.Bsu.VolumeSize,
		SnapshotId:    rootDevice.Bsu.SnapshotId,
		VolumeType:    rootDevice.Bsu.VolumeType,
		Iops:          rootDevice.Bsu.Iops,
	}
	if int(s.RootVolumeSize) > ptr.From(rootDevice.Bsu.VolumeSize) {
		createVolumeInput.Size = new(int(s.RootVolumeSize))
	}

	if s.RootVolumeType == "" || oscgo.VolumeType(s.RootVolumeType) == ptr.From(rootDevice.Bsu.VolumeType) {
		return createVolumeInput, nil
	}

	if s.RootVolumeType == "io1" {
		return nil, fmt.Errorf(
			"root volume type cannot be io1, because existing root volume type was %s",
			ptr.From(rootDevice.Bsu.VolumeType),
		)
	}

	createVolumeInput.VolumeType = new(oscgo.VolumeType(s.RootVolumeType))
	// non io1 cannot set iops
	zeroIops := 0
	createVolumeInput.Iops = &zeroIops

	return createVolumeInput, nil
}
