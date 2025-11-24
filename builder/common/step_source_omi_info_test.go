package common

import (
	"bytes"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/template/config"

	"context"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// Create statebag for running test
func getState() (multistep.StateBag, error) {
	state := new(multistep.BasicStateBag)
	accessConfig := &AccessConfig{}
	oscConn, err := accessConfig.NewOSCClient()
	if err != nil {
		return nil, err
	}
	state.Put("osc", oscConn)
	state.Put("ui", &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	})
	state.Put("accessConfig", accessConfig)
	return state, err
}

func mostRecentOmiFilterStep() StepSourceOMIInfo {
	stepSourceOMIInfo := StepSourceOMIInfo{
		OmiFilters: OmiFilterOptions{
			Owners: []string{"Outscale"},
			NameValueFilter: config.NameValueFilter{
				Filters: map[string]string{"image-name": "Debian-12-*"},
			},
			MostRecent: true,
		},
	}
	return stepSourceOMIInfo
}

func TestMostRecentOmiFilter(t *testing.T) {
	step := mostRecentOmiFilterStep()

	state, err := getState()
	if state == nil {
		t.Fatalf("error retrieving state %s", err.Error())
	}

	action := step.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("should continue, but: %v", action)
	}

}
