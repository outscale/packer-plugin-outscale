package common

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

type StepSecurityGroup struct {
	CommConfig            *communicator.Config
	SecurityGroupFilter   SecurityGroupFilterOptions
	SecurityGroupIds      []string
	TemporarySGSourceCidr string

	createdGroupId string
}

func (s *StepSecurityGroup) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	var (
		ui    = state.Get("ui").(packersdk.Ui)
		conn  = state.Get("osc").(*OscClient)
		netID = state.Get("net_id").(string)
	)

	if len(s.SecurityGroupIds) > 0 {
		req := oscgo.ReadSecurityGroupsRequest{
			Filters: &oscgo.FiltersSecurityGroup{
				SecurityGroupIds: &s.SecurityGroupIds,
			},
		}
		resp, _, err := conn.Api.SecurityGroupApi.ReadSecurityGroups(conn.Auth).ReadSecurityGroupsRequest(req).Execute()

		if err != nil || len(*resp.SecurityGroups) == 0 {
			err := fmt.Errorf("Couldn't find specified security group: %s", err)
			log.Printf("[DEBUG] %s", err.Error())
			state.Put("error", err)

			return multistep.ActionHalt
		}

		log.Printf("Using specified security groups: %v", s.SecurityGroupIds)
		state.Put("securityGroupIds", s.SecurityGroupIds)

		return multistep.ActionContinue
	}

	if !s.SecurityGroupFilter.Empty() {
		filterReq := buildSecurityGroupFilters(s.SecurityGroupFilter.Filters)

		log.Printf("Using SecurityGroup Filters %v", filterReq)
		req := oscgo.ReadSecurityGroupsRequest{
			Filters: &filterReq,
		}
		resp, _, err := conn.Api.SecurityGroupApi.ReadSecurityGroups(conn.Auth).ReadSecurityGroupsRequest(req).Execute()
		if err != nil || len(*resp.SecurityGroups) == 0 {
			err := fmt.Errorf("Couldn't find security groups for filter: %s", err)
			log.Printf("[DEBUG] %s", err.Error())
			state.Put("error", err)

			return multistep.ActionHalt
		}

		securityGroupIds := []string{}
		for _, sg := range *resp.SecurityGroups {
			securityGroupIds = append(securityGroupIds, *sg.SecurityGroupId)
		}

		ui.Message(fmt.Sprintf("Found Security Group(s): %s", strings.Join(securityGroupIds, ", ")))
		state.Put("securityGroupIds", securityGroupIds)

		return multistep.ActionContinue
	}

	/* Create the group */
	groupName := fmt.Sprintf("packer_osc_%s", uuid.TimeOrderedUUID())
	ui.Say(fmt.Sprintf("Creating temporary security group for this instance: %s", groupName))

	createSGReq := oscgo.CreateSecurityGroupRequest{
		SecurityGroupName: groupName,
		Description:       "Temporary group for Packer",
	}

	if netID != "" {
		createSGReq.NetId = &netID
	}

	resp, _, err := conn.Api.SecurityGroupApi.CreateSecurityGroup(conn.Auth).CreateSecurityGroupRequest(createSGReq).Execute()

	if err != nil {
		ui.Error(err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	// Set the group ID so we can delete it later
	s.createdGroupId = *resp.SecurityGroup.SecurityGroupId

	port := int32(s.CommConfig.Port())
	if port == 0 {
		if s.CommConfig.Type != "none" {
			state.Put("error", "port must be set to a non-zero value.")
			return multistep.ActionHalt
		}
	}
	ipProtocol := "tcp"
	// Authorize the SSH access for the security group
	createSGRReq := oscgo.CreateSecurityGroupRuleRequest{
		SecurityGroupId: *resp.SecurityGroup.SecurityGroupId,
		Flow:            "Inbound",
		Rules: &[]oscgo.SecurityGroupRule{
			{
				FromPortRange: &port,
				ToPortRange:   &port,
				IpRanges:      &[]string{s.TemporarySGSourceCidr},
				IpProtocol:    &ipProtocol,
			},
		},
	}

	ui.Say(fmt.Sprintf("Authorizing access to port %d from %s in the temporary security group...", port, s.TemporarySGSourceCidr))

	_, _, err = conn.Api.SecurityGroupRuleApi.CreateSecurityGroupRule(conn.Auth).CreateSecurityGroupRuleRequest(createSGRReq).Execute()
	if err != nil {
		err := fmt.Errorf("Error authorizing temporary security group: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	// Set some state data for use in future steps
	state.Put("securityGroupIds", []string{s.createdGroupId})

	return multistep.ActionContinue
}

func (s *StepSecurityGroup) Cleanup(state multistep.StateBag) {
	if s.createdGroupId == "" {
		return
	}

	var (
		ui   = state.Get("ui").(packersdk.Ui)
		conn = state.Get("osc").(*OscClient)
	)

	ui.Say("Deleting temporary security group...")

	_, _, err := conn.Api.SecurityGroupApi.DeleteSecurityGroup(conn.Auth).DeleteSecurityGroupRequest(oscgo.DeleteSecurityGroupRequest{
		SecurityGroupId: &s.createdGroupId,
	}).Execute()
	if err != nil {
		ui.Error(fmt.Sprintf(
			"Error cleaning up security group. Please delete the group manually: %s", s.createdGroupId))
	}
}

func buildSecurityGroupFilters(input map[string]string) oscgo.FiltersSecurityGroup {
	var filters oscgo.FiltersSecurityGroup

	for k, v := range input {
		filterValue := []string{v}

		switch name := k; name {
		case "security_group_ids":
			filters.SecurityGroupIds = &filterValue
		case "security_group_names":
			filters.SecurityGroupNames = &filterValue
		case "tag_keys":
			filters.TagKeys = &filterValue
		case "tag_values":
			filters.TagValues = &filterValue
		case "tags":
			filters.Tags = &filterValue
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}

	return filters
}
