package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	"github.com/outscale/packer-plugin-outscale/builder/osc/bsu"
	"github.com/outscale/packer-plugin-outscale/builder/osc/bsusurrogate"
	"github.com/outscale/packer-plugin-outscale/builder/osc/bsuvolume"
	"github.com/outscale/packer-plugin-outscale/builder/osc/chroot"
	"github.com/outscale/packer-plugin-outscale/version"
)

func main() {
	pps := plugin.NewSet()
	pps.SetVersion(version.PluginVersion)
	pps.RegisterBuilder("bsu", new(bsu.Builder))
	pps.RegisterBuilder("chroot", new(chroot.Builder))
	pps.RegisterBuilder("bsusurrogate", new(bsusurrogate.Builder))
	pps.RegisterBuilder("bsuvolume", new(bsuvolume.Builder))
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
