packer {
  required_plugins {
    windows-update = {
      version = ">=0.14.0"
      source  = "github.com/rgl/windows-update"
    }
  }
}

source "outscale-bsu" "windows" {
  region  = "eu-west-2"
  vm_type = "t2.micro"
  source_omi_filter {
    filters = {
      image-name = "Windows-10-GOLDEN"
    }
    owners = ["Outscale"]
  }
  ssh_username                = "outscale"
  omi_name                    = "packer-test"
  associate_public_ip_address = true
  force_deregister            = true
}

build {
  sources = ["source.outscale-bsu.windows"]
}