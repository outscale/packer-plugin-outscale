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
  type    = string
  default = env("OSC_REGION")
}

variable "ssh_username" {
  type = string
}

source "outscale-chroot" "chroot" {
  access_key   = var.osc_access_key
  secret_key   = var.osc_secret_key
  region       = var.region
  omi_name     = var.omi_name
  from_scratch = false
  source_omi_filter {
    filters = {
      image-name = "Debian-12-*"
    }
    owners      = ["Outscale"]
    most_recent = true
  }
}
build {
  sources = ["source.outscale-chroot.chroot"]
}
