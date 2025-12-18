package common

import (
	"context"
	"fmt"
	"log"

	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/packer-plugin-outscale/builder/common/retry"
)

func waitUntilForOscVmRunning(conn *OscClient, vmID string) error {
	return waitForState(oscgo.VmStateRunning, waitUntilOscVmStateFunc(conn, vmID))
}

func waitUntilOscVmDeleted(conn *OscClient, vmID string) error {
	return waitForState(oscgo.VmStateTerminated, waitUntilOscVmStateFunc(conn, vmID))
}

func waitUntilOscVmStopped(conn *OscClient, vmID string) error {
	return waitForState(oscgo.VmStateStopped, waitUntilOscVmStateFunc(conn, vmID))
}

func WaitUntilOscSnapshotCompleted(conn *OscClient, id string) error {
	return waitForState(oscgo.SnapshotStateCompleted, waitUntilOscSnapshotStateFunc(conn, id))
}

func WaitUntilOscImageAvailable(conn *OscClient, imageID string) error {
	return waitForState(oscgo.ImageStateAvailable, waitUntilOscImageStateFunc(conn, imageID))
}

func WaitUntilOscVolumeAvailable(conn *OscClient, volumeID string) error {
	return waitForState(oscgo.VolumeStateAvailable, volumeOscWaitFunc(conn, volumeID))
}

func WaitUntilOscVolumeIsLinked(conn *OscClient, volumeID string) error {
	return waitForState(
		oscgo.VolumeStateInUse,
		waitUntilOscVolumeLinkedStateFunc(conn, volumeID),
	)
}

func WaitUntilOscVolumeIsUnlinked(conn *OscClient, volumeID string) error {
	return waitForState(
		"dettached",
		waitUntilOscVolumeUnLinkedStateFunc(conn, volumeID),
	)
}

func WaitUntilOscSnapshotDone(conn *OscClient, snapshotID string) error {
	return waitForState(
		oscgo.SnapshotStateCompleted,
		waitUntilOscSnapshotDoneStateFunc(conn, snapshotID),
	)
}

func waitForState[T comparable](target T, refresh func() (T, error)) error {
	return retry.Run(2, 2, 0, func(_ uint) (bool, error) {
		state, err := refresh()
		if err != nil {
			return false, err
		} else if state == target {
			return true, nil
		}
		return false, nil
	})
}

func waitUntilOscVmStateFunc(conn *OscClient, id string) func() (oscgo.VmState, error) {
	return func() (oscgo.VmState, error) {
		ctx := context.Background()

		log.Printf("[Debug] Retrieving state for VM with id %s", id)
		resp, err := conn.ReadVms(ctx, oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{id},
			},
		})
		if err != nil {
			return "", err
		}

		// TODO: check if needed
		// if resp == nil {
		// 	return "", fmt.Errorf("Vm with ID %s not Found", id)
		// }

		if len(*resp.Vms) == 0 {
			return "pending", nil
		}

		return (*resp.Vms)[0].State, nil
	}
}

func waitUntilOscVolumeLinkedStateFunc(
	conn *OscClient,
	id string,
) func() (oscgo.VolumeState, error) {
	return func() (oscgo.VolumeState, error) {
		ctx := context.Background()

		log.Printf("[Debug] Check if volume with id %s exists", id)
		resp, err := conn.ReadVolumes(ctx, oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds: &[]string{id},
			},
		})
		if err != nil {
			return "", err
		}

		if len(*resp.Volumes) == 0 {
			return "pending", nil
		}

		if len(*(*resp.Volumes)[0].LinkedVolumes) == 0 {
			return "pending", nil
		}
		volume := (*resp.Volumes)[0]
		return volume.State, nil
	}
}

func waitUntilOscVolumeUnLinkedStateFunc(conn *OscClient, id string) func() (string, error) {
	return func() (string, error) {
		ctx := context.Background()
		log.Printf("[Debug] Check if volume with id %s exists", id)
		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{id}},
		}
		resp, err := conn.ReadVolumes(ctx, request)
		if err != nil {
			return "", err
		}

		if len(*resp.Volumes) == 0 {
			return "pending", nil
		}

		if len(*(*resp.Volumes)[0].LinkedVolumes) == 0 {
			return "dettached", nil
		}

		return "failed", nil
	}
}

func waitUntilOscSnapshotStateFunc(conn *OscClient, id string) func() (oscgo.SnapshotState, error) {
	return func() (oscgo.SnapshotState, error) {
		ctx := context.Background()

		log.Printf("[Debug] Check if Snapshot with id %s exists", id)
		resp, err := conn.ReadSnapshots(ctx, oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
		})
		if err != nil {
			return "", err
		}

		if len(*resp.Snapshots) == 0 {
			return "pending", nil
		}

		return (*resp.Snapshots)[0].State, nil
	}
}

func waitUntilOscImageStateFunc(conn *OscClient, id string) func() (oscgo.ImageState, error) {
	return func() (oscgo.ImageState, error) {
		ctx := context.Background()
		log.Printf("[Debug] Check if Image with id %s exists", id)
		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
		}
		resp, err := conn.ReadImages(ctx, filterReq)
		if err != nil {
			return "", err
		}

		if len(*resp.Images) == 0 {
			return "pending", nil
		}

		if (*resp.Images)[0].State == "failed" {
			return (*resp.Images)[0].State, fmt.Errorf("image (%s) creation is failed", id)
		}

		return (*resp.Images)[0].State, nil
	}
}

func waitUntilOscSnapshotDoneStateFunc(
	conn *OscClient,
	id string,
) func() (oscgo.SnapshotState, error) {
	return func() (oscgo.SnapshotState, error) {
		ctx := context.Background()
		log.Printf("[Debug] Check if Snapshot with id %s exists", id)
		resp, err := conn.ReadSnapshots(ctx, oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
		})
		if err != nil {
			return "", err
		}

		if len(*resp.Snapshots) == 0 {
			return "", fmt.Errorf("snapshot with ID %s. Not Found", id)
		}

		if (*resp.Snapshots)[0].State == "error" {
			return (*resp.Snapshots)[0].State, fmt.Errorf("snapshot (%s) creation is failed", id)
		}

		return (*resp.Snapshots)[0].State, nil
	}
}

func volumeOscWaitFunc(conn *OscClient, id string) func() (oscgo.VolumeState, error) {
	return func() (oscgo.VolumeState, error) {
		ctx := context.Background()
		log.Printf("[Debug] Check if SvolumeG with id %s exists", id)
		resp, err := conn.ReadVolumes(ctx, oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{id}},
		})
		if err != nil {
			return "", err
		}

		if len(*resp.Volumes) == 0 {
			return "waiting", nil
		}

		if (*resp.Volumes)[0].State == "error" {
			return (*resp.Volumes)[0].State, fmt.Errorf("volume (%s) creation is failed", id)
		}

		return (*resp.Volumes)[0].State, nil
	}
}
