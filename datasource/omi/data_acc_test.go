package omi

import (
	_ "embed"
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

func TestAccDatasource_OutscaleOmi(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name:     "bsu_data_source_test",
		Template: testBuilderAccBasicDataSource,
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

const testBuilderAccBasicDataSource = `
data "outscale-omi" "test" {
	filters = {
		name                = "RockyLinux-8-*"
		root-device-type    = "ebs"
		virtualization-type = "hvm"
	}
	owners = ["Outscale"]
	most_recent = true
}
  
  source "outscale-bsusurrogate" "basic-data-source-example" {
	launch_block_device_mappings {
        delete_on_vm_deletion = true
        device_name = "/dev/xvdf" 
        iops = 3000
        volume_size = "10"
        volume_type = "io1"
    }
    omi_name = "packer-test{{timestamp}}" 
    omi_root_device {
        delete_on_vm_deletion = true
        device_name = "/dev/sda1"
        source_device_name = "/dev/xvdf"
        volume_size = "10"
        volume_type = "standard"
    }
	region        = "eu-west-2"
	source_omi    = data.outscale-omi.test.id
	vm_type = "t2.micro"
	ssh_username  = "outscale"
  }
  
  build {
	sources = [
	  "source.outscale-bsusurrogate.basic-data-source-example"
	]
  }
`
