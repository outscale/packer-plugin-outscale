// TODO: explain how to delete the image.
package bsu

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestAccBuilder_basic(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsu_basic_test",
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

func TestAccBuilder_GoodProductCode(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsu_product_code_test",
		Template: testBuilderAccWithGoodProductCode,
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
	"builders": [{
		"type": "outscale-bsu",
		"region": "eu-west-2",
		"vm_type": "tinav5.c1r1p1",
		"source_omi_filter": {
	        "filters": {
	          "image-name": "Debian-12-*"
	        },
	        "owners": ["Outscale"],
	        "most_recent": true
        },
		"ssh_username": "outscale",
		"omi_name": "packer-test",
		"associate_public_ip_address": true,
		"force_deregister": true
	}]
}
`

const testBuilderAccWithGoodProductCode = `
{
	"builders": [{
		"type": "outscale-bsu",
		"region": "eu-west-2",
		"vm_type": "tinav5.c1r1p1",
		"source_omi_filter": {
	        "filters": {
	          "image-name": "Debian-12-*"
	        },
	        "owners": ["Outscale"],
	        "most_recent": true
        },
		"ssh_username": "outscale",
		"omi_name": "packer-test",
		"product_codes": ["0001"],
		"associate_public_ip_address": true,
		"force_deregister": true
	}]
}
`
