package chroot_test

import (
	"testing"

	"github.com/outscale/packer-plugin-outscale/builder/chroot"
)

func TestCopyFilesCleanupFunc_ImplementsCleanupFunc(t *testing.T) {
	var raw any = new(chroot.StepCopyFiles)
	if _, ok := raw.(chroot.Cleanup); !ok {
		t.Fatalf("cleanup func should be a CleanupFunc")
	}
}
