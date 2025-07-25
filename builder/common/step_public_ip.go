package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

type StepPublicIp struct {
	AssociatePublicIpAddress bool
	Comm                     *communicator.Config
	publicIpId               string
	Debug                    bool

	doCleanup bool
}

func (s *StepPublicIp) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		ui   = state.Get("ui").(packersdk.Ui)
		conn = state.Get("osc").(*OscClient)
	)

	if !s.AssociatePublicIpAddress {

		// In this case, we are in the public Cloud, so we'll
		// not explicitely allocate a public IP.
		return multistep.ActionContinue
	}

	ui.Say("Creating temporary PublicIp for instance ")
	allocOpts := oscgo.CreatePublicIpRequest{}
	resp, _, err := conn.Api.PublicIpApi.CreatePublicIp(conn.Auth).CreatePublicIpRequest(allocOpts).Execute()

	if err != nil {
		state.Put("error", fmt.Errorf("error creating temporary PublicIp: %w", err))
		return multistep.ActionHalt
	}

	// From there, we have a Public Ip to destroy.
	s.doCleanup = true

	// Set some data for use in future steps
	s.publicIpId = *resp.GetPublicIp().PublicIpId
	state.Put("publicip_id", *resp.GetPublicIp().PublicIpId)

	return multistep.ActionContinue
}

func (s *StepPublicIp) Cleanup(state multistep.StateBag) {
	if !s.doCleanup {
		return
	}

	var (
		conn = state.Get("osc").(*OscClient)
		ui   = state.Get("ui").(packersdk.Ui)
	)

	// Remove the Public IP
	ui.Say("Deleting temporary PublicIp...")
	_, _, err := conn.Api.PublicIpApi.DeletePublicIp(conn.Auth).DeletePublicIpRequest(oscgo.DeletePublicIpRequest{
		PublicIpId: &s.publicIpId,
	}).Execute()

	if err != nil {
		ui.Error(fmt.Sprintf("Error cleaning up PublicIp. Please delete the PublicIp manually: %s", s.publicIpId))
	}
}
