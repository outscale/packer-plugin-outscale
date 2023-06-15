Type: `outscale-omi`

The Outscale OMI Data source will filter and fetch an Outscale OMI, and output all the OMI information that will
be then available to use in the [Outscale builders](https://developer.hashicorp.com/packer/plugins/builder/outscale).

-> **Note:** Data sources is a feature exclusively available to HCL2 templates.

Basic example of usage:

```hcl
data "outscale-omi" "basic-example" {
    filters = {
        virtualization-type = "hvm"
        name = "ubuntu/images/*ubuntu-xenial-16.04-amd64-server-*"
        root-device-type = "ebs"
    }
    owners = ["099720109477"]
    most_recent = true
}
```
This selects the most recent Ubuntu 16.04 HVM EBS OMI from Canonical. Note that the data source will fail unless
*exactly* one OMI is returned. In the above example, `most_recent` will cause this to succeed by selecting the newest image.
