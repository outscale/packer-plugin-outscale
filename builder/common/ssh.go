package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// modified in tests
var sshHostSleepDuration = time.Second

func isNotEmpty(v *string) bool {
	var sv string
	if v != nil {
		sv = *v
	}
	return sv != ""
}

// SSHHost returns a function that can be given to the SSH communicator
// for determining the SSH address based on the vm DNS name.
func SSHHost(conn *OscClient, sshInterface string) func(multistep.StateBag) (string, error) {
	return func(state multistep.StateBag) (string, error) {
		ctx := context.Background()

		const tries = 2
		// <= with current structure to check result of describing `tries` times
		for j := 0; j <= tries; j++ {
			var host string
			i := state.Get("vm").(oscgo.Vm)

			if sshInterface != "" {
				switch sshInterface {
				case "public_ip":
					if isNotEmpty(i.PublicIp) {
						host = *i.PublicIp
					}
				case "public_dns":
					if isNotEmpty(i.PublicDnsName) {
						host = *i.PublicDnsName
					}
				case "private_ip":
					if i.PrivateIp != "" {
						host = i.PrivateIp
					}
				case "private_dns":
					if isNotEmpty(i.PrivateDnsName) {
						host = *i.PrivateDnsName
					}
				default:
					panic(fmt.Sprintf("Unknown interface type: %s", sshInterface))
				}
			} else if isNotEmpty(i.NetId) {
				if isNotEmpty(i.PublicIp) {
					host = *i.PublicIp
				} else if i.PrivateIp != "" {
					host = i.PrivateIp
				}
			} else if isNotEmpty(i.PublicDnsName) {
				host = *i.PublicDnsName
			}

			if host != "" {
				return host, nil
			}

			r, err := conn.ReadVms(ctx, oscgo.ReadVmsRequest{
				Filters: &oscgo.FiltersVm{
					VmIds: &[]string{i.VmId},
				},
			})
			if err != nil {
				return "", err
			}

			if len(*r.Vms) == 0 {
				return "", fmt.Errorf("vm not found: %s", i.VmId)
			}

			state.Put("vm", (*r.Vms)[0])
			time.Sleep(sshHostSleepDuration)
		}

		return "", errors.New("couldn't determine address for vm")
	}
}

// SSHHost returns a function that can be given to the SSH communicator
// for determining the SSH address based on the vm DNS name.
func OscSSHHost(conn *OscClient, sshInterface string) func(multistep.StateBag) (string, error) {
	return SSHHost(conn, sshInterface) // WTF ??
}
