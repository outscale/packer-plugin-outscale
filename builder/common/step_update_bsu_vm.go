package common

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepUpdateBSUBackedVm struct {
	EnableAMIENASupport      *bool
	EnableAMISriovNetSupport bool
}

func (s *StepUpdateBSUBackedVm) Run(
	_ context.Context,
	state multistep.StateBag,
) multistep.StepAction {
	return multistep.ActionContinue
}

func (s *StepUpdateBSUBackedVm) Cleanup(state multistep.StateBag) {
	// No cleanup...
}
