package common

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// BlockDevice
type BlockDevice struct {
	DeleteOnVmDeletion bool   `mapstructure:"delete_on_vm_deletion"`
	DeviceName         string `mapstructure:"device_name"`
	IOPS               int    `mapstructure:"iops"`
	NoDevice           bool   `mapstructure:"no_device"`
	SnapshotId         string `mapstructure:"snapshot_id"`
	VolumeType         string `mapstructure:"volume_type"`
	VolumeSize         int64  `mapstructure:"volume_size"`
}

type BlockDevices struct {
	OMIBlockDevices    `mapstructure:",squash"`
	LaunchBlockDevices `mapstructure:",squash"`
}

type OMIBlockDevices struct {
	OMIMappings []BlockDevice `mapstructure:"omi_block_device_mappings"`
}

type LaunchBlockDevices struct {
	LaunchMappings []BlockDevice `mapstructure:"launch_block_device_mappings"`
}

func setBsuToCreate(blockDevice BlockDevice) *oscgo.BsuToCreate {
	defaultDeleteOnVmDeletion := true
	bsu := &oscgo.BsuToCreate{
		DeleteOnVmDeletion: &defaultDeleteOnVmDeletion,
	}
	if deleteOnVmDeletion := blockDevice.DeleteOnVmDeletion; !deleteOnVmDeletion {
		bsu.DeleteOnVmDeletion = &deleteOnVmDeletion
	}
	if volType := blockDevice.VolumeType; volType != "" {
		bsu.VolumeType = &volType
	}
	if volSize := int(blockDevice.VolumeSize); volSize > 0 {
		bsu.VolumeSize = &volSize
	}
	// IOPS is only valid for io1 type
	if blockDevice.VolumeType == "io1" {
		bsu.Iops = &blockDevice.IOPS
	}
	if snapId := blockDevice.SnapshotId; snapId != "" {
		bsu.SnapshotId = &snapId
	}

	return bsu
}

func buildOscBlockDevicesImage(b []BlockDevice) []oscgo.BlockDeviceMappingImage {
	var blockDevices []oscgo.BlockDeviceMappingImage
	for _, blockDevice := range b {
		mapping := oscgo.BlockDeviceMappingImage{}

		if deviceName := blockDevice.DeviceName; deviceName != "" {
			mapping.DeviceName = &deviceName
		}
		mapping.Bsu = setBsuToCreate(blockDevice)
		blockDevices = append(blockDevices, mapping)
	}
	return blockDevices
}

func buildOscBlockDevicesVmCreation(b []BlockDevice) []oscgo.BlockDeviceMappingVmCreation {
	var blockDevices []oscgo.BlockDeviceMappingVmCreation
	for _, blockDevice := range b {
		mapping := oscgo.BlockDeviceMappingVmCreation{}

		if deviceName := blockDevice.DeviceName; deviceName != "" {
			mapping.DeviceName = &deviceName
		}

		if blockDevice.NoDevice {
			mapping.NoDevice = aws.String("")
		} else {
			mapping.Bsu = setBsuToCreate(blockDevice)
		}
		blockDevices = append(blockDevices, mapping)
	}
	return blockDevices
}

func (b *BlockDevice) Prepare(ctx *interpolate.Context) error {
	if b.DeviceName == "" {
		return errors.New("the `device_name` must be specified " +
			"for every device in the block device mapping")
	}
	return nil
}

func (b *BlockDevices) Prepare(ctx *interpolate.Context) (errs []error) {
	for _, d := range b.OMIMappings {
		if err := d.Prepare(ctx); err != nil {
			errs = append(errs, fmt.Errorf("OMIMapping: %w", err))
		}
	}
	for _, d := range b.LaunchMappings {
		if err := d.Prepare(ctx); err != nil {
			errs = append(errs, fmt.Errorf("LaunchMapping: %w", err))
		}
	}
	return errs
}

func (b *OMIBlockDevices) BuildOscOMIDevices() []oscgo.BlockDeviceMappingImage {
	return buildOscBlockDevicesImage(b.OMIMappings)
}

func (b *LaunchBlockDevices) BuildOSCLaunchDevices() []oscgo.BlockDeviceMappingVmCreation {
	return buildOscBlockDevicesVmCreation(b.LaunchMappings)
}
