
data "outscale-omi" "test" {
	filters = {
		name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
		root-device-type    = "ebs"
		virtualization-type = "hvm"
	}
	owners = ["099720109477"]
	most_recent = true
}
  
  source "outscale-bsusurrogate" "basic-data-source-example" {
	launch_block_device_mappings {
        delete_on_vm_deletion = true
        device_name = "/dev/xvdf" 
        iops = 300
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
	communicator  = "ssh"
	ssh_username  = "ubuntu"
  }
  
  build {
	sources = [
	  "source.outscale-bsusurrogate.basic-data-source-example"
	]
  }
