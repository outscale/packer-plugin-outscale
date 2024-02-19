package common

import (
	"fmt"
	"log"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/packer-plugin-outscale/builder/common/retry"
)

type stateRefreshFunc func() (string, error)

func waitUntilForOscVmRunning(conn *OscClient, vmID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "running", waitUntilOscVmStateFunc(conn, vmID))
	err := <-errCh
	return err
}

func waitUntilOscVmDeleted(conn *OscClient, vmID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "terminated", waitUntilOscVmStateFunc(conn, vmID))
	return <-errCh
}

func waitUntilOscVmStopped(conn *OscClient, vmID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "stopped", waitUntilOscVmStateFunc(conn, vmID))
	return <-errCh
}

func WaitUntilOscSnapshotCompleted(conn *OscClient, id string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "completed", waitUntilOscSnapshotStateFunc(conn, id))
	return <-errCh
}

func WaitUntilOscImageAvailable(conn *OscClient, imageID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "available", waitUntilOscImageStateFunc(conn, imageID))
	return <-errCh
}

func WaitUntilOscVolumeAvailable(conn *OscClient, volumeID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "available", volumeOscWaitFunc(conn, volumeID))
	return <-errCh
}

func WaitUntilOscVolumeIsLinked(conn *OscClient, volumeID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "attached", waitUntilOscVolumeLinkedStateFunc(conn, volumeID))
	return <-errCh
}

func WaitUntilOscVolumeIsUnlinked(conn *OscClient, volumeID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "dettached", waitUntilOscVolumeUnLinkedStateFunc(conn, volumeID))
	return <-errCh
}

func WaitUntilOscSnapshotDone(conn *OscClient, snapshotID string) error {
	errCh := make(chan error, 1)
	go waitForState(errCh, "completed", waitUntilOscSnapshotDoneStateFunc(conn, snapshotID))
	return <-errCh
}

func waitForState(errCh chan<- error, target string, refresh stateRefreshFunc) {
	err := retry.Run(2, 2, 0, func(_ uint) (bool, error) {
		state, err := refresh()
		if err != nil {
			return false, err
		} else if state == target {
			return true, nil
		}
		return false, nil
	})
	errCh <- err
}

func waitUntilOscVmStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Retrieving state for VM with id %s", id)
		resp, _, err := conn.Api.VmApi.ReadVms(conn.Auth).ReadVmsRequest(oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{id},
			},
		}).Execute()

		if err != nil {
			return "", err
		}

		//TODO: check if needed
		// if resp == nil {
		// 	return "", fmt.Errorf("Vm with ID %s not Found", id)
		// }

		if len(resp.GetVms()) == 0 {
			return "pending", nil
		}

		return *resp.GetVms()[0].State, nil
	}
}

func waitUntilOscVolumeLinkedStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Check if volume with id %s exists", id)
		resp, _, err := conn.Api.VolumeApi.ReadVolumes(conn.Auth).ReadVolumesRequest(oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds: &[]string{id},
			},
		}).Execute()

		if err != nil {
			return "", err
		}

		if len(resp.GetVolumes()) == 0 {
			return "pending", nil
		}

		if len(*resp.GetVolumes()[0].LinkedVolumes) == 0 {
			return "pending", nil
		}
		volume := resp.GetVolumes()[0]
		return volume.GetState(), nil
	}
}

func waitUntilOscVolumeUnLinkedStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Check if volume with id %s exists", id)
		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{id}},
		}
		resp, _, err := conn.Api.VolumeApi.ReadVolumes(conn.Auth).ReadVolumesRequest(request).Execute()
		if err != nil {
			return "", err
		}

		if len(resp.GetVolumes()) == 0 {
			return "pending", nil
		}

		if len(*resp.GetVolumes()[0].LinkedVolumes) == 0 {
			return "dettached", nil
		}

		return "failed", nil
	}
}

func waitUntilOscSnapshotStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Check if Snapshot with id %s exists", id)
		resp, _, err := conn.Api.SnapshotApi.ReadSnapshots(conn.Auth).ReadSnapshotsRequest(oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
		}).Execute()

		if err != nil {
			return "", err
		}

		if len(resp.GetSnapshots()) == 0 {
			return "pending", nil
		}

		return *resp.GetSnapshots()[0].State, nil
	}
}

func waitUntilOscImageStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Check if Image with id %s exists", id)
		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
		}
		resp, _, err := conn.Api.ImageApi.ReadImages(conn.Auth).ReadImagesRequest(filterReq).Execute()
		if err != nil {
			return "", err
		}

		if len(resp.GetImages()) == 0 {
			return "pending", nil
		}

		if *resp.GetImages()[0].State == "failed" {
			return *resp.GetImages()[0].State, fmt.Errorf("Image (%s) creation is failed", id)
		}

		return resp.GetImages()[0].GetState(), nil
	}
}

func waitUntilOscSnapshotDoneStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Check if Snapshot with id %s exists", id)
		resp, _, err := conn.Api.SnapshotApi.ReadSnapshots(conn.Auth).ReadSnapshotsRequest(oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
		}).Execute()

		if err != nil {
			return "", err
		}

		if len(resp.GetSnapshots()) == 0 {
			return "", fmt.Errorf("Snapshot with ID %s. Not Found", id)
		}

		if *resp.GetSnapshots()[0].State == "error" {
			return *resp.GetSnapshots()[0].State, fmt.Errorf("Snapshot (%s) creation is failed", id)
		}

		return *resp.GetSnapshots()[0].State, nil
	}
}

func volumeOscWaitFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
		log.Printf("[Debug] Check if SvolumeG with id %s exists", id)
		resp, _, err := conn.Api.VolumeApi.ReadVolumes(conn.Auth).ReadVolumesRequest(oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{id}},
		}).Execute()

		if err != nil {
			return "", err
		}

		if len(resp.GetVolumes()) == 0 {
			return "waiting", nil
		}

		if *resp.GetVolumes()[0].State == "error" {
			return *resp.GetVolumes()[0].State, fmt.Errorf("Volume (%s) creation is failed", id)
		}

		return *resp.GetVolumes()[0].State, nil
	}
}
