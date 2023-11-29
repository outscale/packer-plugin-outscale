# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Outscale"
  description = "Use Packer to create Outscale OMIs."
  identifier = "packer/outscale/outscale"
  component {
    type = "data-source"
    name = "Outscale OMI"
    slug = "omi"
  }
  component {
    type = "builder"
    name = "Outscale BSU Volume"
    slug = "outscale-bsuvolume"
  }
  component {
    type = "builder"
    name = "Outscale chroot"
    slug = "outscale-chroot"
  }
  component {
    type = "builder"
    name = "Outscale BSU Surrogate"
    slug = "outscale-bsusurrogate"
  }
  component {
    type = "builder"
    name = "Outscale BSU"
    slug = "outscale-bsu"
  }
}
