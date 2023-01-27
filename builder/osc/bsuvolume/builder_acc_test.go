// TODO: explain how to delete the image.
package bsuvolume

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestAccBuilder_basic(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsuvolume_basic_test",
		Template: testBuilderAccBasic,
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

const testBuilderAccBasic = `
{
    "builders": [
        {
            "type": "outscale-bsuvolume",
            "region": "eu-west-2",
            "vm_type": "t2.micro",
            "source_omi": "ami-6b9aeecd",
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
