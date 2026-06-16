package common_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
	"github.com/outscale/packer-plugin-outscale/builder/common"
)

func TestAccessconfigConfigPriority(t *testing.T) {
	t.Setenv("OSC_ACCESS_KEY", "env-ak")
	t.Setenv("OSC_SECRET_KEY", "env-sk")
	t.Setenv("OSC_REGION", "env-region")
	t.Setenv("OUTSCALE_ACCESSKEYID", "deprecated-ak")
	t.Setenv("OUTSCALE_SECRETKEYID", "deprecated-sk")
	t.Setenv("OUTSCALE_REGION", "deprecated-region")

	c := &common.AccessConfig{
		AccessKey: "config-ak",
		SecretKey: "config-sk",
		RawRegion: "config-region",
	}

	p, err := c.NewProfile()
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}
	if p.AccessKey != "config-ak" {
		t.Fatalf("expected config access key, got %q", p.AccessKey)
	}
	if p.SecretKey != "config-sk" {
		t.Fatalf("expected config secret key, got %q", p.SecretKey)
	}
	if p.Region != "config-region" {
		t.Fatalf("expected config region, got %q", p.Region)
	}
}

func TestAccessconfigProfilePriority(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	cf := profile.ConfigFile{
		Path: configPath,
		Profiles: map[string]profile.Profile{
			"default": {
				AccessKey: "default-ak",
				SecretKey: "default-sk",
			},
			"test": {
				AccessKey: "file-ak",
				SecretKey: "file-sk",
			},
		},
	}
	if err := cf.Save(); err != nil {
		t.Fatalf("save config: %v", err)
	}

	t.Setenv("OSC_CONFIG_FILE", configPath)
	t.Setenv("OSC_ACCESS_KEY", "env-ak")
	t.Setenv("OSC_SECRET_KEY", "env-sk")
	t.Setenv("OSC_REGION", "env-region")
	t.Setenv("OUTSCALE_ACCESSKEYID", "deprecated-ak")
	t.Setenv("OUTSCALE_SECRETKEYID", "deprecated-sk")
	t.Setenv("OUTSCALE_REGION", "deprecated-region")

	c := &common.AccessConfig{ProfileName: "test"}

	p, err := c.NewProfile()
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}
	if p.AccessKey != "file-ak" {
		t.Fatalf("expected file access key, got %q", p.AccessKey)
	}
	if p.SecretKey != "file-sk" {
		t.Fatalf("expected file secret key, got %q", p.SecretKey)
	}
	if p.Region != "env-region" {
		t.Fatalf("expected env region, got %q", p.Region)
	}
}

func TestAccessconfigEnvPriority(t *testing.T) {
	t.Setenv("OSC_ACCESS_KEY", "env-ak")
	t.Setenv("OSC_SECRET_KEY", "env-sk")
	t.Setenv("OSC_REGION", "env-region")
	t.Setenv("OUTSCALE_ACCESSKEYID", "deprecated-ak")
	t.Setenv("OUTSCALE_SECRETKEYID", "deprecated-sk")
	t.Setenv("OUTSCALE_REGION", "deprecated-region")

	c := &common.AccessConfig{}

	p, err := c.NewProfile()
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}
	if p.AccessKey != "env-ak" {
		t.Fatalf("expected env access key, got %q", p.AccessKey)
	}
	if p.SecretKey != "env-sk" {
		t.Fatalf("expected env secret key, got %q", p.SecretKey)
	}
	if p.Region != "env-region" {
		t.Fatalf("expected env region, got %q", p.Region)
	}
}

func TestAccessconfigDeprecatedEnv(t *testing.T) {
	// Unset env vars that would override deprecated env vars
	t.Setenv("OSC_ACCESS_KEY", "")
	t.Setenv("OSC_SECRET_KEY", "")
	t.Setenv("OSC_REGION", "")
	t.Setenv("OSC_PROFILE", "")
	os.Unsetenv("OSC_ACCESS_KEY")
	os.Unsetenv("OSC_SECRET_KEY")
	os.Unsetenv("OSC_REGION")
	os.Unsetenv("OSC_PROFILE")

	t.Setenv("OUTSCALE_ACCESSKEYID", "deprecated-ak")
	t.Setenv("OUTSCALE_SECRETKEYID", "deprecated-sk")
	t.Setenv("OUTSCALE_REGION", "deprecated-region")

	c := &common.AccessConfig{}

	p, err := c.NewProfile()
	if err != nil {
		t.Fatalf("NewProfile: %v", err)
	}
	if p.AccessKey != "deprecated-ak" {
		t.Fatalf("expected deprecated access key, got %q", p.AccessKey)
	}
	if p.SecretKey != "deprecated-sk" {
		t.Fatalf("expected deprecated secret key, got %q", p.SecretKey)
	}
	if p.Region != "deprecated-region" {
		t.Fatalf("expected deprecated region, got %q", p.Region)
	}
}
