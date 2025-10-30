package common

import (
	"context"
	"fmt"
	"log"

	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
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

func waitUntilOscVolumeLinkedStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
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
		return *volume.State, nil
	}
}

func waitUntilOscVolumeUnLinkedStateFunc(conn *OscClient, id string) stateRefreshFunc {
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
			return "pending", nil
		}

		return "failed", nil
	}
}

func waitUntilOscSnapshotStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
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

func waitUntilOscImageStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
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

func waitUntilOscSnapshotDoneStateFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
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

func volumeOscWaitFunc(conn *OscClient, id string) stateRefreshFunc {
	return func() (string, error) {
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

		if *(*resp.Volumes)[0].State == "error" {
			return *(*resp.Volumes)[0].State, fmt.Errorf("volume (%s) creation is failed", id)
		}

		return *(*resp.Volumes)[0].State, nil
	}
}
