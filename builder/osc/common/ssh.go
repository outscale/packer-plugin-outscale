package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/outscale/osc-sdk-go/osc"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

var (
	// modified in tests
	sshHostSleepDuration = time.Second
)

// SSHHost returns a function that can be given to the SSH communicator
// for determining the SSH address based on the vm DNS name.
func SSHHost(conn *oscgo.APIClient, sshInterface string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		const tries = 2
		// <= with current structure to check result of describing `tries` times
		for j := 0; j <= tries; j++ {
			var host string
			i := state.Get("vm").(osc.Vm)

			if sshInterface != "" {
				switch sshInterface {
				case "public_ip":
					if i.PublicIp != "" {
						host = i.PublicIp
					}
				case "public_dns":
					if i.PublicDnsName != "" {
						host = i.PublicDnsName
					}
				case "private_ip":
					if i.PrivateIp != "" {
						host = i.PrivateIp
					}
				case "private_dns":
					if i.PrivateDnsName != "" {
						host = i.PrivateDnsName
					}
				default:
					panic(fmt.Sprintf("Unknown interface type: %s", sshInterface))
				}
			} else if i.NetId != "" {
				if i.PublicIp != "" {
					host = i.PublicIp
				} else if i.PrivateIp != "" {
					host = i.PrivateIp
				}
			} else if i.PublicDnsName != "" {
				host = i.PublicDnsName
			}

			if host != "" {
				return host, nil
			}
			/*readFilters := oscgo.FiltersVm{
				VmIds: &[]string{i.VmId},
			}*/
			//readOpts := oscgo.ReadVmsRequest{Filters: &readFilters}
			r, _, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{i.VmId},
				}}).Execute()
			//r, _, err := e.ReadVms(context.Background()).ReadVmsRequest(readOpts).Execute()
			/*r, _, err := e.ReadVms(context.Background(), &osc.ReadVmsOpts{
				ReadVmsRequest: optional.NewInterface(osc.ReadVmsRequest{
					Filters: osc.FiltersVm{
						VmIds: []string{i.VmId},
					},
				}),
			})*/
			if err != nil {
				return "", err
			}

			if len(r.GetVms()) == 0 {
				return "", fmt.Errorf("vm not found: %s", i.VmId)
			}

			state.Put("vm", r.GetVms()[0])
			time.Sleep(sshHostSleepDuration)
		}

		return "", errors.New("couldn't determine address for vm")
	}
}

// SSHHost returns a function that can be given to the SSH communicator
// for determining the SSH address based on the vm DNS name.
func OscSSHHost(conn *oscgo.APIClient, sshInterface string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		const tries = 2
		// <= with current structure to check result of describing `tries` times
		for j := 0; j <= tries; j++ {
			var host string
			i := state.Get("vm").(osc.Vm)

			if sshInterface != "" {
				switch sshInterface {
				case "public_ip":
					if i.PublicIp != "" {
						host = i.PublicIp
					}
				case "public_dns":
					if i.PublicDnsName != "" {
						host = i.PublicDnsName
					}
				case "private_ip":
					if i.PrivateIp != "" {
						host = i.PrivateIp
					}
				case "private_dns":
					if i.PrivateDnsName != "" {
						host = i.PrivateDnsName
					}
				default:
					panic(fmt.Sprintf("Unknown interface type: %s", sshInterface))
				}
			} else if i.NetId != "" {
				if i.PublicIp != "" {
					host = i.PublicIp
				} else if i.PrivateIp != "" {
					host = i.PrivateIp
				}
			} else if i.PublicDnsName != "" {
				host = i.PublicDnsName
			}

			if host != "" {
				return host, nil
			}

			r, _, err := conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{i.VmId},
				}}).Execute()

			if err != nil {
				return "", err
			}

			if len(r.GetVms()) == 0 {
				return "", fmt.Errorf("vm not found: %s", i.VmId)
			}

			state.Put("vm", r.GetVms()[0])
			time.Sleep(sshHostSleepDuration)
		}

		return "", errors.New("couldn't determine address for vm")
	}
}
