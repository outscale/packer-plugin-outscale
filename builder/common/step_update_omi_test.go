package common

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"

	"bytes"
	"context"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// Create statebag for running test
func tState() multistep.StateBag {
	state := new(multistep.BasicStateBag)
	accessConfig := &AccessConfig{}

	state.Put("ui", &packersdk.BasicUi{
		Reader: new(bytes.Buffer),
		Writer: new(bytes.Buffer),
	})
	state.Put("omis", map[string]string{"us-west-2": "omi-12345"})
	state.Put("snapshots", map[string][]string{"us-west-2": {"snap-0012345"}})
	state.Put("accessConfig", accessConfig)
	return state
}

func TestUpdateOmi(t *testing.T) {

	stepUpdateOMIAttributes := StepUpdateOMIAttributes{
		AccountIds:         []string{},
		SnapshotAccountIds: []string{},
		RawRegion:          "us-west-2",
		GlobalPermission:   true,
	}
	state := tState()

	action := stepUpdateOMIAttributes.Run(context.Background(), state)
	if err := state.Get("error"); err != nil {
		t.Fatalf("should not error, but: %v", err)
	}

	if action != multistep.ActionContinue {
		t.Fatalf("shoul continue, but: %v", action)
	}

}
