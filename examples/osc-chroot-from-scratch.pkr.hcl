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

variable "region" {
  type = string
}

variable "ssh_username" {
  type = string
}

source "outscale-chroot" "chroot-scratch" {
  access_key   = var.osc_access_key
  secret_key   = var.osc_secret_key
  region       = var.region
  omi_name     = var.omi_name
  from_scratch = true
  pre_mount_commands = [
    "parted {{.Device}} mklabel msdos mkpart primary 1M 100% set 1 boot on print",
    "partprobe",
    "mkfs.ext4 {{.Device}}1"
  ]
  root_volume_size = 50
  root_device_name = "/dev/sda1"
  omi_block_device_mappings {
    device_name = "/dev/sda1"
    volume_type = "gp2"
  }
}
build {
  sources = ["source.outscale-chroot.chroot-scratch"]
}
