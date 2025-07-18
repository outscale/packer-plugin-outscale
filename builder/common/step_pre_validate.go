package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

// StepPreValidate provides an opportunity to pre-validate any configuration for
// the build before actually doing any time consuming work
type StepPreValidate struct {
	DestOmiName     string
	ForceDeregister bool
	API             string
}

func (s *StepPreValidate) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	if s.ForceDeregister {
		ui.Say("Force Deregister flag found, skipping prevalidating OMI Name")
		return multistep.ActionContinue
	}

	var (
		conn   = state.Get("osc").(*OscClient)
		images []interface{}
	)

	ui.Say(fmt.Sprintf("Prevalidating OMI Name: %s", s.DestOmiName))

	accountResp, _, err := conn.Api.AccountApi.ReadAccounts(conn.Auth).ReadAccountsRequest(oscgo.ReadAccountsRequest{}).Execute()
	if err != nil || len(accountResp.GetAccounts()) == 0 {
		err := fmt.Errorf("error querying outscale account: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	resp, _, err := conn.Api.ImageApi.ReadImages(conn.Auth).ReadImagesRequest(oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{
			ImageNames: &[]string{s.DestOmiName},
			AccountIds: &[]string{accountResp.GetAccounts()[0].GetAccountId()},
		},
	}).Execute()

	if err != nil {
		err := fmt.Errorf("error querying OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	for _, omi := range resp.GetImages() {
		if omi.GetImageName() == s.DestOmiName {
			images = append(images, omi)
		}
	}

	if len(images) > 0 {
		err := fmt.Errorf("error: name conflicts with an existing OMI: %s", s.DestOmiName)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *StepPreValidate) Cleanup(multistep.StateBag) {}
