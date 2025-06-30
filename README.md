# Packer Plugin Outscale

[![Project Graduated](https://docs.outscale.com/fr/userguide/_images/Project-Graduated-green.svg)](https://docs.outscale.com/en/userguide/Open-Source-Projects.html)
[![](https://dcbadge.limes.pink/api/server/HUVtY5gT6s?style=flat&theme=default-inverted)](https://discord.gg/HUVtY5gT6s)

<p align="center">
  <img alt="Packer" src="https://github.com/hashicorp/packer/raw/main/website/public/img/logo-packer-padded.svg" width="120px">
</p>

---

## 🌐 Links

* 📘 Plugin Documentation: [Packer Plugin Registry](https://developer.hashicorp.com/packer/integrations/outscale/outscale/latest)
* 📬 Packer official site: [https://www.packer.io](https://www.packer.io)
* 🤝 Contribution Guide: [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## 📄 Table of Contents

* [Overview](#-overview)
* [Requirements](#-requirements)
* [Installation](#-installation)

  * [Using `packer init`](#-using-packer-init)
  * [Manual Installation](#-manual-installation)
  * [From Source](#-from-source)
* [Configuration](#-configuration)
* [Contributing](#-contributing)

---

## 🧭 Overview

The **Outscale Packer Plugin** is a multi-component plugin for [HashiCorp Packer](https://www.packer.io) that enables users to build custom machine images on the [OUTSCALE](https://www.outscale.com) cloud platform.

It supports all standard Packer workflows and integrates seamlessly with the Packer CLI to provision images using your Outscale credentials and regions.

For a full list of features, refer to the [official plugin documentation](https://developer.hashicorp.com/packer/integrations/outscale/outscale/latest).

---

## ✅ Requirements

* [Packer 1.7+](https://www.packer.io/downloads)
* An active OUTSCALE account with access/secret keys
* Go 1.20+ (only if building from source)

---

## 🔨 Installation

### 🧪 Using `packer init`

Packer 1.7+ supports plugin installation via `packer init`.

Add the plugin declaration to your Packer configuration file:

```hcl
packer {
  required_plugins {
    outscale = {
      version = ">= 1.0.0"
      source  = "github.com/outscale/outscale"
    }
  }
}
```

Then initialize:

```bash
packer init .
```

📘 See [Packer plugin installation docs](https://www.packer.io/docs/commands/init) for more details.

---

### 🛠 Manual Installation

1. Download the appropriate binary from the [Releases page](https://github.com/outscale/packer-plugin-outscale/releases).
2. Uncompress the archive to retrieve the `packer-plugin-outscale` binary.
3. Follow the [Packer plugin installation guide](https://www.packer.io/docs/extending/plugins#installing-plugins) to move the binary to the appropriate plugin directory.

---

### 🧬 From Source

To build the plugin from source:

```bash
git clone https://github.com/outscale/packer-plugin-outscale.git
cd packer-plugin-outscale
go build
```

The resulting binary `packer-plugin-outscale` will be available in the root directory.

Follow the [manual plugin installation guide](https://www.packer.io/docs/extending/plugins#installing-plugins) to install the plugin into your environment.

---

## 🔧 Configuration

Configuration options and usage examples are documented in the [`docs/`](./docs) directory and on the [official registry page](https://developer.hashicorp.com/packer/integrations/outscale/outscale/latest).

---

## 🤝 Contributing

We welcome contributions!

* Found a bug or issue? Please [open an issue](https://github.com/outscale/packer-plugin-outscale/issues).
* Want to propose a new feature or fix? Start by opening an issue for discussion, then submit a [Pull Request](https://github.com/outscale/packer-plugin-outscale/pulls).

Please read our [CONTRIBUTING.md](./.github/CONTRIBUTING.md) for guidelines.
