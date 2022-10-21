# The Example Folder
This folder must contain a fully working example of the plugin usage. The example must define the required_plugins block.
The folder can contain multiple HCL2 compatible files. This exapmles should be tested with packer init -upgrade . and packer build ..

If the plugin requires authentication, the configuration should be set as environment variables. 

### env:
	export OSC_ACCESS_KEY=<ACCESS_KEY>
	export OSC_SECRET_KEY=<SECRET_KEY>
	export OSC_REGION=eu-west-2 # Outscale Region
	export OMI_NAME=<OMI_NAME>