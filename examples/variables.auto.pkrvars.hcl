# Name of the new Outscale Machine Image (OMI) to be created.
# This is the custom image that will be generated after the Packer build.
# The builder will automatically append the current datetime to the OMI name.
# You can modify this behavior directly in the builder configuration if needed.
omi_name = "packer-example-{{timestamp}}"

# SSH username used to connect to the instance
# Ensure this user has proper permissions to install and configure software.
ssh_username = "outscale"

# The Outscale region where the VM will be deployed
# Example: "eu-west-2" is a specific Outscale region in Europe.
region = "us-east-2"
