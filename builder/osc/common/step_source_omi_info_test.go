package common

import (
	"bytes"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"

	"context"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// Create statebag for running test
func getState() (multistep.StateBag, error) {
	state := new(multistep.BasicStateBag)
	accessConfig := &AccessConfig{}
	accessConfig.RawRegion = "eu-west-2"
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

func TestMostRecentOmiFilter(t *testing.T) {
	stepSourceOMIInfo := StepSourceOMIInfo{
		SourceOmi: "ami-e58ac287",
		OmiFilters: OmiFilterOptions{
			MostRecent: true,
		},
	}
	state, err := getState()
	if state == nil {
		t.Fatalf("error retrieving state %s", err.Error())
	}

	action := stepSourceOMIInfo.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", action)
	}

}
