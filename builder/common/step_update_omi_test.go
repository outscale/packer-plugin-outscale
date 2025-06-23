package common

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

// Create statebag for running test
func getStateUpdateOMI(omiId string, snapId string, state multistep.StateBag) multistep.StateBag {
	config := state.Get("accessConfig").(*AccessConfig)

	omis := make(map[string]string, 0)
	omis[config.RawRegion] = omiId
	state.Put("omis", omis)
	snaps := make(map[string][]string, 0)
	snaps[config.RawRegion] = []string{snapId}
	state.Put("snapshots", snaps)
	state.Put("accessConfig", config)

	return state
}

func TestUpdateOmi(t *testing.T) {
	state, err := getState()
	if err != nil {
		t.Fatalf("error setting osc client %s", err.Error())
	}
	config := state.Get("accessConfig").(*AccessConfig)
	client := state.Get("osc").(*OscClient)

	stepUpdateOMIAttributes := StepUpdateOMIAttributes{
		AccountIds:         []string{},
		SnapshotAccountIds: []string{},
		RawRegion:          config.RawRegion,
		GlobalPermission:   true,
	}

	createVmOpts := oscgo.CreateVmsRequest{
		ImageId: "ami-0aaec44f",
	}
	createVmOpts.SetVmType("tinav7.c4r8p2")
	createReq, _, err := client.Api.VmApi.CreateVms(client.Auth).CreateVmsRequest(createVmOpts).Execute()
	if err != nil {
		t.Fatalf("error creating VM: %v", err)
	}
	vm := createReq.GetVms()[0]
	vmId := vm.GetVmId()
	volId := vm.GetBlockDeviceMappings()[0].GetBsu().VolumeId

	// defer delete vm
	defer func() {
		_, _, _ = client.Api.VmApi.DeleteVms(client.Auth).
			DeleteVmsRequest(oscgo.DeleteVmsRequest{
				VmIds: []string{vmId},
			}).Execute()
	}()

	err = waitUntilForOscVmRunning(client, vmId)
	if err != nil {
		t.Fatalf("error waiting for VM to be running: %v", err)
	}

	// Create new image
	imgName := "ci-packer"
	bootModes := []oscgo.BootMode{"legacy", "uefi"}
	createRet, _, err := client.Api.ImageApi.CreateImage(client.Auth).
		CreateImageRequest(oscgo.CreateImageRequest{
			ImageName: &imgName,
			VmId:      &vmId,
			BootModes: &bootModes,
		}).Execute()
	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}
	image := createRet.Image
	imgId := image.GetImageId()

	// defer delete image
	defer func() {
		_, _, _ = client.Api.ImageApi.DeleteImage(client.Auth).
			DeleteImageRequest(oscgo.DeleteImageRequest{
				ImageId: imgId,
			}).Execute()
	}()

	if !reflect.DeepEqual(image.GetBootModes(), bootModes) {
		t.Fatalf("created OMI doesn't contains the bootmodes %v: %v", bootModes, err)
	}

	snapshotDesc := "snapshot-ci-packer"
	createSnap, _, err := client.Api.SnapshotApi.CreateSnapshot(client.Auth).
		CreateSnapshotRequest(oscgo.CreateSnapshotRequest{
			Description: &snapshotDesc,
			VolumeId:    volId,
		}).Execute()
	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}
	snapId := createSnap.Snapshot.GetSnapshotId()
	defer func() {
		_, _, _ = client.Api.SnapshotApi.DeleteSnapshot(client.Auth).
			DeleteSnapshotRequest(oscgo.DeleteSnapshotRequest{
				SnapshotId: snapId,
			}).Execute()
	}()

	state = getStateUpdateOMI(imgId, snapId, state)
	action := stepUpdateOMIAttributes.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", action)
	}
}
