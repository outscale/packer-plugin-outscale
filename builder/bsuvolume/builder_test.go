package bsuvolume_test

import (
	"testing"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/outscale/packer-plugin-outscale/builder/bsuvolume"
)

func testConfig() map[string]any {
	return map[string]any{
		"access_key":   "foo",
		"secret_key":   "bar",
		"source_omi":   "foo",
		"vm_type":      "foo",
		"region":       "us-east-1",
		"ssh_username": "root",
	}
}

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw any = &bsuvolume.Builder{}
	if _, ok := raw.(packersdk.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}

func TestBuilder_Prepare_BadType(t *testing.T) {
	b := &bsuvolume.Builder{}
	c := map[string]any{
		"access_key": []string{},
	}

	_, warnings, err := b.Prepare(c)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("prepare should fail")
	}
}

func TestBuilderPrepare_InvalidKey(t *testing.T) {
	var b bsuvolume.Builder
	config := testConfig()

	// Add a random key
	config["i_should_not_be_valid"] = true
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_InvalidShutdownBehavior(t *testing.T) {
	var b bsuvolume.Builder
	config := testConfig()

	// Test good
	config["shutdown_behavior"] = "terminate"
	_, warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test good
	config["shutdown_behavior"] = "stop"
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test bad
	config["shutdown_behavior"] = "foobar"
	_, warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}
