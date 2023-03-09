package common

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/outscale/osc-sdk-go/v2"
	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/outscale/packer-plugin-outscale/builder/osc/common/retry"
)

const (
	RunSourceVmBSUExpectedRootDevice = "bsu"
)

type StepRunSourceVm struct {
	BlockDevices                BlockDevices
	Comm                        *communicator.Config
	Ctx                         interpolate.Context
	Debug                       bool
	BsuOptimized                bool
	EnableT2Unlimited           bool
	ExpectedRootDevice          string
	IamVmProfile                string
	VmInitiatedShutdownBehavior string
	VmType                      string
	IsRestricted                bool
	SourceOMI                   string
	Tags                        TagMap
	UserData                    string
	UserDataFile                string
	VolumeTags                  TagMap
	RawRegion                   string

	vmId string
}

func (s *StepRunSourceVm) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	oscconn := state.Get("osc").(*OscClient)
	securityGroupIds := state.Get("securityGroupIds").([]string)
	ui := state.Get("ui").(packersdk.Ui)

	userData := s.UserData
	if s.UserDataFile != "" {
		contents, err := ioutil.ReadFile(s.UserDataFile)
		if err != nil {
			state.Put("error", fmt.Errorf("Problem reading user data file: %s", err))
			return multistep.ActionHalt
		}

		userData = string(contents)
	}

	// Test if it is encoded already, and if not, encode it
	if _, err := base64.StdEncoding.DecodeString(userData); err != nil {
		log.Printf("[DEBUG] base64 encoding user data...")
		userData = base64.StdEncoding.EncodeToString([]byte(userData))
	}

	ui.Say("Launching a source OUTSCALE vm...")
	image, ok := state.Get("source_image").(oscgo.Image)
	if !ok {
		state.Put("error", fmt.Errorf("source_image type assertion failed"))
		return multistep.ActionHalt
	}
	s.SourceOMI = *image.ImageId

	if s.ExpectedRootDevice != "" && *image.RootDeviceType != s.ExpectedRootDevice {
		state.Put("error", fmt.Errorf(
			"The provided source OMI has an invalid root device type.\n"+
				"Expected '%s', got '%s'.",
			s.ExpectedRootDevice, image.GetRootDeviceType()))
		return multistep.ActionHalt
	}

	var vmId string

	ui.Say("Adding tags to source vm")
	if _, exists := s.Tags["Name"]; !exists {
		s.Tags["Name"] = "Packer Builder"
	}

	rawRegion := s.RawRegion

	oscTags, err := s.Tags.OSCTags(s.Ctx, rawRegion, state)
	if err != nil {
		err := fmt.Errorf("Error tagging source vm: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	volTags, err := s.VolumeTags.OSCTags(s.Ctx, rawRegion, state)
	if err != nil {
		err := fmt.Errorf("Error tagging volumes: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	blockDevice := s.BlockDevices.BuildOSCLaunchDevices()
	vmcount := int32(1)
	subregion := state.Get("subregion_name").(string)
	subnetID := state.Get("subnet_id").(string)

	runOpts := oscgo.CreateVmsRequest{
		ImageId:             s.SourceOMI,
		VmType:              &s.VmType,
		UserData:            &userData,
		MaxVmsCount:         &vmcount,
		MinVmsCount:         &vmcount,
		BsuOptimized:        &s.BsuOptimized,
		BlockDeviceMappings: &blockDevice,
	}

	log.Printf("subregion is %s", subregion)
	if subregion != "" {
		runOpts.Placement = &oscgo.Placement{SubregionName: &subregion}
	}

	if s.Comm.SSHKeyPairName != "" {
		runOpts.KeypairName = &s.Comm.SSHKeyPairName
	}

	runOpts.SubnetId = &subnetID
	runOpts.SecurityGroupIds = &securityGroupIds

	if s.ExpectedRootDevice == "bsu" {
		runOpts.VmInitiatedShutdownBehavior = &s.VmInitiatedShutdownBehavior
	}

	runResp, _, err := oscconn.Api.VmApi.CreateVms(oscconn.Auth).CreateVmsRequest(runOpts).Execute()

	if err != nil {
		err := fmt.Errorf("Error launching source vm: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	vmId = *runResp.GetVms()[0].VmId
	volumeId := *runResp.GetVms()[0].GetBlockDeviceMappings()[0].Bsu.VolumeId

	// Set the vm ID so that the cleanup works properly
	s.vmId = vmId

	ui.Message(fmt.Sprintf("Vm ID: %s", vmId))
	ui.Say(fmt.Sprintf("Waiting for vm (%v) to become ready...", vmId))

	request := oscgo.ReadVmsRequest{
		Filters: &oscgo.FiltersVm{
			VmIds: &[]string{vmId},
		},
	}
	if err := waitUntilForOscVmRunning(oscconn, vmId); err != nil {
		err := fmt.Errorf("Error waiting for vm (%s) to become ready: %s", vmId, err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	//Set Vm tags and vollume tags
	if len(oscTags) > 0 {
		if err := CreateOSCTags(oscconn, s.vmId, ui, oscTags); err != nil {
			err := fmt.Errorf("Error creating tags for vm (%s): %s", s.vmId, err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	if len(volTags) > 0 {
		if err := CreateOSCTags(oscconn, volumeId, ui, volTags); err != nil {
			err := fmt.Errorf("Error creating tags for volume (%s): %s", volumeId, err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	if publicip_id, ok := state.Get("publicip_id").(string); ok {
		ui.Say(fmt.Sprintf("Linking temporary PublicIp %s to instance %s", publicip_id, vmId))
		request := oscgo.LinkPublicIpRequest{PublicIpId: &publicip_id, VmId: &vmId}
		_, _, err := oscconn.Api.PublicIpApi.LinkPublicIp(oscconn.Auth).LinkPublicIpRequest(request).Execute()
		if err != nil {
			state.Put("error", fmt.Errorf("Error linking PublicIp to VM: %s", err))
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	resp, _, err := oscconn.Api.VmApi.ReadVms(oscconn.Auth).ReadVmsRequest(request).Execute()
	r := resp

	if err != nil || len(r.GetVms()) == 0 {
		err := fmt.Errorf("Error finding source vm.")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	vm := r.GetVms()[0]

	if s.Debug {
		if vm.PublicDnsName != nil {
			ui.Message(fmt.Sprintf("Public DNS: %s", vm.GetPublicDnsName()))
		}

		if vm.PublicIp != nil {
			ui.Message(fmt.Sprintf("Public IP: %s", vm.GetPublicIp()))
		}

		if vm.PrivateIp != nil {
			ui.Message(fmt.Sprintf("Private IP: %s", vm.GetPrivateIp()))
		}
	}

	state.Put("vm", vm)
	// instance_id is the generic term used so that users can have access to the
	// instance id inside of the provisioners, used in step_provision.
	state.Put("instance_id", vmId)

	// If we're in a region that doesn't support tagging on vm creation,
	// do that now.

	if s.IsRestricted {
		oscTags.Report(ui)
		// Retry creating tags for about 2.5 minutes
		err = retry.Run(0.2, 30, 11, func(_ uint) (bool, error) {
			request := oscgo.CreateTagsRequest{
				Tags:        oscTags,
				ResourceIds: []string{vmId},
			}
			_, _, err := oscconn.Api.TagApi.CreateTags(oscconn.Auth).CreateTagsRequest(request).Execute()
			if err == nil {
				return true, nil
			}
			//TODO: improve error
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "InvalidVmID.NotFound" {
					return false, nil
				}
			}
			return true, err
		})

		if err != nil {
			err := fmt.Errorf("Error tagging source vm: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}

		// Now tag volumes

		volumeIds := make([]string, 0)
		for _, v := range *vm.BlockDeviceMappings {
			if bsu := v.Bsu; !reflect.DeepEqual(bsu, osc.BsuCreated{}) {
				volumeIds = append(volumeIds, *bsu.VolumeId)
			}
		}

		if len(volumeIds) > 0 && s.VolumeTags.IsSet() {
			ui.Say("Adding tags to source BSU Volumes")

			volumeTags, err := s.VolumeTags.OSCTags(s.Ctx, rawRegion, state)
			if err != nil {
				err := fmt.Errorf("Error tagging source BSU Volumes on %s: %s", vm.GetVmId(), err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}
			volumeTags.Report(ui)
			request := oscgo.CreateTagsRequest{
				ResourceIds: volumeIds,
				Tags:        volumeTags,
			}
			_, _, err = oscconn.Api.TagApi.CreateTags(oscconn.Auth).CreateTagsRequest(request).Execute()

			if err != nil {
				err := fmt.Errorf("Error tagging source BSU Volumes on %s: %s", vm.GetVmId(), err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}
		}
	}

	return multistep.ActionContinue
}

func (s *StepRunSourceVm) Cleanup(state multistep.StateBag) {
	oscconn := state.Get("osc").(*OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	// Terminate the source vm if it exists
	if s.vmId != "" {
		ui.Say("Terminating the source OUTSCALE vm...")
		_, _, err := oscconn.Api.VmApi.DeleteVms(oscconn.Auth).DeleteVmsRequest(oscgo.DeleteVmsRequest{
			VmIds: []string{s.vmId},
		}).Execute()
		if err != nil {
			ui.Error(fmt.Sprintf("Error terminating vm, may still be around: %s", err))
			return
		}

		if err := waitUntilOscVmDeleted(oscconn, s.vmId); err != nil {
			ui.Error(err.Error())
		}
	}
}
