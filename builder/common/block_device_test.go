package common_test

import (
	"reflect"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/packer-plugin-outscale/builder/common"
)

func TestBlockDevice_LaunchDevices(t *testing.T) {
	cases := []struct {
		Config *common.BlockDevice
		Result oscgo.BlockDeviceMappingVmCreation
	}{
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				SnapshotId:         "snap-1234",
				VolumeType:         "standard",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					SnapshotId:         new("snap-1234"),
					VolumeType:         new(oscgo.VolumeTypeStandard),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName: "/dev/sdb",
				VolumeSize: 8,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(false),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "io1",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
				IOPS:               1000,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeIo1),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
					Iops:               new(1000),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeGp2),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeGp2),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "standard",
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeStandard),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
	}

	for _, tc := range cases {
		launchBlockDevices := common.LaunchBlockDevices{
			LaunchMappings: []common.BlockDevice{*tc.Config},
		}

		expected := []oscgo.BlockDeviceMappingVmCreation{tc.Result}

		launchResults := launchBlockDevices.BuildOSCLaunchDevices()
		if !reflect.DeepEqual(expected, launchResults) {
			t.Fatalf("Bad block device, \nexpected: %#v\n\ngot: %#v",
				expected, launchResults)
		}
	}
}

func TestBlockDevice_OMI(t *testing.T) {
	cases := []struct {
		Config *common.BlockDevice
		Result oscgo.BlockDeviceMappingImage
	}{
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				SnapshotId:         "snap-1234",
				VolumeType:         "standard",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					SnapshotId:         new("snap-1234"),
					VolumeType:         new(oscgo.VolumeTypeStandard),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "io1",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
				IOPS:               1000,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeIo1),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
					Iops:               new(1000),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeGp2),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeGp2),
					VolumeSize:         new(8),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
		{
			Config: &common.BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "standard",
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: new("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         new(oscgo.VolumeTypeStandard),
					DeleteOnVmDeletion: new(true),
				},
			},
		},
	}

	for i, tc := range cases {
		omiBlockDevices := common.OMIBlockDevices{
			OMIMappings: []common.BlockDevice{*tc.Config},
		}

		expected := []oscgo.BlockDeviceMappingImage{tc.Result}

		omiResults := omiBlockDevices.BuildOscOMIDevices()
		if !reflect.DeepEqual(expected, omiResults) {
			t.Fatalf("%d - Bad block device, \nexpected: %+#v\n\ngot: %+#v",
				i, expected, omiResults)
		}
	}
}
