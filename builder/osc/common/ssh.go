package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

var (
	// modified in tests
	sshHostSleepDuration = time.Second
)

// SSHHost returns a function that can be given to the SSH communicator
// for determining the SSH address based on the vm DNS name.
func SSHHost(conn *OscClient, sshInterface string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		const tries = 2
		// <= with current structure to check result of describing `tries` times
		for j := 0; j <= tries; j++ {
			var host string
			i := state.Get("vm").(oscgo.Vm)

			if sshInterface != "" {
				switch sshInterface {
				case "public_ip":
					if i.GetPublicIp() != "" {
						host = i.GetPublicIp()
					}
				case "public_dns":
					if i.GetPublicDnsName() != "" {
						host = i.GetPublicDnsName()
					}
				case "private_ip":
					if i.GetPrivateIp() != "" {
						host = i.GetPrivateIp()
					}
				case "private_dns":
					if i.GetPrivateDnsName() != "" {
						host = i.GetPrivateDnsName()
					}
				default:
					panic(fmt.Sprintf("Unknown interface type: %s", sshInterface))
				}
			} else if i.GetNetId() != "" {
				if i.GetPublicIp() != "" {
					host = i.GetPublicIp()
				} else if i.GetPrivateIp() != "" {
					host = i.GetPrivateIp()
				}
			} else if i.GetPublicDnsName() != "" {
				host = i.GetPublicDnsName()
			}

			if host != "" {
				return host, nil
			}

			r, _, err := conn.Api.VmApi.ReadVms(conn.Auth).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{i.GetVmId()},
				}}).Execute()

			if err != nil {
				return "", err
			}

			if len(r.GetVms()) == 0 {
				return "", fmt.Errorf("vm not found: %s", i.GetVmId())
			}

			state.Put("vm", r.GetVms()[0])
			time.Sleep(sshHostSleepDuration)
		}

		return "", errors.New("couldn't determine address for vm")
	}
}

// SSHHost returns a function that can be given to the SSH communicator
// for determining the SSH address based on the vm DNS name.
func OscSSHHost(conn *OscClient, sshInterface string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		const tries = 2
		// <= with current structure to check result of describing `tries` times
		for j := 0; j <= tries; j++ {
			var host string
			i := state.Get("vm").(oscgo.Vm)

			if sshInterface != "" {
				switch sshInterface {
				case "public_ip":
					if i.GetPublicIp() != "" {
						host = i.GetPublicIp()
					}
				case "public_dns":
					if i.GetPublicDnsName() != "" {
						host = i.GetPublicDnsName()
					}
				case "private_ip":
					if i.GetPrivateIp() != "" {
						host = i.GetPrivateIp()
					}
				case "private_dns":
					if i.GetPrivateDnsName() != "" {
						host = i.GetPrivateDnsName()
					}
				default:
					panic(fmt.Sprintf("Unknown interface type: %s", sshInterface))
				}
			} else if i.GetNetId() != "" {
				if i.GetPublicIp() != "" {
					host = i.GetPublicIp()
				} else if i.GetPrivateIp() != "" {
					host = i.GetPrivateIp()
				}
			} else if i.GetPublicDnsName() != "" {
				host = i.GetPublicDnsName()
			}

			if host != "" {
				return host, nil
			}

			r, _, err := conn.Api.VmApi.ReadVms(conn.Auth).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{i.GetVmId()},
				}}).Execute()

			if err != nil {
				return "", err
			}

			if len(r.GetVms()) == 0 {
				return "", fmt.Errorf("vm not found: %s", i.GetVmId())
			}

			state.Put("vm", r.GetVms()[0])
			time.Sleep(sshHostSleepDuration)
		}

		return "", errors.New("couldn't determine address for vm")
	}
}
