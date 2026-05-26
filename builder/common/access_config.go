package common

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/osc-sdk-go/v3/pkg/profile"
)

// AccessConfig is for common configuration related to Outscale API access
type AccessConfig struct {
	AccessKey             string `mapstructure:"access_key"`
	CustomEndpointOAPI    string `mapstructure:"custom_endpoint_oapi"`
	InsecureSkipTLSVerify bool   `mapstructure:"insecure_skip_tls_verify"`
	MFACode               string `mapstructure:"mfa_code"`
	ProfileName           string `mapstructure:"profile"`
	RawRegion             string `mapstructure:"region"`
	SecretKey             string `mapstructure:"secret_key"`
	SkipMetadataApiCheck  bool   `mapstructure:"skip_metadata_api_check"`
	Token                 string `mapstructure:"token"`
	X509certPath          string `mapstructure:"x509_cert_path"`
	X509keyPath           string `mapstructure:"x509_key_path"`
}

type OscClient struct {
	*oscgo.Client
}

func (c *AccessConfig) GetRegion() string {
	return c.RawRegion
}

func (c *AccessConfig) NewOSCClient() (*OscClient, error) {
	profile, err := profile.NewFrom("", "")
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	profile.Protocol = "https"
	profile.TlsSkipVerify = c.InsecureSkipTLSVerify

	if c.RawRegion != "" {
		profile.Region = c.RawRegion
	}
	if c.AccessKey != "" {
		profile.AccessKey = c.AccessKey
	}
	if c.SecretKey != "" {
		profile.SecretKey = c.SecretKey
	}
	if c.X509certPath != "" {
		profile.X509ClientCert = c.X509certPath
	}
	if c.X509keyPath != "" {
		profile.X509ClientKey = c.X509keyPath
	}
	if c.CustomEndpointOAPI != "" {
		profile.Endpoints.API = c.CustomEndpointOAPI
	}

	client, err := oscgo.NewClient(profile)
	if err != nil {
		return nil, err
	}

	return &OscClient{
		Client: client,
	}, nil
}

func (c *AccessConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error

	if c.SkipMetadataApiCheck {
		log.Println("(WARN) skip_metadata_api_check ignored.")
	}
	// Either both access and secret key must be set or neither of them should
	// be.
	if (len(c.AccessKey) > 0) != (len(c.SecretKey) > 0) {
		errs = append(errs,
			errors.New("`access_key` and `secret_key` must both be either set or not set. "))
	}

	return errs
}
