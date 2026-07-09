package common

import (
	"context"
	"errors"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/outscale/goutils/sdk/auth"
)

// StepCheckAuth validates that the authentication is valid
type StepCheckAuth struct{}

var checkCredentials = auth.CheckCredentials

func (s *StepCheckAuth) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	conn := state.Get("osc").(*OscClient)

	ui.Say("Validating configured credentials")

	err := checkCredentials(ctx, conn)
	if err != nil {
		state.Put("error", errors.New("OUTSCALE API authentication failed; verify your credentials and region configuration"))
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepCheckAuth) Cleanup(multistep.StateBag) {}
