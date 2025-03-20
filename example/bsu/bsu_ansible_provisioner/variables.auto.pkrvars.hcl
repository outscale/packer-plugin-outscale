# Name of the new Outscale Machine Image (OMI) to be created.
# This is the custom image that will be generated after the Packer build.
# The builder will automatically append the current datetime to the OMI name.
# You can modify this behavior directly in the builder configuration if needed.
new_omi_name = "Ubuntu-24.04-Apache"

# SSH username used to connect to the instance
# Ensure this user has proper permissions to install and configure software.
ssh_username = "outscale"

# Outscale VM type to be used during the build process
# Example: "tinav6.c4r8p2" -> Tina V6, 4 vCPUs, 8GB RAM
vm_type = "tinav6.c4r8p2"

# The Outscale region where the VM will be deployed
# Example: "eu-west-2" is a specific Outscale region in Europe.
region = "eu-west-2"

# The ID of the base Outscale Machine Image (OMI) used for the build
# This is the source image from which the new image will be created.
# Example: "ami-860c2495" corresponds to Ubuntu 24.04.
osc_source_image_id = "ami-860c2495"
