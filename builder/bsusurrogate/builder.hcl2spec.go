// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package bsusurrogate

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/outscale/packer-plugin-outscale/builder/common"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName             *string                                `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType           *string                                `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion           *string                                `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug                 *bool                                  `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce                 *bool                                  `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError               *string                                `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars              map[string]string                      `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars         []string                               `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
	AccessKey                   *string                                `mapstructure:"access_key" cty:"access_key" hcl:"access_key"`
	CustomEndpointOAPI          *string                                `mapstructure:"custom_endpoint_oapi" cty:"custom_endpoint_oapi" hcl:"custom_endpoint_oapi"`
	InsecureSkipTLSVerify       *bool                                  `mapstructure:"insecure_skip_tls_verify" cty:"insecure_skip_tls_verify" hcl:"insecure_skip_tls_verify"`
	MFACode                     *string                                `mapstructure:"mfa_code" cty:"mfa_code" hcl:"mfa_code"`
	ProfileName                 *string                                `mapstructure:"profile" cty:"profile" hcl:"profile"`
	RawRegion                   *string                                `mapstructure:"region" cty:"region" hcl:"region"`
	SecretKey                   *string                                `mapstructure:"secret_key" cty:"secret_key" hcl:"secret_key"`
	SkipMetadataApiCheck        *bool                                  `mapstructure:"skip_metadata_api_check" cty:"skip_metadata_api_check" hcl:"skip_metadata_api_check"`
	Token                       *string                                `mapstructure:"token" cty:"token" hcl:"token"`
	X509certPath                *string                                `mapstructure:"x509_cert_path" cty:"x509_cert_path" hcl:"x509_cert_path"`
	X509keyPath                 *string                                `mapstructure:"x509_key_path" cty:"x509_key_path" hcl:"x509_key_path"`
	AssociatePublicIpAddress    *bool                                  `mapstructure:"associate_public_ip_address" cty:"associate_public_ip_address" hcl:"associate_public_ip_address"`
	Subregion                   *string                                `mapstructure:"subregion_name" cty:"subregion_name" hcl:"subregion_name"`
	BlockDurationMinutes        *int64                                 `mapstructure:"block_duration_minutes" cty:"block_duration_minutes" hcl:"block_duration_minutes"`
	DisableStopVm               *bool                                  `mapstructure:"disable_stop_vm" cty:"disable_stop_vm" hcl:"disable_stop_vm"`
	BsuOptimized                *bool                                  `mapstructure:"bsu_optimized" cty:"bsu_optimized" hcl:"bsu_optimized"`
	EnableT2Unlimited           *bool                                  `mapstructure:"enable_t2_unlimited" cty:"enable_t2_unlimited" hcl:"enable_t2_unlimited"`
	IamVmProfile                *string                                `mapstructure:"iam_vm_profile" cty:"iam_vm_profile" hcl:"iam_vm_profile"`
	VmInitiatedShutdownBehavior *string                                `mapstructure:"shutdown_behavior" cty:"shutdown_behavior" hcl:"shutdown_behavior"`
	VmType                      *string                                `mapstructure:"vm_type" cty:"vm_type" hcl:"vm_type"`
	SecurityGroupFilter         *common.FlatSecurityGroupFilterOptions `mapstructure:"security_group_filter" cty:"security_group_filter" hcl:"security_group_filter"`
	RunTags                     map[string]string                      `mapstructure:"run_tags" cty:"run_tags" hcl:"run_tags"`
	SecurityGroupId             *string                                `mapstructure:"security_group_id" cty:"security_group_id" hcl:"security_group_id"`
	SecurityGroupIds            []string                               `mapstructure:"security_group_ids" cty:"security_group_ids" hcl:"security_group_ids"`
	SourceOmi                   *string                                `mapstructure:"source_omi" cty:"source_omi" hcl:"source_omi"`
	SourceOmiFilter             *common.FlatOmiFilterOptions           `mapstructure:"source_omi_filter" cty:"source_omi_filter" hcl:"source_omi_filter"`
	SubnetFilter                *common.FlatSubnetFilterOptions        `mapstructure:"subnet_filter" cty:"subnet_filter" hcl:"subnet_filter"`
	SubnetId                    *string                                `mapstructure:"subnet_id" cty:"subnet_id" hcl:"subnet_id"`
	TemporarySGSourceCidr       *string                                `mapstructure:"temporary_security_group_source_cidr" cty:"temporary_security_group_source_cidr" hcl:"temporary_security_group_source_cidr"`
	UserData                    *string                                `mapstructure:"user_data" cty:"user_data" hcl:"user_data"`
	UserDataFile                *string                                `mapstructure:"user_data_file" cty:"user_data_file" hcl:"user_data_file"`
	NetFilter                   *common.FlatNetFilterOptions           `mapstructure:"net_filter" cty:"net_filter" hcl:"net_filter"`
	NetId                       *string                                `mapstructure:"net_id" cty:"net_id" hcl:"net_id"`
	WindowsPasswordTimeout      *string                                `mapstructure:"windows_password_timeout" cty:"windows_password_timeout" hcl:"windows_password_timeout"`
	BootMode                    *string                                `mapstructure:"boot_mode" cty:"boot_mode" hcl:"boot_mode"`
	Type                        *string                                `mapstructure:"communicator" cty:"communicator" hcl:"communicator"`
	PauseBeforeConnect          *string                                `mapstructure:"pause_before_connecting" cty:"pause_before_connecting" hcl:"pause_before_connecting"`
	SSHHost                     *string                                `mapstructure:"ssh_host" cty:"ssh_host" hcl:"ssh_host"`
	SSHPort                     *int                                   `mapstructure:"ssh_port" cty:"ssh_port" hcl:"ssh_port"`
	SSHUsername                 *string                                `mapstructure:"ssh_username" cty:"ssh_username" hcl:"ssh_username"`
	SSHPassword                 *string                                `mapstructure:"ssh_password" cty:"ssh_password" hcl:"ssh_password"`
	SSHKeyPairName              *string                                `mapstructure:"ssh_keypair_name" undocumented:"true" cty:"ssh_keypair_name" hcl:"ssh_keypair_name"`
	SSHTemporaryKeyPairName     *string                                `mapstructure:"temporary_key_pair_name" undocumented:"true" cty:"temporary_key_pair_name" hcl:"temporary_key_pair_name"`
	SSHTemporaryKeyPairType     *string                                `mapstructure:"temporary_key_pair_type" cty:"temporary_key_pair_type" hcl:"temporary_key_pair_type"`
	SSHTemporaryKeyPairBits     *int                                   `mapstructure:"temporary_key_pair_bits" cty:"temporary_key_pair_bits" hcl:"temporary_key_pair_bits"`
	SSHCiphers                  []string                               `mapstructure:"ssh_ciphers" cty:"ssh_ciphers" hcl:"ssh_ciphers"`
	SSHClearAuthorizedKeys      *bool                                  `mapstructure:"ssh_clear_authorized_keys" cty:"ssh_clear_authorized_keys" hcl:"ssh_clear_authorized_keys"`
	SSHKEXAlgos                 []string                               `mapstructure:"ssh_key_exchange_algorithms" cty:"ssh_key_exchange_algorithms" hcl:"ssh_key_exchange_algorithms"`
	SSHPrivateKeyFile           *string                                `mapstructure:"ssh_private_key_file" undocumented:"true" cty:"ssh_private_key_file" hcl:"ssh_private_key_file"`
	SSHCertificateFile          *string                                `mapstructure:"ssh_certificate_file" cty:"ssh_certificate_file" hcl:"ssh_certificate_file"`
	SSHPty                      *bool                                  `mapstructure:"ssh_pty" cty:"ssh_pty" hcl:"ssh_pty"`
	SSHTimeout                  *string                                `mapstructure:"ssh_timeout" cty:"ssh_timeout" hcl:"ssh_timeout"`
	SSHWaitTimeout              *string                                `mapstructure:"ssh_wait_timeout" undocumented:"true" cty:"ssh_wait_timeout" hcl:"ssh_wait_timeout"`
	SSHAgentAuth                *bool                                  `mapstructure:"ssh_agent_auth" undocumented:"true" cty:"ssh_agent_auth" hcl:"ssh_agent_auth"`
	SSHDisableAgentForwarding   *bool                                  `mapstructure:"ssh_disable_agent_forwarding" cty:"ssh_disable_agent_forwarding" hcl:"ssh_disable_agent_forwarding"`
	SSHHandshakeAttempts        *int                                   `mapstructure:"ssh_handshake_attempts" cty:"ssh_handshake_attempts" hcl:"ssh_handshake_attempts"`
	SSHBastionHost              *string                                `mapstructure:"ssh_bastion_host" cty:"ssh_bastion_host" hcl:"ssh_bastion_host"`
	SSHBastionPort              *int                                   `mapstructure:"ssh_bastion_port" cty:"ssh_bastion_port" hcl:"ssh_bastion_port"`
	SSHBastionAgentAuth         *bool                                  `mapstructure:"ssh_bastion_agent_auth" cty:"ssh_bastion_agent_auth" hcl:"ssh_bastion_agent_auth"`
	SSHBastionUsername          *string                                `mapstructure:"ssh_bastion_username" cty:"ssh_bastion_username" hcl:"ssh_bastion_username"`
	SSHBastionPassword          *string                                `mapstructure:"ssh_bastion_password" cty:"ssh_bastion_password" hcl:"ssh_bastion_password"`
	SSHBastionInteractive       *bool                                  `mapstructure:"ssh_bastion_interactive" cty:"ssh_bastion_interactive" hcl:"ssh_bastion_interactive"`
	SSHBastionPrivateKeyFile    *string                                `mapstructure:"ssh_bastion_private_key_file" cty:"ssh_bastion_private_key_file" hcl:"ssh_bastion_private_key_file"`
	SSHBastionCertificateFile   *string                                `mapstructure:"ssh_bastion_certificate_file" cty:"ssh_bastion_certificate_file" hcl:"ssh_bastion_certificate_file"`
	SSHFileTransferMethod       *string                                `mapstructure:"ssh_file_transfer_method" cty:"ssh_file_transfer_method" hcl:"ssh_file_transfer_method"`
	SSHProxyHost                *string                                `mapstructure:"ssh_proxy_host" cty:"ssh_proxy_host" hcl:"ssh_proxy_host"`
	SSHProxyPort                *int                                   `mapstructure:"ssh_proxy_port" cty:"ssh_proxy_port" hcl:"ssh_proxy_port"`
	SSHProxyUsername            *string                                `mapstructure:"ssh_proxy_username" cty:"ssh_proxy_username" hcl:"ssh_proxy_username"`
	SSHProxyPassword            *string                                `mapstructure:"ssh_proxy_password" cty:"ssh_proxy_password" hcl:"ssh_proxy_password"`
	SSHKeepAliveInterval        *string                                `mapstructure:"ssh_keep_alive_interval" cty:"ssh_keep_alive_interval" hcl:"ssh_keep_alive_interval"`
	SSHReadWriteTimeout         *string                                `mapstructure:"ssh_read_write_timeout" cty:"ssh_read_write_timeout" hcl:"ssh_read_write_timeout"`
	SSHRemoteTunnels            []string                               `mapstructure:"ssh_remote_tunnels" cty:"ssh_remote_tunnels" hcl:"ssh_remote_tunnels"`
	SSHLocalTunnels             []string                               `mapstructure:"ssh_local_tunnels" cty:"ssh_local_tunnels" hcl:"ssh_local_tunnels"`
	SSHPublicKey                []byte                                 `mapstructure:"ssh_public_key" undocumented:"true" cty:"ssh_public_key" hcl:"ssh_public_key"`
	SSHPrivateKey               []byte                                 `mapstructure:"ssh_private_key" undocumented:"true" cty:"ssh_private_key" hcl:"ssh_private_key"`
	WinRMUser                   *string                                `mapstructure:"winrm_username" cty:"winrm_username" hcl:"winrm_username"`
	WinRMPassword               *string                                `mapstructure:"winrm_password" cty:"winrm_password" hcl:"winrm_password"`
	WinRMHost                   *string                                `mapstructure:"winrm_host" cty:"winrm_host" hcl:"winrm_host"`
	WinRMNoProxy                *bool                                  `mapstructure:"winrm_no_proxy" cty:"winrm_no_proxy" hcl:"winrm_no_proxy"`
	WinRMPort                   *int                                   `mapstructure:"winrm_port" cty:"winrm_port" hcl:"winrm_port"`
	WinRMTimeout                *string                                `mapstructure:"winrm_timeout" cty:"winrm_timeout" hcl:"winrm_timeout"`
	WinRMUseSSL                 *bool                                  `mapstructure:"winrm_use_ssl" cty:"winrm_use_ssl" hcl:"winrm_use_ssl"`
	WinRMInsecure               *bool                                  `mapstructure:"winrm_insecure" cty:"winrm_insecure" hcl:"winrm_insecure"`
	WinRMUseNTLM                *bool                                  `mapstructure:"winrm_use_ntlm" cty:"winrm_use_ntlm" hcl:"winrm_use_ntlm"`
	SSHInterface                *string                                `mapstructure:"ssh_interface" cty:"ssh_interface" hcl:"ssh_interface"`
	OMIMappings                 []common.FlatBlockDevice               `mapstructure:"omi_block_device_mappings" cty:"omi_block_device_mappings" hcl:"omi_block_device_mappings"`
	LaunchMappings              []common.FlatBlockDevice               `mapstructure:"launch_block_device_mappings" cty:"launch_block_device_mappings" hcl:"launch_block_device_mappings"`
	OMIName                     *string                                `mapstructure:"omi_name" cty:"omi_name" hcl:"omi_name"`
	OMIDescription              *string                                `mapstructure:"omi_description" cty:"omi_description" hcl:"omi_description"`
	OMIAccountIDs               []string                               `mapstructure:"omi_account_ids" cty:"omi_account_ids" hcl:"omi_account_ids"`
	OMIGroups                   []string                               `mapstructure:"omi_groups" cty:"omi_groups" hcl:"omi_groups"`
	OMIProductCodes             []string                               `mapstructure:"omi_product_codes" cty:"omi_product_codes" hcl:"omi_product_codes"`
	OMIRegions                  []string                               `mapstructure:"omi_regions" cty:"omi_regions" hcl:"omi_regions"`
	OMIBootModes                []string                               `mapstructure:"omi_boot_modes" cty:"omi_boot_modes" hcl:"omi_boot_modes"`
	OMISkipRegionValidation     *bool                                  `mapstructure:"skip_region_validation" cty:"skip_region_validation" hcl:"skip_region_validation"`
	OMITags                     common.TagMap                          `mapstructure:"tags" cty:"tags" hcl:"tags"`
	OMIForceDeregister          *bool                                  `mapstructure:"force_deregister" cty:"force_deregister" hcl:"force_deregister"`
	OMIForceDeleteSnapshot      *bool                                  `mapstructure:"force_delete_snapshot" cty:"force_delete_snapshot" hcl:"force_delete_snapshot"`
	SnapshotTags                common.TagMap                          `mapstructure:"snapshot_tags" cty:"snapshot_tags" hcl:"snapshot_tags"`
	SnapshotAccountIDs          []string                               `mapstructure:"snapshot_account_ids" cty:"snapshot_account_ids" hcl:"snapshot_account_ids"`
	GlobalPermission            *bool                                  `mapstructure:"global_permission" cty:"global_permission" hcl:"global_permission"`
	ProductCodes                []string                               `mapstructure:"product_codes" cty:"product_codes" hcl:"product_codes"`
	RootDeviceName              *string                                `mapstructure:"root_device_name" cty:"root_device_name" hcl:"root_device_name"`
	RootDevice                  *FlatRootBlockDevice                   `mapstructure:"omi_root_device" cty:"omi_root_device" hcl:"omi_root_device"`
	VolumeRunTags               common.TagMap                          `mapstructure:"run_volume_tags" cty:"run_volume_tags" hcl:"run_volume_tags"`
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
		"packer_build_name":                    &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":                  &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":                  &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":                         &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":                         &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":                      &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":                &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables":           &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"access_key":                           &hcldec.AttrSpec{Name: "access_key", Type: cty.String, Required: false},
		"custom_endpoint_oapi":                 &hcldec.AttrSpec{Name: "custom_endpoint_oapi", Type: cty.String, Required: false},
		"insecure_skip_tls_verify":             &hcldec.AttrSpec{Name: "insecure_skip_tls_verify", Type: cty.Bool, Required: false},
		"mfa_code":                             &hcldec.AttrSpec{Name: "mfa_code", Type: cty.String, Required: false},
		"profile":                              &hcldec.AttrSpec{Name: "profile", Type: cty.String, Required: false},
		"region":                               &hcldec.AttrSpec{Name: "region", Type: cty.String, Required: false},
		"secret_key":                           &hcldec.AttrSpec{Name: "secret_key", Type: cty.String, Required: false},
		"skip_metadata_api_check":              &hcldec.AttrSpec{Name: "skip_metadata_api_check", Type: cty.Bool, Required: false},
		"token":                                &hcldec.AttrSpec{Name: "token", Type: cty.String, Required: false},
		"x509_cert_path":                       &hcldec.AttrSpec{Name: "x509_cert_path", Type: cty.String, Required: false},
		"x509_key_path":                        &hcldec.AttrSpec{Name: "x509_key_path", Type: cty.String, Required: false},
		"associate_public_ip_address":          &hcldec.AttrSpec{Name: "associate_public_ip_address", Type: cty.Bool, Required: false},
		"subregion_name":                       &hcldec.AttrSpec{Name: "subregion_name", Type: cty.String, Required: false},
		"block_duration_minutes":               &hcldec.AttrSpec{Name: "block_duration_minutes", Type: cty.Number, Required: false},
		"disable_stop_vm":                      &hcldec.AttrSpec{Name: "disable_stop_vm", Type: cty.Bool, Required: false},
		"bsu_optimized":                        &hcldec.AttrSpec{Name: "bsu_optimized", Type: cty.Bool, Required: false},
		"enable_t2_unlimited":                  &hcldec.AttrSpec{Name: "enable_t2_unlimited", Type: cty.Bool, Required: false},
		"iam_vm_profile":                       &hcldec.AttrSpec{Name: "iam_vm_profile", Type: cty.String, Required: false},
		"shutdown_behavior":                    &hcldec.AttrSpec{Name: "shutdown_behavior", Type: cty.String, Required: false},
		"vm_type":                              &hcldec.AttrSpec{Name: "vm_type", Type: cty.String, Required: false},
		"security_group_filter":                &hcldec.BlockSpec{TypeName: "security_group_filter", Nested: hcldec.ObjectSpec((*common.FlatSecurityGroupFilterOptions)(nil).HCL2Spec())},
		"run_tags":                             &hcldec.AttrSpec{Name: "run_tags", Type: cty.Map(cty.String), Required: false},
		"security_group_id":                    &hcldec.AttrSpec{Name: "security_group_id", Type: cty.String, Required: false},
		"security_group_ids":                   &hcldec.AttrSpec{Name: "security_group_ids", Type: cty.List(cty.String), Required: false},
		"source_omi":                           &hcldec.AttrSpec{Name: "source_omi", Type: cty.String, Required: false},
		"source_omi_filter":                    &hcldec.BlockSpec{TypeName: "source_omi_filter", Nested: hcldec.ObjectSpec((*common.FlatOmiFilterOptions)(nil).HCL2Spec())},
		"subnet_filter":                        &hcldec.BlockSpec{TypeName: "subnet_filter", Nested: hcldec.ObjectSpec((*common.FlatSubnetFilterOptions)(nil).HCL2Spec())},
		"subnet_id":                            &hcldec.AttrSpec{Name: "subnet_id", Type: cty.String, Required: false},
		"temporary_security_group_source_cidr": &hcldec.AttrSpec{Name: "temporary_security_group_source_cidr", Type: cty.String, Required: false},
		"user_data":                            &hcldec.AttrSpec{Name: "user_data", Type: cty.String, Required: false},
		"user_data_file":                       &hcldec.AttrSpec{Name: "user_data_file", Type: cty.String, Required: false},
		"net_filter":                           &hcldec.BlockSpec{TypeName: "net_filter", Nested: hcldec.ObjectSpec((*common.FlatNetFilterOptions)(nil).HCL2Spec())},
		"net_id":                               &hcldec.AttrSpec{Name: "net_id", Type: cty.String, Required: false},
		"windows_password_timeout":             &hcldec.AttrSpec{Name: "windows_password_timeout", Type: cty.String, Required: false},
		"boot_mode":                            &hcldec.AttrSpec{Name: "boot_mode", Type: cty.String, Required: false},
		"communicator":                         &hcldec.AttrSpec{Name: "communicator", Type: cty.String, Required: false},
		"pause_before_connecting":              &hcldec.AttrSpec{Name: "pause_before_connecting", Type: cty.String, Required: false},
		"ssh_host":                             &hcldec.AttrSpec{Name: "ssh_host", Type: cty.String, Required: false},
		"ssh_port":                             &hcldec.AttrSpec{Name: "ssh_port", Type: cty.Number, Required: false},
		"ssh_username":                         &hcldec.AttrSpec{Name: "ssh_username", Type: cty.String, Required: false},
		"ssh_password":                         &hcldec.AttrSpec{Name: "ssh_password", Type: cty.String, Required: false},
		"ssh_keypair_name":                     &hcldec.AttrSpec{Name: "ssh_keypair_name", Type: cty.String, Required: false},
		"temporary_key_pair_name":              &hcldec.AttrSpec{Name: "temporary_key_pair_name", Type: cty.String, Required: false},
		"temporary_key_pair_type":              &hcldec.AttrSpec{Name: "temporary_key_pair_type", Type: cty.String, Required: false},
		"temporary_key_pair_bits":              &hcldec.AttrSpec{Name: "temporary_key_pair_bits", Type: cty.Number, Required: false},
		"ssh_ciphers":                          &hcldec.AttrSpec{Name: "ssh_ciphers", Type: cty.List(cty.String), Required: false},
		"ssh_clear_authorized_keys":            &hcldec.AttrSpec{Name: "ssh_clear_authorized_keys", Type: cty.Bool, Required: false},
		"ssh_key_exchange_algorithms":          &hcldec.AttrSpec{Name: "ssh_key_exchange_algorithms", Type: cty.List(cty.String), Required: false},
		"ssh_private_key_file":                 &hcldec.AttrSpec{Name: "ssh_private_key_file", Type: cty.String, Required: false},
		"ssh_certificate_file":                 &hcldec.AttrSpec{Name: "ssh_certificate_file", Type: cty.String, Required: false},
		"ssh_pty":                              &hcldec.AttrSpec{Name: "ssh_pty", Type: cty.Bool, Required: false},
		"ssh_timeout":                          &hcldec.AttrSpec{Name: "ssh_timeout", Type: cty.String, Required: false},
		"ssh_wait_timeout":                     &hcldec.AttrSpec{Name: "ssh_wait_timeout", Type: cty.String, Required: false},
		"ssh_agent_auth":                       &hcldec.AttrSpec{Name: "ssh_agent_auth", Type: cty.Bool, Required: false},
		"ssh_disable_agent_forwarding":         &hcldec.AttrSpec{Name: "ssh_disable_agent_forwarding", Type: cty.Bool, Required: false},
		"ssh_handshake_attempts":               &hcldec.AttrSpec{Name: "ssh_handshake_attempts", Type: cty.Number, Required: false},
		"ssh_bastion_host":                     &hcldec.AttrSpec{Name: "ssh_bastion_host", Type: cty.String, Required: false},
		"ssh_bastion_port":                     &hcldec.AttrSpec{Name: "ssh_bastion_port", Type: cty.Number, Required: false},
		"ssh_bastion_agent_auth":               &hcldec.AttrSpec{Name: "ssh_bastion_agent_auth", Type: cty.Bool, Required: false},
		"ssh_bastion_username":                 &hcldec.AttrSpec{Name: "ssh_bastion_username", Type: cty.String, Required: false},
		"ssh_bastion_password":                 &hcldec.AttrSpec{Name: "ssh_bastion_password", Type: cty.String, Required: false},
		"ssh_bastion_interactive":              &hcldec.AttrSpec{Name: "ssh_bastion_interactive", Type: cty.Bool, Required: false},
		"ssh_bastion_private_key_file":         &hcldec.AttrSpec{Name: "ssh_bastion_private_key_file", Type: cty.String, Required: false},
		"ssh_bastion_certificate_file":         &hcldec.AttrSpec{Name: "ssh_bastion_certificate_file", Type: cty.String, Required: false},
		"ssh_file_transfer_method":             &hcldec.AttrSpec{Name: "ssh_file_transfer_method", Type: cty.String, Required: false},
		"ssh_proxy_host":                       &hcldec.AttrSpec{Name: "ssh_proxy_host", Type: cty.String, Required: false},
		"ssh_proxy_port":                       &hcldec.AttrSpec{Name: "ssh_proxy_port", Type: cty.Number, Required: false},
		"ssh_proxy_username":                   &hcldec.AttrSpec{Name: "ssh_proxy_username", Type: cty.String, Required: false},
		"ssh_proxy_password":                   &hcldec.AttrSpec{Name: "ssh_proxy_password", Type: cty.String, Required: false},
		"ssh_keep_alive_interval":              &hcldec.AttrSpec{Name: "ssh_keep_alive_interval", Type: cty.String, Required: false},
		"ssh_read_write_timeout":               &hcldec.AttrSpec{Name: "ssh_read_write_timeout", Type: cty.String, Required: false},
		"ssh_remote_tunnels":                   &hcldec.AttrSpec{Name: "ssh_remote_tunnels", Type: cty.List(cty.String), Required: false},
		"ssh_local_tunnels":                    &hcldec.AttrSpec{Name: "ssh_local_tunnels", Type: cty.List(cty.String), Required: false},
		"ssh_public_key":                       &hcldec.AttrSpec{Name: "ssh_public_key", Type: cty.List(cty.Number), Required: false},
		"ssh_private_key":                      &hcldec.AttrSpec{Name: "ssh_private_key", Type: cty.List(cty.Number), Required: false},
		"winrm_username":                       &hcldec.AttrSpec{Name: "winrm_username", Type: cty.String, Required: false},
		"winrm_password":                       &hcldec.AttrSpec{Name: "winrm_password", Type: cty.String, Required: false},
		"winrm_host":                           &hcldec.AttrSpec{Name: "winrm_host", Type: cty.String, Required: false},
		"winrm_no_proxy":                       &hcldec.AttrSpec{Name: "winrm_no_proxy", Type: cty.Bool, Required: false},
		"winrm_port":                           &hcldec.AttrSpec{Name: "winrm_port", Type: cty.Number, Required: false},
		"winrm_timeout":                        &hcldec.AttrSpec{Name: "winrm_timeout", Type: cty.String, Required: false},
		"winrm_use_ssl":                        &hcldec.AttrSpec{Name: "winrm_use_ssl", Type: cty.Bool, Required: false},
		"winrm_insecure":                       &hcldec.AttrSpec{Name: "winrm_insecure", Type: cty.Bool, Required: false},
		"winrm_use_ntlm":                       &hcldec.AttrSpec{Name: "winrm_use_ntlm", Type: cty.Bool, Required: false},
		"ssh_interface":                        &hcldec.AttrSpec{Name: "ssh_interface", Type: cty.String, Required: false},
		"omi_block_device_mappings":            &hcldec.BlockListSpec{TypeName: "omi_block_device_mappings", Nested: hcldec.ObjectSpec((*common.FlatBlockDevice)(nil).HCL2Spec())},
		"launch_block_device_mappings":         &hcldec.BlockListSpec{TypeName: "launch_block_device_mappings", Nested: hcldec.ObjectSpec((*common.FlatBlockDevice)(nil).HCL2Spec())},
		"omi_name":                             &hcldec.AttrSpec{Name: "omi_name", Type: cty.String, Required: false},
		"omi_description":                      &hcldec.AttrSpec{Name: "omi_description", Type: cty.String, Required: false},
		"omi_account_ids":                      &hcldec.AttrSpec{Name: "omi_account_ids", Type: cty.List(cty.String), Required: false},
		"omi_groups":                           &hcldec.AttrSpec{Name: "omi_groups", Type: cty.List(cty.String), Required: false},
		"omi_product_codes":                    &hcldec.AttrSpec{Name: "omi_product_codes", Type: cty.List(cty.String), Required: false},
		"omi_regions":                          &hcldec.AttrSpec{Name: "omi_regions", Type: cty.List(cty.String), Required: false},
		"omi_boot_modes":                       &hcldec.AttrSpec{Name: "omi_boot_modes", Type: cty.List(cty.String), Required: false},
		"skip_region_validation":               &hcldec.AttrSpec{Name: "skip_region_validation", Type: cty.Bool, Required: false},
		"tags":                                 &hcldec.AttrSpec{Name: "tags", Type: cty.Map(cty.String), Required: false},
		"force_deregister":                     &hcldec.AttrSpec{Name: "force_deregister", Type: cty.Bool, Required: false},
		"force_delete_snapshot":                &hcldec.AttrSpec{Name: "force_delete_snapshot", Type: cty.Bool, Required: false},
		"snapshot_tags":                        &hcldec.AttrSpec{Name: "snapshot_tags", Type: cty.Map(cty.String), Required: false},
		"snapshot_account_ids":                 &hcldec.AttrSpec{Name: "snapshot_account_ids", Type: cty.List(cty.String), Required: false},
		"global_permission":                    &hcldec.AttrSpec{Name: "global_permission", Type: cty.Bool, Required: false},
		"product_codes":                        &hcldec.AttrSpec{Name: "product_codes", Type: cty.List(cty.String), Required: false},
		"root_device_name":                     &hcldec.AttrSpec{Name: "root_device_name", Type: cty.String, Required: false},
		"omi_root_device":                      &hcldec.BlockSpec{TypeName: "omi_root_device", Nested: hcldec.ObjectSpec((*FlatRootBlockDevice)(nil).HCL2Spec())},
		"run_volume_tags":                      &hcldec.AttrSpec{Name: "run_volume_tags", Type: cty.Map(cty.String), Required: false},
	}
	return s
}

// FlatRootBlockDevice is an auto-generated flat version of RootBlockDevice.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatRootBlockDevice struct {
	SourceDeviceName   *string `mapstructure:"source_device_name" cty:"source_device_name" hcl:"source_device_name"`
	DeviceName         *string `mapstructure:"device_name" cty:"device_name" hcl:"device_name"`
	DeleteOnVmDeletion *bool   `mapstructure:"delete_on_vm_deletion" cty:"delete_on_vm_deletion" hcl:"delete_on_vm_deletion"`
	IOPS               *int64  `mapstructure:"iops" cty:"iops" hcl:"iops"`
	VolumeType         *string `mapstructure:"volume_type" cty:"volume_type" hcl:"volume_type"`
	VolumeSize         *int64  `mapstructure:"volume_size" cty:"volume_size" hcl:"volume_size"`
}

// FlatMapstructure returns a new FlatRootBlockDevice.
// FlatRootBlockDevice is an auto-generated flat version of RootBlockDevice.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*RootBlockDevice) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatRootBlockDevice)
}

// HCL2Spec returns the hcl spec of a RootBlockDevice.
// This spec is used by HCL to read the fields of RootBlockDevice.
// The decoded values from this spec will then be applied to a FlatRootBlockDevice.
func (*FlatRootBlockDevice) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"source_device_name":    &hcldec.AttrSpec{Name: "source_device_name", Type: cty.String, Required: false},
		"device_name":           &hcldec.AttrSpec{Name: "device_name", Type: cty.String, Required: false},
		"delete_on_vm_deletion": &hcldec.AttrSpec{Name: "delete_on_vm_deletion", Type: cty.Bool, Required: false},
		"iops":                  &hcldec.AttrSpec{Name: "iops", Type: cty.Number, Required: false},
		"volume_type":           &hcldec.AttrSpec{Name: "volume_type", Type: cty.String, Required: false},
		"volume_size":           &hcldec.AttrSpec{Name: "volume_size", Type: cty.Number, Required: false},
	}
	return s
}
