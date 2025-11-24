package common

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"

	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
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
	ctx := context.Background()
	state, err := getState()
	if err != nil {
		t.Fatalf("error setting osc client %s", err.Error())
	}
	config := state.Get("accessConfig").(*AccessConfig)
	client := state.Get("osc").(*OscClient)

	stepSourceOMIInfo := mostRecentOmiFilterStep()
	imageAction := stepSourceOMIInfo.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if imageAction != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", imageAction)
	}
	sourceImage := state.Get("source_image").(oscgo.Image)

	stepUpdateOMIAttributes := StepUpdateOMIAttributes{
		AccountIds:         []string{},
		SnapshotAccountIds: []string{},
		RawRegion:          config.RawRegion,
		GlobalPermission:   true,
	}

	createVmOpts := oscgo.CreateVmsRequest{
		ImageId: sourceImage.ImageId,
	}
	defaultVmType := "tinav7.c4r8p2"
	createVmOpts.VmType = &defaultVmType
	createReq, err := client.CreateVms(ctx, createVmOpts)
	if err != nil {
		t.Fatalf("error creating VM: %v", err)
	}
	vm := (*createReq.Vms)[0]
	vmId := vm.VmId
	volId := vm.BlockDeviceMappings[0].Bsu.VolumeId

	// defer delete vm
	defer func() {
		_, _ = client.DeleteVms(ctx, oscgo.DeleteVmsRequest{
			VmIds: []string{vmId},
		})
	}()

	err = waitUntilForOscVmRunning(client, vmId)
	if err != nil {
		t.Fatalf("error waiting for VM to be running: %v", err)
	}

	// Create new image
	imgName := "ci-packer"
	bootModes := []oscgo.BootMode{"legacy", "uefi"}
	createRet, err := client.CreateImage(ctx, oscgo.CreateImageRequest{
		ImageName: &imgName,
		VmId:      &vmId,
		BootModes: &bootModes,
	})
	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}
	image := createRet.Image
	imgId := image.ImageId

	// defer delete image
	defer func() {
		_, _ = client.DeleteImage(ctx, oscgo.DeleteImageRequest{
			ImageId: imgId,
		})
	}()

	if !reflect.DeepEqual(image.BootModes, bootModes) {
		t.Fatalf("created OMI doesn't contains the bootmodes %v: %v", bootModes, err)
	}

	snapshotDesc := "snapshot-ci-packer"
	createSnap, err := client.CreateSnapshot(ctx, oscgo.CreateSnapshotRequest{
		Description: &snapshotDesc,
		VolumeId:    &volId,
	})
	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}
	snapId := createSnap.Snapshot.SnapshotId
	defer func() {
		_, _ = client.DeleteSnapshot(ctx, oscgo.DeleteSnapshotRequest{
			SnapshotId: snapId,
		})
	}()

	state = getStateUpdateOMI(imgId, snapId, state)
	action := stepUpdateOMIAttributes.Run(ctx, state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", action)
	}
}
