# Outscale packer examples

This folder contains a number of examples to show how to use Outscale packer with bsu builder.

# How to test examples

First, make sure you have [installed packer](https://developer.hashicorp.com/packer/install).

Each folder is self-contained example.
You will need to setup your credentials through environement variables:
```bash
export PKR_VAR_access_key="myaccesskey"
export PKR_VAR_secret_key="mysecretkey"
export PKR_VAR_region="eu-west-2"
```

If you want to write your credentials or some parameters in variables, just edit `variables.auto.pkrvars.hcl` file.

Once your credentials are configured, you can go to any example folder and test them:
```hcl
packer init .
# Check configuration before build
packer validate .
# Create your OMI
packer build .
```
