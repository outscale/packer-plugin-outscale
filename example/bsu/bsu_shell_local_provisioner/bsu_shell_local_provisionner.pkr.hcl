packer {
  required_plugins {
    outscale = {
      version = ">= 1.4.0"
      source  = "github.com/outscale/outscale"
    }
  }
}

variable "access_key" {
    type    = string
    description = "your access_key"
    sensitive   = true
}
variable "secret_key" {
    type    = string
    description = "your secret_key"
    sensitive   = true
}
variable "region" {
    description = "your region"
    type    = string
}

variable "new_omi_name" {
    type    = string
    description = "name of the OMI to be create"
}

variable "vm_type" {
    type    = string
    description = "The Outscale VM type to use while building the OMI"
}

variable "osc_source_image_id" {
    type    = string
    description = "the existing image (OMI) to be use to create the new one"
}

variable "ssh_username" {
    description = "the existing image (OMI) to be use to create the new one"
    type    = string
}

source "outscale-bsu" "create-omi" {

    launch_block_device_mappings {
	delete_on_vm_deletion = true
	device_name           = "/dev/sda1"
	volume_size           = 100
	volume_type           = "gp2"
    }
    communicator     = "ssh"
    region           = "${var.region}"
    vm_type          = "${var.vm_type}"
    source_omi       = "${var.osc_source_image_id}"
    omi_name         = "${var.new_omi_name}"
    ssh_username     =  "${var.ssh_username}"
    root_device_name = "/dev/sda1" 
    ssh_interface               = "public_ip"
    associate_public_ip_address = true
   
    # list of accounts that can access the new OMI :
    #omi_account_ids = [
    #	"xxxxxxxxxx",
    #]
}

build {
    sources = ["source.outscale-bsu.create-omi"]
    
    # Copy file to home directory in the new omi
    provisioner "file" {
	source = "../README.md"
	destination = "/tmp/README.md"
    }
    
    # Install oapi-cli using bash file script
    provisioner "shell-local" {
	script       = "install_oapi-cli.sh"
	pause_before = "3s"
    }

    # Install nginx using bash command line 
    provisioner "shell-local" {
	inline = [
	    "sudo apt-get -y update",
	    "sudo apt-get -y install nginx",
	    "sudo systemctl enable nginx",
	    "nginx -v"
    ]
    }
}
