Type: `outscale-chroot`
Artifact BuilderId: `oapi.outscale.chroot`

The `outscale-chroot` Packer builder is able to create Outscale Machine Images (OMIs) backed by an
BSU volume as the root device. For more information on the difference between
instance storage and BSU-backed instances, see the ["storage for the root
device" section in the Outscale
documentation](https://docs.outscale.com/en/userguide/Home.html).

The difference between this builder and the `outscale-bsu` builder is that this
builder is able to build an BSU-backed OMI without launching a new Outscale
VM. This can dramatically speed up OMI builds for organizations who need
the extra fast build.

~> **This is an advanced builder** If you're just getting started with
Packer, we recommend starting with the [outscale-bsu
builder](/docs/builder/outscale-bsu), which is much easier to use.

The builder does _not_ manage OMIs. Once it creates an OMI and stores it in
your account, it is up to you to use, delete, etc., the OMI.

## How Does it Work?

This builder works by creating a new BSU volume from an existing source OMI and
attaching it into an already-running Outscale VM. Once attached, a
[chroot](https://en.wikipedia.org/wiki/Chroot) is used to provision the system
within that volume. After provisioning, the volume is detached, snapshotted,
and an OMI is made.

Using this process, minutes can be shaved off the OMI creation process because
a new Outscale VM doesn't need to be launched.

There are some restrictions, however. The host Outscale instance where the volume is
attached to must be a similar system (generally the same OS version, kernel
versions, etc.) as the OMI being built. Additionally, this process is much more
expensive because the Outscale VM must be kept running persistently in order
to build OMIs, whereas the other OMI builders start VMs on-demand to
build OMIs as needed.

## Configuration Reference

There are many configuration options available for the builder. They are
segmented below into two categories: required and optional parameters. Within
each category, the available configuration keys are alphabetized.

### Required:

- `access_key` (string) - The access key used to communicate with OUTSCALE. [Learn how to set this](/docs/builder/outscale#authentication)

- `omi_name` (string) - The name of the resulting OMIS that will appear when managing OMIs in the Outscale console or via APIs. This must be unique. To help make this unique, use a function like `timestamp` (see [template engine](/docs/templates/legacy_json_templates/engine) for more info).

- `secret_key` (string) - The secret key used to communicate with Outscale. [Learn how to set this](/docs/builder/outscale#authentication)

- `source_omi` (string) - The initial OMI used as a base for the newly created machine. `source_omi_filter` may be used instead to populate this automatically.

### Optional:

- `omi_description` (string) - The description to set for the resulting OMI(s).
  By default this description is empty. This is a [template engine](/docs/templates/legacy_json_templates/engine),
  see [Build template data](#build-template-data) for more information.

- `omi_account_ids` (array of strings) - A list of account IDs that have access to launch the resulting OMI(s). By default no additional users other than the user creating the OMIS has permissions to launch it.

- `omi_boot_modes` (array of strings) - A list of boot modes (`legacy` and/or `uefi`) to enable on the resulting OMI.

- `boot_mode` (strings) - The boot mode (`legacy` or `uefi`) on the VM to use while building the OMI. (It must be active on the `source_omi`).

- `product_codes` ([]string) - A list of product codes to associate with the OMI. By default no product codes are associated with the OMI.

- `global_permission` (boolean) - This option is useful to make the OMI publicly accessible.

- `omi_virtualization_type` (string) - The type of virtualization for the OMI you are building. This option must match the supported virtualization type of `source_omi`. Can be `paravirtual` or `hvm`.

- `chroot_mounts` (array of array of strings) - This is a list of devices to
  mount into the chroot environment. This configuration parameter requires
  some additional documentation which is in the [Chroot
  Mounts](#chroot-mounts) section. Please read that section for more
  information on how to use this.

- `command_wrapper` (string) - How to run shell commands. This defaults to
  `{{.Command}}`. This may be useful to set if you want to set environmental
  variables or perhaps run it with `sudo` or so on. This is a configuration
  template where the `.Command` variable is replaced with the command to be
  run. Defaults to `{{.Command}}`.

- `copy_files` (array of strings) - Paths to files on the running Outscale
  VM that will be copied into the chroot environment prior to
  provisioning. Defaults to `/etc/resolv.conf` so that DNS lookups work. Pass
  an empty list to skip copying `/etc/resolv.conf`. You may need to do this
  if you're building an image that uses systemd.

- `custom_endpoint_oapi` (string) - This option is useful if you use a cloud
  provider whose API is compatible with Outscale OAPI. Specify another endpoint
  like this `outscale.com/oapi/latest`.

- `device_path` (string) - The path to the device where the root volume of
  the source OMI will be attached. This defaults to "" (empty string), which
  forces Packer to find an open device automatically.

- `force_deregister` (boolean) - Force Packer to first deregister an existing
  OMIS if one with the same name already exists. Default `false`.

- `force_delete_snapshot` (boolean) - Force Packer to delete snapshots
  associated with OMIs, which have been deregistered by `force_deregister`.
  Default `false`.

- `insecure_skip_tls_verify` (boolean) - This allows skipping TLS
  verification of the OAPI endpoint. The default is `false`.

- `from_scratch` (boolean) - Build a new volume instead of starting from an
  existing OMI root volume snapshot. Default `false`. If `true`, `source_omi`
  is no longer used and the following options become required:
  `omi_virtualization_type`, `pre_mount_commands` and `root_volume_size`. The
  below options are also required in this mode only:

- `omi_block_device_mappings` (array of block device mappings) - Add one or more [block device mappings](https://docs.outscale.com/en/userguide/Defining-Block-Device-Mappings.html) to the OMI. These will be attached when booting a new VM from your OMI. To add a block device during the Packer build see `launch_block_device_mappings` below. Your options here may vary depending on the type of VM you use. The block device mappings allow for the following configuration:

  - `delete_on_vm_deletion` (boolean) - Indicates whether the BSU volume is deleted on VM termination. Default `false`. **NOTE**: If this value is not explicitly set to `true` and volumes are not cleaned up by an alternative method, additional volumes will accumulate after every build.

  - `device_name` (string) - The device name exposed to the VM (for example, `/dev/sdh` or `xvdh`). Required for every device in the block device mapping.

  - `iops` (number) - The number of I/O operations per second (IOPS) that the volume supports. See the documentation on
    [IOPs](https://docs.outscale.com/en/userguide/About-Volumes.html#_volume_types_and_iops)
    for more information

  - `no_device` (boolean) - Suppresses the specified device included in the
    block device mapping of the OMI

  - `snapshot_id` (string) - The ID of the snapshot

  - `virtual_name` (string) - The virtual device name. See the documentation on [Block Device Mapping](https://docs.outscale.com/en/userguide/Defining-Block-Device-Mappings.html) for more information

  - `volume_size` (number) - The size of the volume, in GiB. Required if not specifying a `snapshot_id`

  - `volume_type` (string) - The volume type. `gp2` for General Purpose (SSD) volumes, `io1` for Provisioned IOPS (SSD) volumes, and `standard` for Magnetic volumes

- `root_device_name` (string) - The root device name. For example, `xvda`.

- `mount_path` (string) - The path where the volume will be mounted. This is
  where the chroot environment will be. This defaults to
  `/mnt/packer-amazon-chroot-volumes/{{.Device}}`. This is a configuration
  template where the `.Device` variable is replaced with the name of the
  device where the volume is attached.

- `mount_partition` (string) - The partition number containing the /
  partition. By default this is the first partition of the volume, (for
  example, `xvdf1`) but you can designate the entire block device by setting
  `"mount_partition": "0"` in your config, which will mount `xvdf` instead.

- `mount_options` (array of strings) - Options to supply the `mount` command
  when mounting devices. Each option will be prefixed with `-o` and supplied
  to the `mount` command ran by Packer. Because this command is ran in a
  shell, user discretion is advised. See [this manual page for the mount
  command](https://linux.die.net/man/8/mount) for valid file
  system specific options.

- `nvme_device_path` (string) - When we call the mount command (by default
  `mount -o device dir`), the string provided in `nvme_mount_path` will
  replace `device` in that command. When this option is not set, `device` in
  that command will be something like `/dev/sdf1`, mirroring the attached
  device name. This assumption works for most instances but will fail with c5
  and m5 instances. In order to use the chroot builder with c5 and m5
  instances, you must manually set `nvme_device_path` and `device_path`.

- `pre_mount_commands` (array of strings) - A series of commands to execute
  after attaching the root volume and before mounting the chroot. This is not
  required unless using `from_scratch`. If so, this should include any
  partitioning and filesystem creation commands. The path to the device is
  provided by `{{.Device}}`.

- `post_mount_commands` (array of strings) - As `pre_mount_commands`, but the
  commands are executed after mounting the root device and before the extra
  mount and copy steps. The device and mount path are provided by
  `{{.Device}}` and `{{.MountPath}}`.

- `root_volume_size` (number) - The size of the root volume in GB for the
  chroot environment and the resulting OMI. Default size is the snapshot size
  of the `source_omi` unless `from_scratch` is `true`, in which case this
  field must be defined.

- `root_volume_type` (string) - The type of BSU volume for the chroot
  environment and resulting OMI. The default value is the type of the
  `source_omi`, unless `from_scratch` is `true`, in which case the default
  value is `gp2`. You can only specify `io1` if building based on top of a
  `source_omi` which is also `io1`.

- `root_volume_tags` (object of key/value strings) - Tags to apply to the
  volumes that are _launched_. This is a [template
  engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

- `skip_region_validation` (boolean) - Set to true if you want to skip
  validation of the region configuration option. Default `false`.

- `snapshot_tags` (object of key/value strings) - Tags to apply to snapshot.
  They will override OMI tags if already applied to snapshot. This is a
  [template engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

- `source_omi_filter` (object) - Filters used to populate the `source_omi` field.

  - `filters` (map of strings) - filters used to select a `source_omi`.
  - `owners` (array of strings) - Filters the images by their owner. You may specify one or more Outscale account IDs, "self" (which will use the account whose credentials you are using to run Packer) or an Outscale owner alias. This option is required for security reasons.

    Example:

    ```json
    {
      "source_omi_filter": {
        "filters": {
          "virtualization-type": "hvm",
          "image-name": "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*",
          "root-device-type": "ebs"
        },
        "owners": ["099720109477"]
      }
    }
    ```

    This selects an Ubuntu 16.04 HVM BSU OMIS from Canonical. NOTE:
    This will fail unless _exactly_ one OMIS is returned. In the above example,
    `most_recent` will cause this to succeed by selecting the newest image.

    You may set this in place of `source_omi` or in conjunction with it. If you
    set this in conjunction with `source_omi`, the `source_omi` will be added
    to the filter. The provided `source_omi` must meet all of the filtering
    criteria provided in `source_omi_filter`; this pins the OMI returned by the
    filter, but will cause Packer to fail if the `source_omi` does not exist.

- `tags` (object of key/value strings) - Tags applied to the OMI. This is a
  [template engine](/docs/templates/legacy_json_templates/engine), see [Build template
  data](#build-template-data) for more information.

## Basic Example

Here is a basic example. It is completely valid except for the access keys:

```json
{
  "type": "outscale-chroot",
  "access_key": "YOUR KEY HERE",
  "secret_key": "YOUR SECRET KEY HERE",
  "source_omi": "ami-3e158364",
  "omi_name": "packer-outscale-chroot {{timestamp}}"
}
```

## Chroot Mounts

The `chroot_mounts` configuration can be used to mount specific devices within
the chroot. By default, the following additional mounts are added into the
chroot by Packer:

- `/proc` (proc)
- `/sys` (sysfs)
- `/dev` (bind to real `/dev`)
- `/dev/pts` (devpts)
- `/proc/sys/fs/binfmt_misc` (binfmt_misc)

These default mounts are usually good enough for anyone and are sane defaults.
However, if you want to change or add the mount points, you may using the
`chroot_mounts` configuration. Here is an example configuration which only
mounts `/proc` and `/dev`:

```json
{
  "chroot_mounts": [
    ["proc", "proc", "/proc"],
    ["bind", "/dev", "/dev"]
  ]
}
```

`chroot_mounts` is a list of a 3-tuples of strings. The three components of the
3-tuple, in order, are:

- The filesystem type. If this is "bind", then Packer will properly bind the
  filesystem to another mount point.

- The source device.

- The mount directory.

## Parallelism

A quick note on parallelism: it is perfectly safe to run multiple _separate_
Packer processes with the `outscale-chroot` builder on the same Outscale VM. In
fact, this is recommended as a way to push the most performance out of your OMI
builds.

Packer properly obtains a process lock for the parallelism-sensitive parts of
its internals such as finding an available device.

## Gotchas

### Unmounting the Filesystem

One of the difficulties with using the chroot builder is that your provisioning
scripts must not leave any processes running or packer will be unable to
unmount the filesystem.

For debian based distributions you can setup a
[policy-rc.d](http://people.debian.org/~hmh/invokerc.d-policyrc.d-specification.txt)
file which will prevent packages installed by your provisioners from starting
services:

```json
({
  "type": "shell",
  "inline": [
    "echo '#!/bin/sh' > /usr/sbin/policy-rc.d",
    "echo 'exit 101' >> /usr/sbin/policy-rc.d",
    "chmod a+x /usr/sbin/policy-rc.d"
  ]
},
{
  "type": "shell",
  "inline": ["rm -f /usr/sbin/policy-rc.d"]
})
```

### Ansible provisioner

Running `ansible` against `outscale-chroot` requires changing the Ansible connection
to chroot and running Ansible as root/sudo.

## Building From Scratch

This example demonstrates the essentials of building an image from scratch. A
15G gp2 (SSD) device is created (overriding the default of standard/magnetic).
The device setup commands partition the device with one partition for use as an
HVM image and format it ext4. This builder block should be followed by
provisioning commands to install the os and bootloader.

```json
{
  "type": "outscale-chroot",
  "ami_name": "packer-from-scratch {{timestamp}}",
  "from_scratch": true,
  "ami_virtualization_type": "hvm",
  "pre_mount_commands": [
    "parted {{.Device}} mklabel msdos mkpart primary 1M 100% set 1 boot on print",
    "mkfs.ext4 {{.Device}}1"
  ],
  "root_volume_size": 15,
  "root_device_name": "xvdf",
  "ami_block_device_mappings": [
    {
      "device_name": "xvdf",
      "delete_on_termination": true,
      "volume_type": "gp2"
    }
  ]
}
```

## Build template data

In configuration directives marked as a template engine above, the following
variables are available:

- `BuildRegion` - The region (for example `eu-west-2`) where Packer is building the OMI.
- `SourceOMI` - The source OMIS ID (for example `ami-a2412fcd`) used to build the OMI.
- `SourceOMIName` - The source OMIS Name (for example `ubuntu-390`) used to build the OMI.
- `SourceOMITags` - The source OMIS Tags, as a `map[string]string` object
