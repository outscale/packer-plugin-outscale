package bsusurrogate_test

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/outscale/packer-plugin-outscale/internal/testacc"
)

func TestAccBuilder_basic(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsusurrogate_basic_test",
		Template: testBuilderAccBasic,
		Check:    testacc.CheckWithLogs,
	}
	acctest.TestPlugin(t, testCase)
}

func TestAccBuilder_VmTerminate(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsusurrogate_vm_terminate_test",
		Template: testBuilderAccVmTerminate,
		Check:    testacc.CheckWithLogs,
	}
	acctest.TestPlugin(t, testCase)
}

const testBuilderAccBasic = `
{
	"builders": [{
		"type": "outscale-bsusurrogate",
		"vm_type": "tinav7.c1r1p1",
		"source_omi_filter": {
		    "filters": {
	       		"image-name": "Debian-12-*"
	        },
	        "owners": ["Outscale"],
	        "most_recent": true
        },
		"ssh_username": "outscale",
		"omi_name": "packer-test {{timestamp}}",
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
		"vm_type": "tinav7.c1r1p1",
		"source_omi_filter": {
		    "filters": {
	       		"image-name": "Debian-12-*"
	        },
	        "owners": ["Outscale"],
	        "most_recent": true
        },
		"ssh_username": "outscale",
		"omi_name": "packer-test {{timestamp}}",
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
