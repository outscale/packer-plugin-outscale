package common

import (
	"testing"
)

func testAccessConfig() *AccessConfig {
	return &AccessConfig{}
}

func TestAccessConfigPrepare_Region(t *testing.T) {

	c := testAccessConfig()
	c.RawRegion = "eu-west-2"

	err := c.ValidateOSCRegion(c.RawRegion)
	if err != nil {
		t.Fatalf("should not have region validation for %s err: %s", c.RawRegion, err)
	}

	region := "us-east-12"
	err = c.ValidateOSCRegion(region)
	if err == nil {
		t.Fatalf("should have region validation err: %s", region)
	}

	region = "us-east-1"
	err = c.ValidateOSCRegion(region)
	if err == nil {
		t.Fatalf("should have region validation err: %s", region)
	}

	region = "custom"
	err = c.ValidateOSCRegion(region)
	if err == nil {
		t.Fatalf("should have region validation err: %s", region)
	}

	region = ""
	err = c.ValidateOSCRegion(region)
	if err == nil {
		t.Fatalf("should have region validation err: %s", region)
	}
	t.Logf("Error %v", err)

	region = "custom"
	c.SkipValidation = true
	// testing whole prepare func here; this is checking that validation is
	// skipped, so we don't need a mock connection
	if err := c.Prepare(nil); err != nil {
		t.Fatalf("shouldn't have err: %s", err)
	}

	c.SkipValidation = false
	region = ""
	if err := c.Prepare(nil); err != nil {
		t.Fatalf("shouldn't have err: %s", err)
	}
}
