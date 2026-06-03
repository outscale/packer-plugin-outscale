packer {
  required_plugins {
    outscale = {
      version = ">=v0.1.0"
      source  = "github.com/outscale/outscale"
    }
  }
}

variable "osc_access_key" {
  type    = string
  default = "${env("OSC_ACCESS_KEY")}"
}

variable "osc_secret_key" {
  type    = string
  default = "${env("OSC_SECRET_KEY")}"
}

variable "omi_name" {
  type = string
}

variable "ssh_username" {
  type = string
}

variable "region" {
  type = string
}
source "outscale-bsu" "bsu" {
  region  = var.region
  vm_type = "tinav7.c1r1p1"
  source_omi_filter {
    filters = {
      image-name = "Ubuntu-24.04-*"
    }
    owners      = ["Outscale"]
    most_recent = true
  }
  ssh_username                = var.ssh_username
  omi_name                    = var.omi_name
  associate_public_ip_address = true
  force_deregister            = true
}

source "outscale-bsusurrogate" "rhel" {
  region = var.region
  launch_block_device_mappings {
    delete_on_vm_deletion = true
    device_name           = "/dev/xvdf"
    iops                  = 3000
    volume_size           = "20"
    volume_type           = "io1"
  }
  omi_name = var.omi_name
  omi_root_device {
    delete_on_vm_deletion = true
    device_name           = "/dev/sda1"
    source_device_name    = "/dev/xvdf"
    volume_size           = "20"
    volume_type           = "standard"
  }
  source_omi_filter {
    filters = {
      image-name = "RHEL-10-*"
    }
    owners      = ["Outscale"]
    most_recent = true
  }
  ssh_interface = "public_ip"
  ssh_username  = var.ssh_username
  vm_type       = "tinav7.c2r4p1"
}

build {
  sources = ["source.outscale-bsusurrogate.rhel"]
}
