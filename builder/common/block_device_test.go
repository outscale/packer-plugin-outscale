package common

import (
	"reflect"
	"testing"

	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func PtrBool(boolVal bool) *bool {
	return &boolVal
}

func PtrString(strVal string) *string {
	return &strVal
}

func PtrInt(intVal int) *int {
	return &intVal
}

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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					SnapshotId:         PtrString("snap-1234"),
					VolumeType:         PtrString("standard"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
				},
			},
		},
		{
			Config: &BlockDevice{
				DeviceName: "/dev/sdb",
				VolumeSize: 8,
			},

			Result: oscgo.BlockDeviceMappingVmCreation{
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(false),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("io1"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
					Iops:               PtrInt(1000),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("gp2"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("gp2"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("standard"),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					SnapshotId:         PtrString("snap-1234"),
					VolumeType:         PtrString("standard"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("io1"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
					Iops:               PtrInt(1000),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("gp2"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("gp2"),
					VolumeSize:         PtrInt(8),
					DeleteOnVmDeletion: PtrBool(true),
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
				DeviceName: PtrString("/dev/sdb"),
				Bsu: &oscgo.BsuToCreate{
					VolumeType:         PtrString("standard"),
					DeleteOnVmDeletion: PtrBool(true),
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
