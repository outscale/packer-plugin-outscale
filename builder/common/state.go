package common

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/outscale/goutils/sdk/ptr"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/packer-plugin-outscale/builder/common/retry"
)

func WaitUntilForOscVmRunning(ctx context.Context, conn *OscClient, vmID string) error {
	return waitForState(waitUntilOscVmStateFunc(ctx, conn, vmID), oscgo.VmStateRunning)
}

func waitUntilOscVmDeleted(ctx context.Context, conn *OscClient, vmID string) error {
	return waitForState(waitUntilOscVmStateFunc(ctx, conn, vmID), oscgo.VmStateTerminated)
}

func waitUntilOscVmStopped(ctx context.Context, conn *OscClient, vmID string) error {
	return waitForState(waitUntilOscVmStateFunc(ctx, conn, vmID), oscgo.VmStateStopped)
}

func WaitUntilOscSnapshotCompleted(ctx context.Context, conn *OscClient, id string) error {
	return waitForState(waitUntilOscSnapshotStateFunc(ctx, conn, id), oscgo.SnapshotStateCompleted, oscgo.SnapshotStateError)
}

func WaitUntilOscImageAvailable(ctx context.Context, conn *OscClient, imageID string) error {
	return waitForState(waitUntilOscImageStateFunc(ctx, conn, imageID), oscgo.ImageStateAvailable, oscgo.ImageStateFailed)
}

func WaitUntilOscVolumeAvailable(ctx context.Context, conn *OscClient, volumeID string) error {
	return waitForState(volumeOscWaitFunc(ctx, conn, volumeID), oscgo.VolumeStateAvailable)
}

func WaitUntilOscVolumeIsLinked(ctx context.Context, conn *OscClient, volumeID string) error {
	return waitForState(
		waitUntilOscVolumeLinkedStateFunc(ctx, conn, volumeID),
		oscgo.VolumeStateInUse,
		oscgo.VolumeStateError,
	)
}

func WaitUntilOscVolumeIsUnlinked(ctx context.Context, conn *OscClient, volumeID string) error {
	return waitForState(
		waitUntilOscVolumeUnLinkedStateFunc(ctx, conn, volumeID),
		oscgo.VolumeStateAvailable,
		oscgo.VolumeStateError,
	)
}

func WaitUntilOscSnapshotDone(ctx context.Context, conn *OscClient, snapshotID string) error {
	return waitForState(
		waitUntilOscSnapshotDoneStateFunc(ctx, conn, snapshotID),
		oscgo.SnapshotStateCompleted,
		oscgo.SnapshotStateError,
	)
}

func waitForState[T ~string](refresh func() (T, error), target T, failStates ...T) error {
	return retry.Run(2, 2, 0, func(_ uint) (bool, error) {
		state, err := refresh()
		if err != nil {
			return false, err
		} else if state == target {
			return true, nil
		}
		if slices.Contains(failStates, state) {
			return false, fmt.Errorf("resource in failing state: %s", state)
		}

		return false, nil
	})
}

func waitUntilOscVmStateFunc(ctx context.Context, conn *OscClient, id string) func() (oscgo.VmState, error) {
	return func() (oscgo.VmState, error) {
		log.Printf("[DEBUG] Retrieving state for VM with id %s", id)
		resp, err := conn.ReadVms(ctx, oscgo.ReadVmsRequest{
			Filters: &oscgo.FiltersVm{
				VmIds: &[]string{id},
			},
		})
		if err != nil {
			return "", err
		}
		vms := ptr.From(resp.Vms)

		if len(vms) == 0 {
			return "", ErrNotFound
		}

		return vms[0].State, nil
	}
}

func waitUntilOscVolumeLinkedStateFunc(ctx context.Context, conn *OscClient, id string) func() (oscgo.VolumeState, error) {
	return func() (oscgo.VolumeState, error) {
		log.Printf("[DEBUG] Check if volume with id %s exists", id)
		resp, err := conn.ReadVolumes(ctx, oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds: &[]string{id},
			},
		})
		if err != nil {
			return "", err
		}
		volumes := ptr.From(resp.Volumes)

		if len(volumes) == 0 {
			return "", ErrNotFound
		}

		// We want the volume to be 'in-use' and all linked volumes to be 'attached'
		volume := volumes[0]
		for _, link := range volume.LinkedVolumes {
			if link.State != oscgo.LinkedVolumeStateAttached {
				return oscgo.VolumeStateCreating, nil
			}
		}

		return volume.State, nil
	}
}

func waitUntilOscVolumeUnLinkedStateFunc(ctx context.Context, conn *OscClient, id string) func() (oscgo.VolumeState, error) {
	return func() (oscgo.VolumeState, error) {
		log.Printf("[DEBUG] Check if volume with id %s exists", id)
		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{id}},
		}
		resp, err := conn.ReadVolumes(ctx, request)
		if err != nil {
			return "", err
		}
		volumes := ptr.From(resp.Volumes)

		if len(volumes) == 0 {
			return "", ErrNotFound
		}

		return volumes[0].State, nil
	}
}

func waitUntilOscSnapshotStateFunc(ctx context.Context, conn *OscClient, id string) func() (oscgo.SnapshotState, error) {
	return func() (oscgo.SnapshotState, error) {
		log.Printf("[DEBUG] Check if Snapshot with id %s exists", id)
		resp, err := conn.ReadSnapshots(ctx, oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
		})
		if err != nil {
			return "", err
		}
		snapshots := ptr.From(resp.Snapshots)

		if len(snapshots) == 0 {
			return "", ErrNotFound
		}

		return snapshots[0].State, nil
	}
}

func waitUntilOscImageStateFunc(ctx context.Context, conn *OscClient, id string) func() (oscgo.ImageState, error) {
	return func() (oscgo.ImageState, error) {
		log.Printf("[DEBUG] Check if Image with id %s exists", id)
		filterReq := oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{ImageIds: &[]string{id}},
		}
		resp, err := conn.ReadImages(ctx, filterReq)
		if err != nil {
			return "", err
		}
		images := ptr.From(resp.Images)

		if len(images) == 0 {
			return "", ErrNotFound
		}

		return images[0].State, nil
	}
}

func waitUntilOscSnapshotDoneStateFunc(
	ctx context.Context,
	conn *OscClient,
	id string,
) func() (oscgo.SnapshotState, error) {
	return func() (oscgo.SnapshotState, error) {
		log.Printf("[DEBUG] Check if Snapshot with id %s exists", id)
		resp, err := conn.ReadSnapshots(ctx, oscgo.ReadSnapshotsRequest{
			Filters: &oscgo.FiltersSnapshot{SnapshotIds: &[]string{id}},
		})
		if err != nil {
			return "", err
		}
		snapshots := ptr.From(resp.Snapshots)

		if len(snapshots) == 0 {
			return "", ErrNotFound
		}

		return snapshots[0].State, nil
	}
}

func volumeOscWaitFunc(ctx context.Context, conn *OscClient, id string) func() (oscgo.VolumeState, error) {
	return func() (oscgo.VolumeState, error) {
		log.Printf("[DEBUG] Check if volume with id %s exists", id)
		resp, err := conn.ReadVolumes(ctx, oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{id}},
		})
		if err != nil {
			return "", err
		}
		volumes := ptr.From(resp.Volumes)

		if len(volumes) == 0 {
			return "", ErrNotFound
		}

		return volumes[0].State, nil
	}
}
