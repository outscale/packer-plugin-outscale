
{
	"builders": [{
		"type": "outscale-bsusurrogate",
		"region": "eu-west-2",
		"vm_type": "tinav5.c1r1p1",
		"source_omi_filter": {
			"filters" = {
				"name"                = "Debian-12-*"
			}
			owners = ["Outscale"]
			most_recent = true
		},
		"ssh_username": "outscale",
		"omi_name": "packer-test {{timestamp}}",
		"subregion_name": "eu-west-2a",
		"launch_block_device_mappings" : [
			{
			"volume_type" : "io1",
			"device_name" : "/dev/xvdf",
			"delete_on_vm_deletion" : false,
			"volume_size" : 10,
			"iops": 300
			}
		],
		"omi_root_device":{
			"source_device_name": "/dev/xvdf",
			"device_name": "/dev/sda1",
			"delete_on_vm_deletion": true,
			"volume_size": 10,
			"volume_type": "standard"
		}

	}]
}
