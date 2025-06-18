package chroot

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
)

// StepVmInfo verifies that this builder is running on an Outscale vm.
type StepVmInfo struct{}

func (s *StepVmInfo) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	oscconn := state.Get("osc").(*osccommon.OscClient)
	//session := state.Get("clientConfig").(*session.Session)
	ui := state.Get("ui").(packersdk.Ui)

	// Get our own vm ID
	ui.Say("Gathering information about this Outscale vm...")

	cmd := ShellCommand("curl http://169.254.169.254/latest/meta-data/instance-id")

	vmID, err := cmd.Output()
	if err != nil {
		err := fmt.Errorf(
			"error retrieving the ID of the vm Packer is running on.\n" +
				"Please verify Packer is running on a proper Outscale vm")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	log.Printf("[Debug] VmID got: %s", string(vmID))

	// Query the entire vm metadata

	resp, _, err := oscconn.Api.VmApi.ReadVms(oscconn.Auth).ReadVmsRequest(oscgo.ReadVmsRequest{
		Filters: &oscgo.FiltersVm{
			VmIds: &[]string{string(vmID)},
		},
	}).Execute()
	if err != nil {
		err := fmt.Errorf("error getting vm data: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	vmsResp := resp.GetVms()

	if len(vmsResp) == 0 {
		err := errors.New("error getting vm data: no vm found")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	vm := vmsResp[0]
	state.Put("vm", vm)

	return multistep.ActionContinue
}

func (s *StepVmInfo) Cleanup(multistep.StateBag) {}
