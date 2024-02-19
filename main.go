package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"

	"github.com/outscale/packer-plugin-outscale/builder/bsu"
	"github.com/outscale/packer-plugin-outscale/builder/bsusurrogate"
	"github.com/outscale/packer-plugin-outscale/builder/bsuvolume"
	"github.com/outscale/packer-plugin-outscale/builder/chroot"
	"github.com/outscale/packer-plugin-outscale/datasource/omi"
	"github.com/outscale/packer-plugin-outscale/version"
)

func main() {
	pps := plugin.NewSet()
	pps.SetVersion(version.PluginVersion)
	pps.RegisterBuilder("bsu", new(bsu.Builder))
	pps.RegisterBuilder("chroot", new(chroot.Builder))
	pps.RegisterBuilder("bsusurrogate", new(bsusurrogate.Builder))
	pps.RegisterBuilder("bsuvolume", new(bsuvolume.Builder))
	pps.RegisterDatasource("omi", new(omi.Datasource))
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
