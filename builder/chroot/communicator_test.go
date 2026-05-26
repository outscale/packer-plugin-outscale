package chroot_test

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/outscale/packer-plugin-outscale/builder/chroot"
)

func TestCommunicator_ImplementsCommunicator(t *testing.T) {
	var raw any = &chroot.Communicator{}
	if _, ok := raw.(packersdk.Communicator); !ok {
		t.Fatalf("Communicator should be a communicator")
	}
}
