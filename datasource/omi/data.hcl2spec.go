// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package omi

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName       *string                `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType     *string                `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion     *string                `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug           *bool                  `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce           *bool                  `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError         *string                `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars        map[string]string      `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars   []string               `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
	AccessKey             *string                `mapstructure:"access_key" cty:"access_key" hcl:"access_key"`
	CustomEndpointOAPI    *string                `mapstructure:"custom_endpoint_oapi" cty:"custom_endpoint_oapi" hcl:"custom_endpoint_oapi"`
	InsecureSkipTLSVerify *bool                  `mapstructure:"insecure_skip_tls_verify" cty:"insecure_skip_tls_verify" hcl:"insecure_skip_tls_verify"`
	MFACode               *string                `mapstructure:"mfa_code" cty:"mfa_code" hcl:"mfa_code"`
	ProfileName           *string                `mapstructure:"profile" cty:"profile" hcl:"profile"`
	RawRegion             *string                `mapstructure:"region" cty:"region" hcl:"region"`
	SecretKey             *string                `mapstructure:"secret_key" cty:"secret_key" hcl:"secret_key"`
	SkipMetadataApiCheck  *bool                  `mapstructure:"skip_metadata_api_check" cty:"skip_metadata_api_check" hcl:"skip_metadata_api_check"`
	Token                 *string                `mapstructure:"token" cty:"token" hcl:"token"`
	X509certPath          *string                `mapstructure:"x509_cert_path" cty:"x509_cert_path" hcl:"x509_cert_path"`
	X509keyPath           *string                `mapstructure:"x509_key_path" cty:"x509_key_path" hcl:"x509_key_path"`
	Filters               map[string]string      `cty:"filters" hcl:"filters"`
	Filter                []config.FlatNameValue `cty:"filter" hcl:"filter"`
	Owners                []string               `cty:"owners" hcl:"owners"`
	MostRecent            *bool                  `mapstructure:"most_recent" cty:"most_recent" hcl:"most_recent"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"access_key":                 &hcldec.AttrSpec{Name: "access_key", Type: cty.String, Required: false},
		"custom_endpoint_oapi":       &hcldec.AttrSpec{Name: "custom_endpoint_oapi", Type: cty.String, Required: false},
		"insecure_skip_tls_verify":   &hcldec.AttrSpec{Name: "insecure_skip_tls_verify", Type: cty.Bool, Required: false},
		"mfa_code":                   &hcldec.AttrSpec{Name: "mfa_code", Type: cty.String, Required: false},
		"profile":                    &hcldec.AttrSpec{Name: "profile", Type: cty.String, Required: false},
		"region":                     &hcldec.AttrSpec{Name: "region", Type: cty.String, Required: false},
		"secret_key":                 &hcldec.AttrSpec{Name: "secret_key", Type: cty.String, Required: false},
		"skip_metadata_api_check":    &hcldec.AttrSpec{Name: "skip_metadata_api_check", Type: cty.Bool, Required: false},
		"token":                      &hcldec.AttrSpec{Name: "token", Type: cty.String, Required: false},
		"x509_cert_path":             &hcldec.AttrSpec{Name: "x509_cert_path", Type: cty.String, Required: false},
		"x509_key_path":              &hcldec.AttrSpec{Name: "x509_key_path", Type: cty.String, Required: false},
		"filters":                    &hcldec.AttrSpec{Name: "filters", Type: cty.Map(cty.String), Required: false},
		"filter":                     &hcldec.BlockListSpec{TypeName: "filter", Nested: hcldec.ObjectSpec((*config.FlatNameValue)(nil).HCL2Spec())},
		"owners":                     &hcldec.AttrSpec{Name: "owners", Type: cty.List(cty.String), Required: false},
		"most_recent":                &hcldec.AttrSpec{Name: "most_recent", Type: cty.Bool, Required: false},
	}
	return s
}

// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatDatasourceOutput struct {
	ID           *string           `mapstructure:"id" cty:"id" hcl:"id"`
	Name         *string           `mapstructure:"name" cty:"name" hcl:"name"`
	CreationDate *string           `mapstructure:"creation_date" cty:"creation_date" hcl:"creation_date"`
	Owner        *string           `mapstructure:"owner" cty:"owner" hcl:"owner"`
	OwnerName    *string           `mapstructure:"owner_name" cty:"owner_name" hcl:"owner_name"`
	Tags         map[string]string `mapstructure:"tags" cty:"tags" hcl:"tags"`
}

// FlatMapstructure returns a new FlatDatasourceOutput.
// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*DatasourceOutput) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatDatasourceOutput)
}

// HCL2Spec returns the hcl spec of a DatasourceOutput.
// This spec is used by HCL to read the fields of DatasourceOutput.
// The decoded values from this spec will then be applied to a FlatDatasourceOutput.
func (*FlatDatasourceOutput) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"id":            &hcldec.AttrSpec{Name: "id", Type: cty.String, Required: false},
		"name":          &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
		"creation_date": &hcldec.AttrSpec{Name: "creation_date", Type: cty.String, Required: false},
		"owner":         &hcldec.AttrSpec{Name: "owner", Type: cty.String, Required: false},
		"owner_name":    &hcldec.AttrSpec{Name: "owner_name", Type: cty.String, Required: false},
		"tags":          &hcldec.AttrSpec{Name: "tags", Type: cty.Map(cty.String), Required: false},
	}
	return s
}
