Type: `outscale-bsu`
Artifact BuilderId: `oapi.outscale.bsu`

The `outscale-bsu` Packer builder is able to create Outscale OMIs backed by BSU
volumes for use in [Flexible Compute Unit](https://docs.outscale.com/en/userguide/Flexible-Compute-Unit-(FCU).html). For more information on
the difference between BSU-backed VMs and VM-store backed
VMs, see the ["storage for the root device" section in the Outscale
documentation](https://docs.outscale.com/en/userguide/Defining-Block-Device-Mappings.html).

This builder builds an OMI by launching an Outscale VM from a source OMI,
provisioning that running machine, and then creating an OMI from that machine.
This is all done in your own Outscale account. The builder will create temporary
keypairs, security group rules, etc. that provide it temporary access to the
VM while the image is being created. This simplifies configuration quite
a bit.

The builder does _not_ manage OMIs. Once it creates an OMI and stores it in
your account, it is up to you to use, delete, etc. the OMI.

-> **Note:** Temporary resources are, by default, all created with the
prefix `packer`. This can be useful if you want to restrict the security groups
and key pairs Packer is able to operate on.

## Configuration Reference

There are many configuration options available for the builder. They are
segmented below into two categories: required and optional parameters. Within
each category, the available configuration keys are alphabetized.

In addition to the options listed here, a
[communicator](/docs/templates/legacy_json_templates/communicator) can be configured for this
builder.

### Required:

- `access_key` (string) - The access key used to communicate with OUTSCALE. [Learn how to set this](/docs/builder/outscale#authentication)

- `omi_name` (string) - The name of the resulting OMIS that will appear when managing OMIs in the Outscale console or via APIs. This must be unique. To help make this unique, use a function like `timestamp` (see [template engine](/docs/templates/legacy_json_templates/engine) for more info).

- `vm_type` (string) - The Outscale VM type to use while building the OMI, such as `tinav6.c4r8p2`.

- `region` (string) - The name of the region, such as `us-east-1`, in which to launch the Outscale VM to create the OMI.

- `secret_key` (string) - The secret key used to communicate with Outscale. [Learn how to set this](/docs/builder/outscale#authentication)

- `source_omi` (string) - The initial OMI used as a base for the newly created machine. `source_omi_filter` may be used instead to populate this automatically.

### Optional:

- `omi_block_device_mappings` (array of block device mappings) - Add one or more [block device mappings](https://docs.outscale.com/en/userguide/Defining-Block-Device-Mappings.html) to the OMI. These will be attached when booting a new VM from your OMI. To add a block device during the Packer build see `launch_block_device_mappings` below. Your options here may vary depending on the type of VM you use. The block device mappings allow for the following configuration:

  - `delete_on_vm_deletion` (boolean) - Indicates whether the BSU volume is deleted on VM termination. Default `false`. **NOTE**: If this value is not explicitly set to `true` and volumes are not cleaned up by an alternative method, additional volumes will accumulate after every build.

  - `device_name` (string) - The device name exposed to the VM (for example, `/dev/sdh` or `xvdh`). Required for every device in the block device mapping.

  - `iops` (number) - The number of I/O operations per second (IOPS) that the volume supports. See the documentation on
    [IOPs](https://docs.outscale.com/en/userguide/About-Volumes.html#_volume_types_and_iops)
    for more information

  - `no_device` (boolean) - Suppresses the specified device included in the
    block device mapping of the OMI

  - `snapshot_id` (string) - The ID of the snapshot

  - `volume_size` (number) - The size of the volume, in GiB. Required if not specifying a `snapshot_id`

  - `volume_type` (string) - The volume type. `gp2` for General Purpose (SSD) volumes, `io1` for Provisioned IOPS (SSD) volumes, and `standard` for Magnetic volumes

- `omi_description` (string) - The description to set for the resulting OMI(s). By default this description is empty. This is a [template engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

- `omi_account_ids` (array of strings) - A list of account IDs that have access to launch the resulting OMI(s). By default no additional users other than the user creating the OMIS has permissions to launch it.

- `omi_boot_modes` (array of strings) - A list of boot modes (`legacy` and/or `uefi`) to enable on the resulting OMI.

- `boot_mode` (strings) - The boot mode (`legacy` or `uefi`) on the VM to use while building the OMI. (It must be active on the `source_omi`).

- `product_codes` ([]string) - A list of product codes to associate with the OMI. By default no product codes are associated with the OMI.

- `global_permission` (boolean) - This option is useful to make the OMI publicly accessible.

- `omi_virtualization_type` (string) - The type of virtualization for the OMI you are building. This option must match the supported virtualization type of `source_omi`. Can be `paravirtual` or `hvm`.

- `associate_public_ip_address` (boolean) - If using a non-default Net, public IP addresses are not provided by default. If this is toggled, your new VM will get a Public IP.

- `subregion_name` (string) - Destination subregion to launch VM in. Leave this empty to allow Outscale to auto-assign.

- `custom_endpoint_oapi` (string) - This option is useful if you use a cloud
  provider whose API is compatible with Outscale OAPI. Specify another endpoint
  like this `outscale.com/oapi/latest`.

- `disable_stop_vm` (boolean) - Packer normally stops the build
  VM after all provisioners have run. For Windows VMs, it is
  sometimes desirable to [run Sysprep](<https://docs.microsoft.com/en-us/previous-versions/windows/it-pro/windows-vista/cc721940(v=ws.10)>) which will stop the VM for you. If this is set to `true`, Packer
  _will not_ stop the VM but will assume that you will send the stop
  signal yourself through your final provisioner. You can do this with a
  [windows-shell provisioner](/docs/provisioner/windows-shell).

  Note that Packer will still wait for the VM to be stopped, and
  failing to send the stop signal yourself, when you have set this flag to
  `true`, will cause a timeout.

- `bsu_optimized` (boolean) - If true, the VM is created with optimized BSU I/O.

- `force_delete_snapshot` (boolean) - Force Packer to delete snapshots
  associated with OMIs, which have been deregistered by `force_deregister`.
  Default `false`.

- `force_deregister` (boolean) - Force Packer to first deregister an existing
  OMIS if one with the same name already exists. Default `false`.

- `insecure_skip_tls_verify` (boolean) - This allows skipping TLS
  verification of the OAPI endpoint. The default is `false`.

- `launch_block_device_mappings` (array of block device mappings) - Add one
  or more block devices before the Packer build starts. If you add VM
  store volumes or BSU volumes in addition to the root device volume, the
  created OMIS will contain block device mapping information for those
  volumes. Outscale creates snapshots of the source VM's root volume and
  any other BSU volumes described here. When you launch an VM from this
  new OMI, the VM automatically launches with these additional volumes,
  and will restore them from snapshots taken from the source VM.

- `run_tags` (object of key/value strings) - Tags to apply to the VM
  that is _launched_ to create the OMI. These tags are _not_ applied to the
  resulting OMIS unless they're duplicated in `tags`. This is a [template
  engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

- `run_volume_tags` (object of key/value strings) - Tags to apply to the
  volumes that are _launched_ to create the OMI. These tags are _not_ applied
  to the resulting OMIS unless they're duplicated in `tags`. This is a
  [template engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

- `security_group_id` (string) - The ID (_not_ the name) of the security
  group to assign to the VM. By default this is not set and Packer will
  automatically create a new temporary security group to allow SSH access.
  Note that if this is specified, you must be sure the security group allows
  access to the `ssh_port` given below.

- `security_group_ids` (array of strings) - A list of security groups as
  described above. Note that if this is specified, you must omit the
  `security_group_id`.

- `shutdown_behavior` (string) - Automatically terminate VMs on
  shutdown in case Packer exits ungracefully. Possible values are "stop" and
  "terminate", default is `stop`.

- `skip_create_omi` (boolean) - Set to true if you want to skip snapshot creation.
  No image will be created if set to true. Default `false`.

- `skip_region_validation` (boolean) - Set to true if you want to skip
  validation of the region configuration option. Default `false`.

- `snapshot_groups` (array of strings) - A list of groups that have access to
  create volumes from the snapshot(s). By default no groups have permission
  to create volumes from the snapshot(s). `all` will make the snapshot
  publicly accessible.

- `snapshot_tags` (object of key/value strings) - Tags to apply to snapshot.
  They will override OMIS tags if already applied to snapshot. This is a
  [template engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

- `source_omi_filter` (object) - Filters used to populate the `source_omi` field.

  - `filters` (map of strings) - filters used to select a `source_omi`.
  - `owners` (array of strings) - Filters the images by their owner. You may specify one or more Outscale account IDs, "self" (which will use the account whose credentials you are using to run Packer). This option is required for security reasons.

    Example:
    #### HCL
    ```hcl
    source_omi_filter {
      filters = {
        image-name = "image-name-in-account"
        root-device-type    = "ebs" # or "bsu"
        virtualization-type = "hvm"
      }
      owners      = ["339215505907"]
    }
    ```

    #### JSON
    ```json
    {
      "source_omi_filter": {
        "filters": {
          "virtualization-type": "hvm",
          "image-name": "image-name-in-account",
          "root-device-type": "ebs"
        },
        "owners": ["099720109477"]
      }
    }
    ```

    This selects an Ubuntu 16.04 HVM BSU OMIS from Canonical. NOTE:
    This will fail unless _exactly_ one OMIS is returned. In the above example,
    `most_recent` will cause this to succeed by selecting the newest image.

- `ssh_private_key_file` (string) - Path to a PEM encoded private key file to use to authenticate with SSH.
  The `~` can be used in path and will be expanded to the home directory
  of current user.


- `ssh_keypair_name` (string) - If specified, this is the key that will be used for SSH with the
  machine. The key must match a key pair name loaded up into the remote.
  By default, this is blank, and Packer will generate a temporary keypair
  unless [`ssh_password`](#ssh_password) is used.
  [`ssh_private_key_file`](#ssh_private_key_file) or
  [`ssh_agent_auth`](#ssh_agent_auth) must be specified when
  [`ssh_keypair_name`](#ssh_keypair_name) is utilized.


- `ssh_agent_auth` (bool) - If true, the local SSH agent will be used to authenticate connections to
  the source instance. No temporary keypair will be created, and the
  values of [`ssh_password`](#ssh_password) and
  [`ssh_private_key_file`](#ssh_private_key_file) will be ignored. The
  environment variable `SSH_AUTH_SOCK` must be set for this option to work
  properly.


- `ssh_interface` (string) - One of `public_ip`, `private_ip`, `public_dns`, or `private_dns`. If set, either the public IP address, private IP address, public DNS name or private DNS name will used as the host for SSH. The default behaviour if inside a Net is to use the public IP address if available, otherwise the private IP address will be used. If not in a Net the public DNS name will be used. Also works for WinRM.

  Where Packer is configured for an outbound proxy but WinRM traffic should be direct, `ssh_interface` must be set to `private_dns` and `<region>.compute.internal` included in the `NO_PROXY` environment variable.

- `subnet_id` (string) - If using Net, the ID of the subnet, such as `subnet-12345def`, where Packer will launch the VM. This field is required if you are using an non-default Net.

- `tags` (object of key/value strings) - Tags applied to the OMIS and relevant snapshots. This is a [template engine](/docs/templates/legacy_json_templates/engine), see [Build template data](#build-template-data) for more information.

- `temporary_key_pair_name` (string) - The name of the temporary key pair to generate. By default, Packer generates a name that looks like `packer_<UUID>`, where &lt;UUID&gt; is a 36 character unique identifier.

- `temporary_security_group_source_cidr` (string) - An IPv4 CIDR block to be authorized access to the VM, when packer is creating a temporary security group. The default is `0.0.0.0/0` (i.e., allow any IPv4 source). This is only used when `security_group_id` or `security_group_ids` is not specified.

- `user_data` (string) - User data to apply when launching the VM. Note that you need to be careful about escaping characters due to the templates being JSON. It is often more convenient to use `user_data_file`, instead. Packer will not automatically wait for a user script to finish before shutting down the VM this must be handled in a provisioner.

- `user_data_file` (string) - Path to a file that will be used for the user data when launching the VM.

- `net_id` (string) - If launching into a Net subnet, Packer needs the Net ID in order to create a temporary security group within the Net. Requires `subnet_id` to be set. If this field is left blank, Packer will try to get the Net ID from the `subnet_id`.

- `net_filter` (object) - Filters used to populate the `net_id` field.
  Example:

  ```json
  {
    "net_filter": {
      "filters": {
        "is-default": "false",
        "ip-range": "/24"
      }
    }
  }
  ```

  This selects the Net with a IPv4 CIDR block of `/24`. NOTE: This will fail unless _exactly_ one Net is returned.

  - `filters` (map of strings) - filters used to select a `vpc_id`. NOTE: This will fail unless _exactly_ one Net is returned.

    `net_id` take precedence over this.

- `windows_password_timeout` (string) - The timeout for waiting for a Windows password for Windows VMs. Defaults to 20 minutes. Example value: `10m`

## Basic Example

Here is a basic example. You will need to provide access keys, and may need to change the OMIS IDs according to what images exist at the time the template is run:
#### HCL
```hcl
source "outscale-bsu" "basic-example" {
   region             = "us-east-1"
   vm_type            = "tinav6.c4r8p2"
   source_omi         = "ami-abcfd0283"
   omi_name           = "packer_osc_{{timestamp}}"
   omi_boot_modes     = ["uefi"]
   ssh_username       = "outscale"
   ssh_interface      = "public_ip"
}

```
#### JSON
```json
{
  "variables": {
    "access_key": "{{env `OUTSCALE_ACCESSKEYID`}}",
    "secret_key": "{{env `OUTSCALE_SECRETKEYID`}}"
  },
  "builders": [
    {
      "type": "outscale-bsu",
      "access_key": "{{user `access_key`}}",
      "secret_key": "{{user `secret_key`}}",
      "region": "us-east-1",
      "source_omi": "ami-abcfd0283",
      "vm_type": "tinav6.c4r8p2",
      "ssh_username": "outscale",
      "omi_name": "packer_osc {{timestamp}}"
    }
  ]
}
```

-> **Note:** Packer can also read the access key and secret access key from
environmental variables. See the configuration reference in the section above
for more information on what environmental variables Packer will look for.

Further information on locating OMIS IDs and their relationship to VM
types and regions can be found in the Outscale Documentation [reference](https://docs.outscale.com/en/userguide/Official-OMIs-Reference.html).

## Accessing the Instance to Debug

If you need to access the VM to debug for some reason, run the builder
with the `-debug` flag. In debug mode, the Outscale builder will save the private key in the current directory and will output the DNS or IP information as well.
You can use this information to access the VM as it is running.

## OMIS Block Device Mappings Example

Here is an example using the optional OMIS block device mappings. Our
configuration of `launch_block_device_mappings` will expand the root volume
(`/dev/sda`) to 40gb during the build (up from the default of 8gb). With
`omi_block_device_mappings` Outscale will attach additional volumes `/dev/sdb` and
`/dev/sdc` when we boot a new VM of our OMI.

#### HCL
##### with `launch_block_device_mappings`
```hcl
// export osc_access_key=$YOURKEY
variable "osc_access_key" {
  type = string
  // default = "hardcoded_key"
}

// export osc_secret_key=$YOURSECRETKEY
variable "osc_secret_key" {
  type = string
  // default = "hardcoded_secret_key"
}

source "outscale-bsu" "basic-example" {
  region             = "us-east-1"
  vm_type            = "tinav6.c4r8p2"
  source_omi         = "ami-abcfd0283"
  omi_name           = "packer_osc_{{timestamp}}"
  ssh_username       = "outscale"
  ssh_interface      = "public_ip"
  omi_boot_modes     = ["legacy","uefi"]
  boot_mode          = "legacy"

  launch_block_device_mappings {
    delete_on_vm_deletion = false
    device_name           = "/dev/sda1"
    volume_size           = 40
    volume_type           = "gp2"
  }
  launch_block_device_mappings {
    device_name = "/dev/sdc"
    volume_size = 50
    volume_type = "gp2"
  }
  launch_block_device_mappings {
    device_name = "/dev/sdc"
    volume_size = 100
    volume_type = "gp2"
  }
}

```
##### with `omi_block_device_mappings`
```hcl
// export osc_access_key=$YOURKEY
variable "aws_access_key" {
  type = string
  // default = "hardcoded_key"
}

// export osc_secret_key=$YOURSECRETKEY
variable "aws_secret_key" {
  type = string
  // default = "hardcoded_secret_key"
}

source "outscale-bsu" "basic-example" {
  region             = "us-east-1"
  vm_type            = "tinav6.c4r8p2"
  source_omi         = "ami-abcfd0283"
  omi_name           = "packer_osc_{{timestamp}}"
  root_device_name   = "/dev/sda1"
  ssh_username       = "outscale"
  ssh_interface      = "public_ip"

  omi_block_device_mappings {
    delete_on_vm_deletion = false
    device_name           = "/dev/sda1"
    snapshot_id           = "snap-792fce69"
    volume_size           = 40
    volume_type           = "gp2"
  }
  omi_block_device_mappings {
    device_name = "/dev/sdc"
    snapshot_id = "snap-792fce69"
    volume_size = 50
    volume_type = "gp2"
  }
}

```

```json
{
  "type": "outscale-bsu",
  "access_key": "YOUR KEY HERE",
  "secret_key": "YOUR SECRET KEY HERE",
  "region": "us-east-1",
  "source_omi": "ami-fce3c696",
  "vm_type": "tinav6.c4r8p2",
  "ssh_username": "ubuntu",
  "omi_name": "packer-quick-start {{timestamp}}",
  "launch_block_device_mappings": [
    {
      "device_name": "/dev/sda1",
      "volume_size": 40,
      "volume_type": "gp2",
      "delete_on_vm_deletion": true
    }
  ],
  "launch_block_device_mappings": [
    {
      "device_name": "/dev/sdb",
      "volume_size": 50,
      "volume_type": "gp2"
    },
    {
      "device_name": "/dev/sdc",
      "volume_size": 100,
      "volume_type": "gp2"
    }
  ]
}
```

## Build template data

In configuration directives marked as a template engine above, the following variables are available:

- `BuildRegion` - The region (for example `eu-west-2`) where Packer is building the OMI.
- `SourceOMI` - The source OMIS ID (for example `ami-a2412fcd`) used to build the OMI.
- `SourceOMIName` - The source OMIS Name (for example `ubutu-390`) used to build the OMI.
- `SourceOMITags` - The source OMIS Tags, as a `map[string]string` object.

## Tag Example

Here is an example using the optional OMIS tags. This will add the tags `OS_Version` and `Release` to the finished OMI. As before, you will need to provide your access keys, and may need to change the source OMIS ID based on what images exist when this template is run:

```json
{
  "type": "outscale-bsu",
  "access_key": "YOUR KEY HERE",
  "secret_key": "YOUR SECRET KEY HERE",
  "region": "us-east-1",
  "source_omi": "ami-fce3c696",
  "vm_type": "tinav6.c4r8p2",
  "ssh_username": "ubuntu",
  "omi_name": "packer-quick-start {{timestamp}}",
  "tags": {
    "OS_Version": "Ubuntu",
    "Release": "Latest",
    "Base_OMI_Name": "{{ .SourceOMIName }}",
    "Extra": "{{ .SourceOMITags.TagName }}"
  }
}
```

-> **Note:** Packer uses pre-built OMIs as the source for building images.
These source OMIs may include volumes that are not flagged to be destroyed on
termination of the VM building the new image. Packer will attempt to
clean up all residual volumes that are not designated by the user to remain
after termination. If you need to preserve those source volumes, you can
overwrite the termination setting by specifying `delete_on_vm_deletion=false`
in the `launch_block_device_mappings` block for the device.
