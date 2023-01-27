package bsusurrogate

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestAccBuilder_basic(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsusurrogate_basic_test",
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
func TestAccBuilder_VmTerminate(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsusurrogate_vm_terminate_test",
		Template: testBuilderAccVmTerminate,
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
		"type": "outscale-bsusurrogate",
		"region": "eu-west-2",
		"vm_type": "t2.micro",
		"source_omi": "ami-6b9aeecd",
		"ssh_username": "outscale",
		"omi_name": "packer-test {{timestamp}}",
		"subregion_name": "eu-west-2a",
		"launch_block_device_mappings" : [
			{
			"volume_type" : "io1",
			"device_name" : "/dev/xvdf",
			"delete_on_vm_deletion" : false,
			"volume_size" : 10,
			"iops": 300
			}
		],
		"omi_root_device":{
			"source_device_name": "/dev/xvdf",
			"device_name": "/dev/sda1",
			"delete_on_vm_deletion": true,
			"volume_size": 10,
			"volume_type": "standard"
		}

	}]
}
`

const testBuilderAccVmTerminate = `
{
	"builders": [{
		"type": "outscale-bsusurrogate",
		"region": "eu-west-2",
		"vm_type": "t2.micro",
		"source_omi": "ami-6b9aeecd",
		"ssh_username": "outscale",
		"omi_name": "packer-test {{timestamp}}",
		"subregion_name": "eu-west-2a",
		"launch_block_device_mappings" : [
			{
			"volume_type" : "io1",
			"device_name" : "/dev/xvdf",
			"delete_on_vm_deletion" : false,
			"volume_size" : 10,
			"iops": 300
			}
		],
		"omi_root_device":{
			"source_device_name": "/dev/xvdf",
			"device_name": "/dev/sda1",
			"delete_on_vm_deletion": true,
			"volume_size": 10,
			"volume_type": "standard"
		},
		"shutdown_behavior": "terminate"
	}]
}
`
