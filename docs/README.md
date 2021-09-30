# Outscale Plugin

The Outscale Packer Plugin is able to create Outscale OMIs. To achieve this, the plugin comes with
multiple builders depending on the strategy you want to use to build the OMI.
Packer supports the following builders at the moment:

- [osc-bsu](/docs/builders/osc-bsu) - Create BSU-backed OMIs by
  launching a source OMI and re-packaging it into a new OMI after
  provisioning. If in doubt, use this builder, which is the easiest to get
  started with.

- [osc-chroot](/docs/builders/osc-chroot) - Create EBS-backed OMIs
  from an existing OUTSCALE VM by mounting the root device and using a
  [Chroot](https://en.wikipedia.org/wiki/Chroot) environment to provision
  that device. This is an **advanced builder and should not be used by
  newcomers**. However, it is also the fastest way to build an EBS-backed OMI
  since no new OUTSCALE VM needs to be launched.

- [osc-bsusurrogate](/docs/builders/osc-bsusurrogate) - Create BSU-backed OMIs from scratch. Works similarly to the `chroot` builder but does
  not require running in Outscale VM. This is an **advanced builder and should not be
  used by newcomers**.

-> **Don't know which builder to use?** If in doubt, use the [osc-bsu
builder](/docs/builders/osc-bsu). It is much easier to use and Outscale generally recommends BSU-backed images nowadays.

### Outscale BSU Volume Builder

Packer is able to create Outscale BSU Volumes which are preinitialized with a filesystem and data.

- [osc-bsuvolume](/docs/builders/osc-bsuvolume) - Create EBS volumes by launching a source OMI with block devices mapped. Provision the VM, then destroy it, retaining the EBS volumes.


## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    outscale = {
      version = ">= 1.0.0"
      source  = "github.com/hashicorp/outscale"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/hashicorp/packer-plugin-name/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


#### From Source

If you prefer to build the plugin from its source code, clone the GitHub
repository locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-outscale` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).
