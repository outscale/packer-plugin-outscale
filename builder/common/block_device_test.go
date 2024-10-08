package common

import (
	"reflect"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func TestBlockDevice_LaunchDevices(t *testing.T) {
	cases := []struct {
		Config *BlockDevice
		Result oscgo.BlockDeviceMappingVmCreation
	}{
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				SnapshotId:         "snap-1234",
				VolumeType:         "standard",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					SnapshotId:         oscgo.PtrString("snap-1234"),
					VolumeType:         oscgo.PtrString("standard"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName: "/dev/sdb",
				VolumeSize: 8,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(false),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "io1",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
				IOPS:               1000,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("io1"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
					Iops:               oscgo.PtrInt32(1000),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("gp2"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("gp2"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "standard",
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("standard"),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
	}

	for _, tc := range cases {

		launchBlockDevices := LaunchBlockDevices{
			LaunchMappings: []BlockDevice{*tc.Config},
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
		Config *BlockDevice
		Result oscgo.BlockDeviceMappingImage
	}{
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				SnapshotId:         "snap-1234",
				VolumeType:         "standard",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					SnapshotId:         oscgo.PtrString("snap-1234"),
					VolumeType:         oscgo.PtrString("standard"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "io1",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
				IOPS:               1000,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("io1"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
					Iops:               oscgo.PtrInt32(1000),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("gp2"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "gp2",
				VolumeSize:         8,
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("gp2"),
					VolumeSize:         oscgo.PtrInt32(8),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName:         "/dev/sdb",
				VolumeType:         "standard",
				DeleteOnVmDeletion: true,
			},

			Result: oscgo.BlockDeviceMappingImage{
				DeviceName: oscgo.PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         oscgo.PtrString("standard"),
					DeleteOnVmDeletion: oscgo.PtrBool(true),
				},
			},
		},
	}

	for i, tc := range cases {
		omiBlockDevices := OMIBlockDevices{
			OMIMappings: []BlockDevice{*tc.Config},
		}

		expected := []oscgo.BlockDeviceMappingImage{tc.Result}

		omiResults := omiBlockDevices.BuildOscOMIDevices()
		if !reflect.DeepEqual(expected, omiResults) {
			t.Fatalf("%d - Bad block device, \nexpected: %+#v\n\ngot: %+#v",
				i, expected, omiResults)
		}
	}
}
