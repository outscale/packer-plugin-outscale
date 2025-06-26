packer {
  required_plugins {
    outscale = {
      version = ">=v0.1.0"
      source  = "github.com/outscale/outscale"
    }
  }
}

variable "omi_name" {
  type    = string
  default = "${env("OMI_NAME")}"
}

variable "username" {
  type    = string
  default = "outscale"
}
variable "omi" {
  type    = string
  default = "${env("SOURCE_OMI")}"
}
variable "volsize" {
  type    = string
  default = "10"
}
source "outscale-bsusurrogate" "test" {
  launch_block_device_mappings {
    delete_on_vm_deletion = true
    device_name           = "/dev/xvdf"
    iops                  = 3000
    volume_size           = "${var.volsize}"
    volume_type           = "io1"
  }
  omi_name = "${var.omi_name}"
  source_omi_filter {
    filters = {
      name                = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    owners      = ["099720109477"]
    most_recent = true
  }
  omi_root_device {
    delete_on_vm_deletion = true
    device_name           = "/dev/sda1"
    source_device_name    = "/dev/xvdf"
    volume_size           = "${var.volsize}"
    volume_type           = "standard"
  }
  source_omi    = "${var.omi}"
  ssh_interface = "public_ip"
  ssh_username  = "${var.username}"
  vm_type       = "tinav4.c2r4p1"
}

build {
  sources = ["source.outscale-bsusurrogate.test"]
}
