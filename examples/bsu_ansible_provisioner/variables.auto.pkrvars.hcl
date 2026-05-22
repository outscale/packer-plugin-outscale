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
region = "us-east-2"

# The name pattern of the base Outscale Machine Image (OMI) used for the build.
# The most recent Outscale-owned image matching this pattern will be selected.
osc_source_image_name = "Ubuntu-24.04-*"

# The ID of the base Outscale Machine Image (OMI) used for the build.
# Uncomment this and the matching `source_omi` line in the template to pin a specific image.
# osc_source_image_id = "ami-b29cea33"
