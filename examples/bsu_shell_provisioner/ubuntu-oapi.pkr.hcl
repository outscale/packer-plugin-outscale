packer {
  required_plugins {
    outscale = {
      version = ">=1.0.7" # Ensures the Outscale plugin is a ([latest version](latest version https://developer.hashicorp.com/packer/integrations/outscale/outscale))
      source  = "github.com/outscale/outscale"
    }
  }
}

# Define Variables for the Packer Build
variable "region" {
  description = "Outscale region to deploy the instance"
  type        = string
}

variable "new_omi_name" {
  description = "Name of the new Outscale Machine Image (OMI) to be created"
  type        = string
}

variable "vm_type" {
  description = "The Outscale VM type to use while building the OMI"
  type        = string
}

variable "osc_source_image_id" {
  description = "ID of the existing Outscale Machine Image (OMI) used as a base"
  type        = string
}

variable "ssh_username" {
  description = "SSH username used to connect to the instance"
  type        = string
}

# Define the Outscale Builder (Creating a VM and New OMI)
source "outscale-bsu" "create-omi" {
  launch_block_device_mappings {
    delete_on_vm_deletion = true # Ensures the disk is deleted when the instance is terminated
    device_name           = "/dev/sda1"
    volume_size           = 100   # Disk size in GB
    volume_type           = "gp2" # General Purpose SSD
  }

  omi_boot_modes = ["legacy", "uefi"] # Enable boot modes on the new Outscale Machine Image (OMI) to be created
  boot_mode      = "legacy"           # Enable boot mode on the VM to use while building the OMI.(It must be active on the `source_omi`).
  communicator   = "ssh"              # Specifies SSH as the communication method
  region         = var.region
  vm_type        = var.vm_type
  source_omi     = var.osc_source_image_id
  omi_name       = "${var.new_omi_name}-${formatdate("YYYYMMDD-HHmmss", timestamp())}" # Append timestamp to OMI name
  ssh_username   = var.ssh_username

  ssh_interface               = "public_ip" # Use the public IP for SSH connection
  associate_public_ip_address = true        # Ensures the instance gets a public IP
  force_deregister            = false       # If set to true, Packer will first deregister any existing OMI with the same name before creating a new one.
}

# Define the Build Steps
build {
  sources = ["source.outscale-bsu.create-omi"]

  # Install oapi-cli using a shell script
  provisioner "shell" {
    script       = "./scripts/install_oapi-cli.sh"
    pause_before = "3s"
  }

  # Install nginx
  provisioner "shell" {
    inline = [
      "sudo apt-get -y update",
      "sudo apt-get -y install nginx",
      "sudo systemctl enable nginx",
      "nginx -v"
    ]
  }
}
