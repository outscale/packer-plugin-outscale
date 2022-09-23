package common

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/outscale/osc-sdk-go/osc"

	"bytes"
	"context"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// Create statebag for running test
func getState() multistep.StateBag {
	state := new(multistep.BasicStateBag)
	accessConfig := &AccessConfig{}
	accessConfig.RawRegion = "eu-west-2"
	var oscConn *osc.APIClient
	var err error
	if oscConn, err = accessConfig.NewOSCClient(); err != nil {
		return nil
	}
	state.Put("osc", oscConn)
	state.Put("ui", &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	})
	state.Put("accessConfig", accessConfig)
	return state
}

func TestMostRecentOmiFilter(t *testing.T) {
	stepSourceOMIInfo := StepSourceOMIInfo{
		SourceOmi: "ami-7cab7c18",
		OmiFilters: OmiFilterOptions{
			MostRecent: true,
		},
	}
	state := getState()
	if state == nil {
		t.Fatalf("error retrieving state, but")
	}

	action := stepSourceOMIInfo.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", action)
	}

}
