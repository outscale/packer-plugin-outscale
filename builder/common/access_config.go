package common

import (
	"errors"
	"fmt"
	"log"
	"os"

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

func getValueFromEnvVariables(envVariables []string) (string, bool) {
	for _, envVariable := range envVariables {
		if value, ok := os.LookupEnv(envVariable); ok && value != "" {
			return value, true
		}
	}

	return "", false
}

type OscClient struct {
	*oscgo.Client
}

// NewOSCClient retrieves the Outscale OSC-SDK client
func (c *AccessConfig) NewOSCClient() (*OscClient, error) {
	if c.AccessKey == "" {
		var ok bool
		if c.AccessKey, ok = getValueFromEnvVariables([]string{"OSC_ACCESS_KEY", "OUTSCALE_ACCESSKEYID"}); !ok {
			return nil, errors.New(
				"no access key has been setted (configuration file, environment variable : OSC_ACCESS_KEY or OUTSCALE_ACCESSKEYID)",
			)
		}
	}

	if c.SecretKey == "" {
		var ok bool
		if c.SecretKey, ok = getValueFromEnvVariables([]string{"OSC_SECRET_KEY", "OUTSCALE_SECRETKEYID"}); !ok {
			return nil, errors.New(
				"no secret key has been setted (configuration file, environment variable : OSC_SECRET_KEY or OUTSCALE_SECRETKEYID)",
			)
		}
	}

	if c.RawRegion == "" {
		var ok bool
		if c.RawRegion, ok = getValueFromEnvVariables([]string{"OSC_REGION", "OUTSCALE_REGION"}); !ok {
			return nil, errors.New(
				"no region has been setted (configuration file, environment variable : OSC_REGION or OUTSCALE_REGION)",
			)
		}
	}

	if c.CustomEndpointOAPI == "" {
		var ok bool
		if c.CustomEndpointOAPI, ok = getValueFromEnvVariables([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); !ok {
			log.Printf("No Custom Endpoint has been setted")
		}
	}

	if c.RawRegion == "cn-southeast-1" {
		c.CustomEndpointOAPI = fmt.Sprintf("https://api.%s.outscale.hk/api/v1", c.RawRegion)
	}

	if c.X509certPath == "" {
		var ok bool
		if c.X509certPath, ok = getValueFromEnvVariables([]string{"OSC_X509_CLIENT_CERT", "OUTSCALE_X509CERT"}); !ok {
			log.Printf("No Certificat Path has been setted")
		}
	}

	if c.X509keyPath == "" {
		var ok bool
		if c.X509certPath, ok = getValueFromEnvVariables([]string{"OSC_X509_CLIENT_KEY", "OUTSCALE_X509KEY"}); !ok {
			log.Printf("No Key Path has been setted")
		}
	}
	return c.NewOSCClientByRegion(c.RawRegion)
}

// GetRegion retrieves the Outscale OSC-SDK Region set
func (c *AccessConfig) GetRegion() string {
	return c.RawRegion
}

// NewOSCClientByRegion returns the connection depdending of the region given
func (c *AccessConfig) NewOSCClientByRegion(region string) (*OscClient, error) {
	profile := profile.Profile{
		Region:         c.RawRegion,
		AccessKey:      c.AccessKey,
		SecretKey:      c.SecretKey,
		X509ClientCert: c.X509certPath,
		X509ClientKey:  c.X509keyPath,
		Endpoints: profile.Endpoint{
			API: c.CustomEndpointOAPI,
		},
	}

	client, err := oscgo.NewClient(&profile)
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
