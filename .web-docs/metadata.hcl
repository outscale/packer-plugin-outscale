# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Outscale"
  description = "TODO"
  identifier = "packer/BrandonRomano/outscale"
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
