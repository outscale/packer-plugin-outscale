package common

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

type TagMap map[string]string
type OSCTags []oscgo.ResourceTag

func (t OSCTags) Report(ui packersdk.Ui) {
	for _, tag := range t {
		ui.Message(fmt.Sprintf("Adding tag: \"%s\": \"%s\"",
			tag.Key, tag.Value))
	}
}

func (t TagMap) IsSet() bool {
	return len(t) > 0
}

func (t TagMap) OSCTags(ctx interpolate.Context, region string, state multistep.StateBag) (OSCTags, error) {
	var oscTags []oscgo.ResourceTag
	ctx.Data = extractBuildInfo(region, state)

	for key, value := range t {
		interpolatedKey, err := interpolate.Render(key, &ctx)
		if err != nil {
			return nil, fmt.Errorf("error processing tag: %s:%s - %w", key, value, err)
		}
		interpolatedValue, err := interpolate.Render(value, &ctx)
		if err != nil {
			return nil, fmt.Errorf("error processing tag: %s:%s - %w", key, value, err)
		}
		oscTags = append(oscTags, oscgo.ResourceTag{
			Key:   interpolatedKey,
			Value: interpolatedValue,
		})
	}
	return oscTags, nil
}

func CreateOSCTags(conn *OscClient, resourceID string, ui packersdk.Ui, tags OSCTags) error {
	tags.Report(ui)
	request := oscgo.CreateTagsRequest{
		ResourceIds: []string{resourceID},
		Tags:        tags,
	}
	_, _, err := conn.Api.TagApi.CreateTags(conn.Auth).CreateTagsRequest(request).Execute()
	return err
}
