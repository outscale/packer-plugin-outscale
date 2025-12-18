//go:generate go run -modfile=../../go.mod github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc mapstructure-to-hcl2 -type SecurityGroupFilterOptions,OmiFilterOptions,SubnetFilterOptions,NetFilterOptions,BlockDevice

package common

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// ShutDown behavior possible
var (
	StopShutdownBehavior      string = "stop"
	TerminateShutdownBehavior string = "terminate"
	reShutdownBehavior               = regexp.MustCompile(
		"^(" + StopShutdownBehavior + "|" + TerminateShutdownBehavior + ")$",
	)
)

// docs at
// https://docs.outscale.com/en/userguide/Getting-Information-About-Your-OMIs.html
type OmiFilterOptions struct {
	config.NameValueFilter `         mapstructure:",squash"`
	Owners                 []string
	MostRecent             bool `mapstructure:"most_recent"`
}

func (d *OmiFilterOptions) Empty() bool {
	return len(d.Owners) == 0 && d.NameValueFilter.Empty()
}

func (d *OmiFilterOptions) NoOwner() bool {
	return len(d.Owners) == 0
}

// docs at
// https://docs.outscale.com/en/userguide/Getting-Information-About-Your-Subnets.html
type SubnetFilterOptions struct {
	config.NameValueFilter `     mapstructure:",squash"`
	MostFree               bool `mapstructure:"most_free"`
	Random                 bool `mapstructure:"random"`
}

// docs at https://docs.outscale.com/api#tocsfiltersnet
type NetFilterOptions struct {
	config.NameValueFilter `mapstructure:",squash"`
}

// docs at
// https://docs.outscale.com/en/userguide/Getting-Information-About-Your-Security-Groups.html
type SecurityGroupFilterOptions struct {
	config.NameValueFilter `mapstructure:",squash"`
}

// RunConfig contains configuration for running an vm from a source
// AMI and details on how to access that launched image.
type RunConfig struct {
	AssociatePublicIpAddress    bool                       `mapstructure:"associate_public_ip_address"`
	Subregion                   string                     `mapstructure:"subregion_name"`
	BlockDurationMinutes        int64                      `mapstructure:"block_duration_minutes"`
	DisableStopVm               bool                       `mapstructure:"disable_stop_vm"`
	BsuOptimized                bool                       `mapstructure:"bsu_optimized"`
	EnableT2Unlimited           bool                       `mapstructure:"enable_t2_unlimited"`
	IamVmProfile                string                     `mapstructure:"iam_vm_profile"`
	VmInitiatedShutdownBehavior string                     `mapstructure:"shutdown_behavior"`
	VmType                      string                     `mapstructure:"vm_type"`
	SecurityGroupFilter         SecurityGroupFilterOptions `mapstructure:"security_group_filter"`
	RunTags                     map[string]string          `mapstructure:"run_tags"`
	SecurityGroupId             string                     `mapstructure:"security_group_id"`
	SecurityGroupIds            []string                   `mapstructure:"security_group_ids"`
	SourceOmi                   string                     `mapstructure:"source_omi"`
	SourceOmiFilter             OmiFilterOptions           `mapstructure:"source_omi_filter"`
	SubnetFilter                SubnetFilterOptions        `mapstructure:"subnet_filter"`
	SubnetId                    string                     `mapstructure:"subnet_id"`
	TemporarySGSourceCidr       string                     `mapstructure:"temporary_security_group_source_cidr"`
	UserData                    string                     `mapstructure:"user_data"`
	UserDataFile                string                     `mapstructure:"user_data_file"`
	NetFilter                   NetFilterOptions           `mapstructure:"net_filter"`
	NetId                       string                     `mapstructure:"net_id"`
	WindowsPasswordTimeout      time.Duration              `mapstructure:"windows_password_timeout"`
	BootMode                    oscgo.BootMode             `mapstructure:"boot_mode"`
	// Communicator settings
	Comm         communicator.Config `mapstructure:",squash"`
	SSHInterface string              `mapstructure:"ssh_interface"`
}

func (c *RunConfig) Prepare(ctx *interpolate.Context) []error {
	// If we are not given an explicit ssh_keypair_name or
	// ssh_private_key_file, then create a temporary one, but only if the
	// temporary_key_pair_name has not been provided and we are not using
	// ssh_password.
	if c.Comm.SSHKeyPairName == "" && c.Comm.SSHTemporaryKeyPairName == "" &&
		c.Comm.SSHPrivateKeyFile == "" && c.Comm.SSHPassword == "" {

		c.Comm.SSHTemporaryKeyPairName = fmt.Sprintf("packer_%s", uuid.TimeOrderedUUID())
	}

	if c.WindowsPasswordTimeout == 0 {
		c.WindowsPasswordTimeout = 20 * time.Minute
	}

	if c.RunTags == nil {
		c.RunTags = make(map[string]string)
	}
	// Validation
	errs := c.Comm.Prepare(ctx)

	for _, preparer := range []interface{ Prepare() []error }{
		&c.SourceOmiFilter,
		&c.SecurityGroupFilter,
		&c.SubnetFilter,
		&c.NetFilter,
	} {
		errs = append(errs, preparer.Prepare()...)
	}

	// Validating ssh_interface
	if c.SSHInterface != "public_ip" &&
		c.SSHInterface != "private_ip" &&
		c.SSHInterface != "public_dns" &&
		c.SSHInterface != "private_dns" &&
		c.SSHInterface != "" {
		errs = append(errs, fmt.Errorf("unknown interface type: %s", c.SSHInterface))
	}

	if c.Comm.SSHKeyPairName != "" {
		if c.Comm.Type == "winrm" && c.Comm.WinRMPassword == "" && c.Comm.SSHPrivateKeyFile == "" {
			errs = append(
				errs,
				errors.New(
					"ssh_private_key_file must be provided to retrieve the winrm password when using ssh_keypair_name",
				),
			)
		} else if c.Comm.SSHPrivateKeyFile == "" && !c.Comm.SSHAgentAuth {
			errs = append(errs, errors.New("ssh_private_key_file must be provided or ssh_agent_auth enabled when ssh_keypair_name is specified"))
		}
	}
	if c.BootMode != "" {
		bootModesSupported := []oscgo.BootMode{"legacy", "uefi"}
		if !slices.Contains(bootModesSupported, c.BootMode) {
			errs = append(
				errs,
				fmt.Errorf("the vm boot_Mode '%v' is not supported yet", c.BootMode),
			)
		}
	}
	if c.SourceOmi == "" && c.SourceOmiFilter.Empty() {
		errs = append(errs, errors.New("a source_omi or source_omi_filter must be specified"))
	}

	if c.SourceOmi == "" && c.SourceOmiFilter.NoOwner() {
		errs = append(
			errs,
			errors.New("for security reasons, your source AMI filter must declare an owner"),
		)
	}

	if c.VmType == "" {
		errs = append(errs, errors.New("an vm_type must be specified"))
	}

	if c.BlockDurationMinutes%60 != 0 {
		errs = append(errs, errors.New(
			"block_duration_minutes must be multiple of 60"))
	}

	if c.UserData != "" && c.UserDataFile != "" {
		errs = append(errs, errors.New("only one of user_data or user_data_file can be specified"))
	} else if c.UserDataFile != "" {
		if _, err := os.Stat(c.UserDataFile); err != nil {
			errs = append(errs, fmt.Errorf("user_data_file not found: %s", c.UserDataFile))
		}
	}

	if c.SecurityGroupId != "" {
		if len(c.SecurityGroupIds) > 0 {
			errs = append(
				errs,
				errors.New("only one of security_group_id or security_group_ids can be specified"),
			)
		} else {
			c.SecurityGroupIds = []string{c.SecurityGroupId}
			c.SecurityGroupId = ""
		}
	}

	if c.TemporarySGSourceCidr == "" {
		c.TemporarySGSourceCidr = "0.0.0.0/0"
	} else {
		if _, _, err := net.ParseCIDR(c.TemporarySGSourceCidr); err != nil {
			errs = append(errs, fmt.Errorf("error parsing temporary_security_group_source_cidr: %w", err))
		}
	}

	if c.VmInitiatedShutdownBehavior == "" {
		c.VmInitiatedShutdownBehavior = StopShutdownBehavior
	} else if !reShutdownBehavior.MatchString(c.VmInitiatedShutdownBehavior) {
		errs = append(errs, errors.New("shutdown_behavior only accepts 'stop' or 'terminate' values"))
	}

	if c.EnableT2Unlimited {
		firstDotIndex := strings.Index(c.VmType, ".")
		if firstDotIndex == -1 {
			errs = append(errs, fmt.Errorf("error determining main Vm Type from: %s", c.VmType))
		} else if c.VmType[0:firstDotIndex] != "t2" {
			errs = append(errs, fmt.Errorf("error: T2 Unlimited enabled with a non-T2 Vm Type: %s", c.VmType))
		}
	}

	return errs
}
