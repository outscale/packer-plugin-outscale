packer {
  required_plugins {
    outscale = {
      version = ">=v0.1.0"
      source  = "github.com/outscale/outscale"
    }
  }
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
      image-name = "RockyLinux-10-*"
    }
    owners      = ["Outscale"]
    most_recent = true
  }

  ssh_interface               = "public_ip"
  ssh_username                = var.ssh_username
  omi_name                    = var.omi_name
  associate_public_ip_address = true
  force_deregister            = true
}

build {
  sources = ["source.outscale-bsu.bsu"]
}
