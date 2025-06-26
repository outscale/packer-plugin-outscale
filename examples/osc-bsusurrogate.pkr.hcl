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

variable "volsize" {
  type    = string
  default = "10"
}

variable "region" {
  type    = string
  default = "${env("OUTSCALE_REGION")}"
}
variable "username" {
  type    = string
  default = "outscale"
}

source "outscale-bsusurrogate" "centos8" {
  launch_block_device_mappings {
    delete_on_vm_deletion = true
    device_name           = "/dev/xvdf"
    iops                  = 3000
    volume_size           = "${var.volsize}"
    volume_type           = "io1"
  }
  omi_name = "packer-test"
  omi_root_device {
    delete_on_vm_deletion = true
    device_name           = "/dev/sda1"
    source_device_name    = "/dev/xvdf"
    volume_size           = "${var.volsize}"
    volume_type           = "standard"
  }
  source_omi    = "${var.omi_name}"
  ssh_interface = "public_ip"
  ssh_username  = "${var.username}"
  vm_type       = "tinav4.c2r4p1"
}

build {
  sources = ["source.outscale-bsusurrogate.centos8"]
}