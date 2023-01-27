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

const testBuilderAccBasic = `
{
	"builders": [{
		"type": "outscale-bsu",
		"region": "eu-west-2",
		"vm_type": "t2.micro",
		"source_omi": "ami-68ed4301",
		"ssh_username": "outscale",
		"omi_name": "packer-test",
		"associate_public_ip_address": true,
		"force_deregister": true
	}]
}
`
