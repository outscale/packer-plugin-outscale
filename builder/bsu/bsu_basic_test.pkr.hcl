
{
	"builders": [{
		"type": "outscale-bsu",
		"region": "eu-west-2",
		"vm_type": "tinav5.c1r1p1",
		"source_omi_filter": {
        "filters": {
          "image-name": "Debian-12-*"
        },
        "owners": ["Outscale"]
        "most_recent": true
        },
		"ssh_username": "outscale",
		"omi_name": "packer-test",
		"associate_public_ip_address": true,
		"force_deregister": true
	}]
}
