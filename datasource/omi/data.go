//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type DatasourceOutput,Config
package omi

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	osccommon "github.com/outscale/packer-plugin-outscale/builder/common"
	"github.com/zclconf/go-cty/cty"
)

type Datasource struct {
	config Config
}

type Config struct {
	common.PackerConfig        `mapstructure:",squash"`
	osccommon.AccessConfig     `mapstructure:",squash"`
	osccommon.OmiFilterOptions `mapstructure:",squash"`

	ctx interpolate.Context
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	var errs *packersdk.MultiError
	errs = packersdk.MultiErrorAppend(errs, d.config.AccessConfig.Prepare(&d.config.ctx)...)

	if len(d.config.Owners) == 0 && d.config.NameValueFilter.Empty() {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("The `filters` must be specified"))
	}
	if len(d.config.Owners) == 0 {
		errs = packersdk.MultiErrorAppend(
			errs,
			fmt.Errorf("For security reasons, you must declare an owner."),
		)
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

type DatasourceOutput struct {
	// The ID of the OMI.
	ID string `mapstructure:"id"`
	// The name of the OMI.
	Name string `mapstructure:"name"`
	// The date of creation of the OMI.
	CreationDate string `mapstructure:"creation_date"`
	// The Outscale account Id of the owner.
	Owner string `mapstructure:"owner"`
	// The owner alias.
	OwnerName string `mapstructure:"owner_name"`
	// The key/value combination of the tags assigned to the OMI.
	Tags map[string]string `mapstructure:"tags"`
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	accessConfig := &osccommon.AccessConfig{}
	oscConn, err := accessConfig.NewOSCClient()
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}
	image, err := d.config.OmiFilterOptions.GetFilteredImage(oscgo.ReadImagesRequest{}, oscConn)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	imageTags := make(map[string]string, len(*image.Tags))
	for _, tag := range *image.Tags {
		imageTags[tag.Key] = tag.Value
	}

	output := DatasourceOutput{
		ID:           image.ImageId,
		Name:         *image.ImageName,
		CreationDate: image.CreationDate.String(),
		Owner:        image.AccountId,
		OwnerName:    *image.AccountAlias,
		Tags:         imageTags,
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
