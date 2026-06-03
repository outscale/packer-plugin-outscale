package chroot

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type mountPathData struct {
	Device string
}

// StepMountDevice mounts the attached device.
//
// Produces:
//
//	mount_path string - The location where the volume was mounted.
//	mount_device_cleanup CleanupFunc - To perform early cleanup
type StepMountDevice struct {
	MountOptions   []string
	MountPartition string

	mountPath string
}

func (s *StepMountDevice) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packersdk.Ui)
	device := state.Get("device").(string)
	if config.NVMEDevicePath != "" {
		// customizable device path for mounting NVME block devices on c5 and m5 HVM
		device = config.NVMEDevicePath
	}
	wrappedCommand := state.Get("wrappedCommand").(CommandWrapper)

	interCtx := config.ctx
	interCtx.Data = &mountPathData{Device: filepath.Base(device)}
	mountPath, err := interpolate.Render(config.MountPath, &interCtx)
	if err != nil {
		err := fmt.Errorf("error preparing mount directory: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	mountPath, err = filepath.Abs(mountPath)
	if err != nil {
		err := fmt.Errorf("error preparing mount directory: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("[DEBUG] Device: %s", device)
	log.Printf("[DEBUG] Mount path: %s", mountPath)

	if err := os.MkdirAll(mountPath, 0o750); err != nil {
		err := fmt.Errorf("error creating mount directory: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	deviceMount := device
	if s.MountPartition != "" && s.MountPartition != "0" {
		deviceMount += s.MountPartition
	}
	state.Put("deviceMount", deviceMount)

	log.Printf("[DEBUG] s.MountPartition  = %s", s.MountPartition)
	log.Printf("[DEBUG] DeviceMount: %s", deviceMount)

	ui.Say("Mounting the root device...")
	stderr := new(bytes.Buffer)

	// build mount options from mount_options config, useful for nouuid options
	// or other specific device type settings for mount
	opts := ""
	if len(s.MountOptions) > 0 {
		opts = "-o " + strings.Join(s.MountOptions, " -o ")
	}
	mountCommand, err := wrappedCommand(
		fmt.Sprintf("mount %s %s %s", opts, deviceMount, mountPath))
	if err != nil {
		err := fmt.Errorf("error creating mount command: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	log.Printf("[DEBUG] (step mount) mount command is %s", mountCommand)
	cmd := ShellCommand(ctx, mountCommand)
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		err := fmt.Errorf(
			"error mounting root volume: %w\nStderr: %s", err, stderr.String())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Set the mount path so we remember to unmount it later
	s.mountPath = mountPath
	state.Put("mount_path", s.mountPath)
	state.Put("mount_device_cleanup", s)

	return multistep.ActionContinue
}

func (s *StepMountDevice) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packersdk.Ui)
	if err := s.CleanupFunc(state); err != nil {
		ui.Error(err.Error())
	}
}

func (s *StepMountDevice) CleanupFunc(state multistep.StateBag) error {
	if s.mountPath == "" {
		return nil
	}

	ui := state.Get("ui").(packersdk.Ui)
	wrappedCommand := state.Get("wrappedCommand").(CommandWrapper)

	ui.Say("Unmounting the root device...")
	unmountCommand, err := wrappedCommand("umount " + s.mountPath)
	if err != nil {
		return fmt.Errorf("error creating unmount command: %w", err)
	}

	cmd := ShellCommand(context.Background(), unmountCommand)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error unmounting root device: %w", err)
	}

	s.mountPath = ""
	return nil
}
