package common

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"

	"bytes"
	"context"

	"os"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

// Create statebag for running test
func tState(omiId string, snapId string) (multistep.StateBag, error) {
	region := os.Getenv("OSC_REGION")
	state := new(multistep.BasicStateBag)
	accessConfig := &AccessConfig{
		RawRegion: region,
	}

	oscConn, err := accessConfig.NewOSCClient()
	if err != nil {
		return nil, err
	}
	state.Put("osc", oscConn)
	state.Put("ui", &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	})

	omis := make(map[string]string, 0)
	omis[region] = omiId
	state.Put("omis", omis)
	snaps := make(map[string][]string, 0)
	snaps[region] = []string{snapId}
	state.Put("snapshots", snaps)
	state.Put("accessConfig", accessConfig)
	return state, err
}

func TestUpdateOmi(t *testing.T) {
	region := os.Getenv("OSC_REGION")
	stepUpdateOMIAttributes := StepUpdateOMIAttributes{
		AccountIds:         []string{},
		SnapshotAccountIds: []string{},
		RawRegion:          region,
		GlobalPermission:   true,
	}
	config := &AccessConfig{
		RawRegion: region,
	}

	// Read vms to get 1 to copy
	regionconn, _ := config.NewOSCClient()
	readRet, _, err := regionconn.Api.VmApi.ReadVms(regionconn.Auth).ReadVmsRequest(oscgo.ReadVmsRequest{}).Execute()

	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	vm := readRet.GetVms()[0]
	vmId := vm.GetVmId()
	volId := vm.GetBlockDeviceMappings()[0].GetBsu().VolumeId

	// Create new image
	imgName := "ci-packer"
	createRet, _, err := regionconn.Api.ImageApi.CreateImage(regionconn.Auth).
		CreateImageRequest(oscgo.CreateImageRequest{
			ImageName: &imgName,
			VmId:      &vmId,
		}).Execute()

	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	imgId := createRet.Image.GetImageId()

	// defer delete image
	defer func() {
		_, _, _ = regionconn.Api.ImageApi.DeleteImage(regionconn.Auth).
			DeleteImageRequest(oscgo.DeleteImageRequest{
				ImageId: imgId,
			}).Execute()
	}()

	snapshotDesc := "snapshot-ci-packer"
	createSnap, _, err := regionconn.Api.SnapshotApi.CreateSnapshot(regionconn.Auth).
		CreateSnapshotRequest(oscgo.CreateSnapshotRequest{
			Description: &snapshotDesc,
			VolumeId:    volId,
		}).Execute()
	if err != nil {
		t.Fatalf("should not error, but: %v", err)
	}
	snapId := createSnap.Snapshot.GetSnapshotId()
	defer func() {
		_, _, _ = regionconn.Api.SnapshotApi.DeleteSnapshot(regionconn.Auth).
			DeleteSnapshotRequest(oscgo.DeleteSnapshotRequest{
				SnapshotId: snapId,
			}).Execute()
	}()

	state, err := tState(imgId, snapId)
	if state == nil {
		t.Fatalf("error retrieving state %s", err.Error())
	}

	action := stepUpdateOMIAttributes.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", action)
	}

}
