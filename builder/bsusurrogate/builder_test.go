package bsusurrogate_test

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/outscale/packer-plugin-outscale/builder/bsusurrogate"
)

func testConfig() map[string]any {
	return map[string]any{
		"access_key":   "foo",
		"secret_key":   "bar",
		"source_omi":   "foo",
		"vm_type":      "foo",
		"region":       "us-east-1",
		"ssh_username": "root",
		"omi_name":     "foo",
		"omi_root_device": map[string]any{
			"device_name":        "/dev/sda1",
			"source_device_name": "/dev/xvdf",
		},
		"launch_block_device_mappings": map[string]any{
			"device_name": "/dev/xvdf",
		},
	}
}

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw any = &bsusurrogate.Builder{}
	if _, ok := raw.(packersdk.Builder); !ok {
		t.Fatal("Builder should be a builder")
	}
}

func TestBuilder_ShutdownBehavior_BsuDeletion(t *testing.T) {
	var b bsusurrogate.Builder
	config := testConfig()

	// Test good (terminate and keep bsu)
	config["shutdown_behavior"] = "terminate"
	config["launch_block_device_mappings"].(map[string]any)["delete_on_vm_deletion"] = false
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test KO (terminate and delete bsu with vm deletion)
	config["shutdown_behavior"] = "terminate"
	config["launch_block_device_mappings"].(map[string]any)["delete_on_vm_deletion"] = true
	_, _, err = b.Prepare(config)
	if err == nil {
		t.Fatalf("should  have failed")
	}

	// Test OK (stop and delete bsu with vm deletion)
	config["shutdown_behavior"] = "stop"
	config["launch_block_device_mappings"].(map[string]any)["delete_on_vm_deletion"] = true
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test OK (stop and keep bsu)
	config["shutdown_behavior"] = "stop"
	config["launch_block_device_mappings"].(map[string]any)["delete_on_vm_deletion"] = false
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}
}
