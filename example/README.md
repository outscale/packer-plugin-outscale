The Example Folder
This folder must contain a fully working example of the plugin usage. The example must define the required_plugins block.
The folder can contain multiple HCL2 compatible files. This exapmles should be tested with packer init -upgrade . and packer build ..

If the plugin requires authentication, the configuration should be set as environment variables. Example:
      env:
	export OUTSCALE_ACCESSKEYID=<ACCESS_KEY>
	export OUTSCALE_SECRETKEYID=<SECRET_KEY>
	export OUTSCALE_REGION=eu-west-2 # Outscale Region
	export OMI_NAME=<OMI_NAME>

