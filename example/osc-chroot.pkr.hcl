variable "osc_access_key" {
    type    = string
    default = "${env("OSC_ACCESS_KEY")}"
}

variable "osc_secret_key" {
    type    = string
    default = "${env("OSC_SECRET_KEY")}"
}
packer {
  required_plugins {
    outscale = {
      version = ">=v0.1.0"
      source  = "github.com/outscale/outscale"
    }
  }
}

source "outscale-chroot" "windows" {
  access_key = "${var.osc_access_key}"
  secret_key = "${var.osc_secret_key}"
  omi_name = "packer-outscale-chroot {{timestamp}}"
  from_scratch = true
  pre_mount_commands = [
    "parted {{.Device}} mklabel msdos mkpart primary 1M 100% set 1 boot on print",
    "mkfs.ext4 {{.Device}}1"
  ]
  root_volume_size = 15
  root_device_name = "xvdf"
  omi_block_device_mappings {
      device_name = "xvdf"
      volume_type = "gp2"
  }
}
build {
    sources = [ "source.outscale-chroot.windows" ]
}