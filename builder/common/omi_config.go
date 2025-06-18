package common

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

// OMIConfig is for common configuration related to creating OMIs.
type OMIConfig struct {
	OMIName                 string   `mapstructure:"omi_name"`
	OMIDescription          string   `mapstructure:"omi_description"`
	OMIAccountIDs           []string `mapstructure:"omi_account_ids"`
	OMIGroups               []string `mapstructure:"omi_groups"`
	OMIProductCodes         []string `mapstructure:"omi_product_codes"`
	OMIRegions              []string `mapstructure:"omi_regions"`
	OMIBootModes            []string `mapstructure:"omi_boot_modes"`
	OMISkipRegionValidation bool     `mapstructure:"skip_region_validation"`
	OMITags                 TagMap   `mapstructure:"tags"`
	OMIForceDeregister      bool     `mapstructure:"force_deregister"`
	OMIForceDeleteSnapshot  bool     `mapstructure:"force_delete_snapshot"`
	SnapshotTags            TagMap   `mapstructure:"snapshot_tags"`
	SnapshotAccountIDs      []string `mapstructure:"snapshot_account_ids"`
	GlobalPermission        bool     `mapstructure:"global_permission"`
	ProductCodes            []string `mapstructure:"product_codes"`
	RootDeviceName          string   `mapstructure:"root_device_name"`
}

func (c *OMIConfig) Prepare(accessConfig *AccessConfig, ctx *interpolate.Context) []error {
	var errs []error

	if c.OMIName == "" {
		errs = append(errs, errors.New("omi_name must be specified"))
	}

	errs = append(errs, c.prepareRegions(accessConfig)...)

	if len(c.OMIName) < 3 || len(c.OMIName) > 128 {
		errs = append(errs, errors.New("omi_name must be between 3 and 128 characters long"))
	}
	if len(c.OMIBootModes) > 0 {
		bootModesSupported := []oscgo.BootMode{"legacy", "uefi"}
		for _, booModeValue := range c.OMIBootModes {
			var bootMode oscgo.BootMode = (oscgo.BootMode)(booModeValue)
			if !slices.Contains(bootModesSupported, bootMode) {
				errs = append(errs, fmt.Errorf("the omi_boot_Modes ['%v'] is not supported yet", bootMode))
			}
		}
	}

	if c.OMIName != templateCleanResourceName(c.OMIName) {
		errs = append(errs, errors.New("the parameter 'OMIName' should only contain"+
			" alphanumeric characters, parentheses (()), square brackets ([]), spaces "+
			"( ), periods (.), slashes (/), dashes (-), single quotes ('), at-signs "+
			"(@), or underscores(_). You can use the `clean_omi_name` template "+
			"filter to automatically clean your omi name. "))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (c *OMIConfig) prepareRegions(accessConfig *AccessConfig) (errs []error) {
	if len(c.OMIRegions) > 0 {
		regionSet := make(map[string]struct{})
		regions := make([]string, 0, len(c.OMIRegions))

		for _, region := range c.OMIRegions {
			// If we already saw the region, then don't look again
			if _, ok := regionSet[region]; ok {
				continue
			}

			// Mark that we saw the region
			regionSet[region] = struct{}{}

			if (accessConfig != nil) && (region == accessConfig.RawRegion) {
				// make sure we don't try to copy to the region we originally
				// create the OMI in.
				log.Printf("Cannot copy OMI to OUTSCALE session region '%s', deleting it from `omi_regions`.", region)
				continue
			}
			regions = append(regions, region)
		}

		c.OMIRegions = regions
	}
	return errs
}
func (c *OMIConfig) GetBootModes() (bootModes []oscgo.BootMode) {
	if len(c.OMIBootModes) > 0 {
		for _, bootModeValue := range c.OMIBootModes {
			bootModes = append(bootModes, (oscgo.BootMode)(bootModeValue))
		}
	}
	return
}
