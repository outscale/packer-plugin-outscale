package common

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/outscale/osc-sdk-go/v2"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/packer-plugin-outscale/version"
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
	SkipValidation        bool   `mapstructure:"skip_region_validation"`
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
	Auth context.Context
	Api  *oscgo.APIClient
}

// NewOSCClient retrieves the Outscale OSC-SDK client
func (c *AccessConfig) NewOSCClient() (*OscClient, error) {
	if c.AccessKey == "" {
		var ok bool
		if c.AccessKey, ok = getValueFromEnvVariables([]string{"OSC_ACCESS_KEY", "OUTSCALE_ACCESSKEYID"}); !ok {
			return nil, errors.New("No access key has been setted (configuration file, environment variable : OSC_ACCESS_KEY or OUTSCALE_ACCESSKEYID)")
		}
	}

	if c.SecretKey == "" {
		var ok bool
		if c.SecretKey, ok = getValueFromEnvVariables([]string{"OSC_SECRET_KEY", "OUTSCALE_SECRETKEYID"}); !ok {
			return nil, errors.New("No secret key has been setted (configuration file, environment variable : OSC_SECRET_KEY or OUTSCALE_SECRETKEYID)")
		}
	}

	if c.RawRegion == "" {
		var ok bool
		if c.RawRegion, ok = getValueFromEnvVariables([]string{"OSC_REGION", "OUTSCALE_REGION"}); !ok {
			return nil, errors.New("No region has been setted (configuration file, environment variable : OSC_REGION or OUTSCALE_REGION)")
		}
	}

	if c.CustomEndpointOAPI == "" {
		var ok bool
		if c.CustomEndpointOAPI, ok = getValueFromEnvVariables([]string{"OSC_ENDPOINT_API", "OUTSCALE_OAPI_URL"}); !ok {
			log.Printf("No Custom Endpoint has been setted")
		}
	}

	if c.CustomEndpointOAPI == "" {
		c.CustomEndpointOAPI = "outscale.com/oapi/latest"

		if c.RawRegion == "cn-southeast-1" {
			c.CustomEndpointOAPI = "outscale.hk/oapi/latest"
		}

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
	return c.NewOSCClientByRegion(c.RawRegion), nil
}

// GetRegion retrieves the Outscale OSC-SDK Region set
func (c *AccessConfig) GetRegion() string {
	return c.RawRegion
}

// NewOSCClientByRegion returns the connection depdending of the region given
func (c *AccessConfig) NewOSCClientByRegion(region string) *OscClient {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.InsecureSkipTLSVerify},
		Proxy:           http.ProxyFromEnvironment,
	}

	if c.X509certPath != "" && c.X509keyPath != "" {
		cert, err := tls.LoadX509KeyPair(c.X509certPath, c.X509keyPath)
		if err == nil {
			transport.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: c.InsecureSkipTLSVerify,
				Certificates:       []tls.Certificate{cert},
			}
		}
	}

	skipClient := &http.Client{
		Transport: transport,
	}

	skipClient.Transport = NewTransport(c.AccessKey, c.SecretKey, c.RawRegion, skipClient.Transport)

	config := oscgo.NewConfiguration()
	config.Debug = true
	config.UserAgent = fmt.Sprintf("packer-plugin-outscale/%s", version.PluginVersion.String())
	auth := context.WithValue(context.Background(), oscgo.ContextAWSv4, osc.AWSv4{
		AccessKey: os.Getenv("OSC_ACCESS_KEY"),
		SecretKey: os.Getenv("OSC_SECRET_KEY"),
	})
	config.HTTPClient = skipClient
	auth = context.WithValue(auth, oscgo.ContextServerIndex, 0)
	auth = context.WithValue(auth, oscgo.ContextServerVariables, map[string]string{"region": os.Getenv("OSC_REGION")})

	oscClient := &OscClient{
		Api:  oscgo.NewAPIClient(config),
		Auth: auth,
	}
	return oscClient
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
			fmt.Errorf("`access_key` and `secret_key` must both be either set or not set."))
	}

	return errs
}
