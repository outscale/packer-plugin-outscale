package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

type StepPublicIp struct {
	AssociatePublicIpAddress bool
	Comm                     *communicator.Config
	publicIpId               string
	Debug                    bool

	doCleanup bool
}

func (s *StepPublicIp) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
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
	resp, err := conn.CreatePublicIp(ctx, allocOpts)
	if err != nil {
		state.Put("error", fmt.Errorf("error creating temporary PublicIp: %w", err))
		return multistep.ActionHalt
	}

	// From there, we have a Public Ip to destroy.
	s.doCleanup = true

	// Set some data for use in future steps
	s.publicIpId = *resp.PublicIp.PublicIpId
	state.Put("publicip_id", *resp.PublicIp.PublicIpId)

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
	ctx := context.Background()
	_, err := conn.DeletePublicIp(ctx, oscgo.DeletePublicIpRequest{
		PublicIpId: &s.publicIpId,
	})
	if err != nil {
		ui.Error(
			fmt.Sprintf(
				"Error cleaning up PublicIp. Please delete the PublicIp manually: %s",
				s.publicIpId,
			),
		)
	}
}
