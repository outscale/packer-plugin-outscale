package common

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
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

func fromDeprecatedEnv() profile.Option {
	return func(p *profile.Profile) error {
		logDeprecation := func(deprecatedVar, replacementVar string) {
			log.Printf("%s environment variable is deprecated and support will be dropped in the next major version, use the %s environment variable instead", deprecatedVar, replacementVar)
		}

		if ak, ok := os.LookupEnv("OUTSCALE_ACCESSKEYID"); ok {
			p.AccessKey = ak
			logDeprecation("OUTSCALE_ACCESSKEYID", "OSC_ACCESS_KEY")
		}
		if sk, ok := os.LookupEnv("OUTSCALE_SECRETKEYID"); ok {
			p.SecretKey = sk
			logDeprecation("OUTSCALE_SECRETKEYID", "OSC_SECRET_KEY")
		}
		if cert, ok := os.LookupEnv("OUTSCALE_X509CERT"); ok {
			p.X509ClientCert = cert
			logDeprecation("OUTSCALE_X509CERT", "OSC_X509_CLIENT_CERT")
		}
		if key, ok := os.LookupEnv("OUTSCALE_X509KEY"); ok {
			p.X509ClientKey = key
			logDeprecation("OUTSCALE_X509KEY", "OSC_X509_CLIENT_KEY")
		}
		if region, ok := os.LookupEnv("OUTSCALE_REGION"); ok {
			p.Region = region
			logDeprecation("OUTSCALE_REGION", "OSC_REGION")
		}
		if endpoint, ok := os.LookupEnv("OUTSCALE_OAPI_URL"); ok {
			p.Endpoints.API = endpoint
			logDeprecation("OUTSCALE_OAPI_URL", "OSC_ENDPOINT_API")
		}

		return nil
	}
}

func fromConfig(c *AccessConfig) profile.Option {
	return func(profile *profile.Profile) error {
		profile.Protocol = "https"
		profile.TlsSkipVerify = c.InsecureSkipTLSVerify
		profile.Region = c.RawRegion
		profile.AccessKey = c.AccessKey
		profile.SecretKey = c.SecretKey
		profile.X509ClientCert = c.X509certPath
		profile.X509ClientKey = c.X509keyPath
		profile.Endpoints.API = c.CustomEndpointOAPI

		return nil
	}
}

func (c *AccessConfig) NewProfile() (*profile.Profile, error) {
	opts := []profile.Option{fromConfig(c), profile.MergeWith(profile.FromEnv()), profile.MergeWith(fromDeprecatedEnv())}
	if c.ProfileName != "" {
		opts = []profile.Option{fromConfig(c), profile.MergeWith(profile.FromFile(c.ProfileName, "")), profile.MergeWith(profile.FromEnv()), profile.MergeWith(fromDeprecatedEnv())}
	}

	return profile.New(opts...)
}

func (c *AccessConfig) NewOSCClient() (*OscClient, error) {
	profile, err := c.NewProfile()
	if err != nil {
		return nil, fmt.Errorf("new profile: %w", err)
	}

	client, err := oscgo.NewClient(profile, options.WithLogging(newLogger()))
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
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
