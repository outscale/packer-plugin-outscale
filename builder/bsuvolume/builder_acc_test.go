// TODO: explain how to delete the image.
package bsuvolume

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/outscale/packer-plugin-outscale/internal/testacc"
)

func TestAccBuilder_basic(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsuvolume_basic_test",
		Template: testBuilderAccBasic,
		Check:    testacc.CheckWithLogs,
	}
	acctest.TestPlugin(t, testCase)
}

const testBuilderAccBasic = `
{
    "builders": [
        {
            "type": "outscale-bsuvolume",
            "vm_type": "tinav7.c1r1p1",
            "source_omi_filter": {
			    "filters": {
		       		"image-name": "Debian-12-*"
		        },
		        "owners": ["Outscale"],
		        "most_recent": true
            },
            "ssh_username": "outscale",
            "bsu_volumes": [
                {
                    "volume_type": "gp2",
                    "device_name": "/dev/xvdf",
                    "delete_on_vm_deletion": false,
                    "tags": {
                        "zpool": "data",
                        "Name": "Data1"
                    },
                    "volume_size": 10
                },
                {
                    "volume_type": "gp2",
                    "device_name": "/dev/xvdg",
                    "tags": {
                        "zpool": "data",
                        "Name": "Data2"
                    },
                    "delete_on_vm_deletion": false,
                    "volume_size": 10
                },
                {
                    "volume_size": 10,
                    "tags": {
                        "Name": "Data3",
                        "zpool": "data"
                    },
                    "delete_on_vm_deletion": false,
                    "device_name": "/dev/xvdh",
                    "volume_type": "gp2"
                }
            ]
        }
    ]
}
`
