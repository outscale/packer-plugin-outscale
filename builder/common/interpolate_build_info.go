package common

import (
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

type BuildInfoTemplate struct {
	BuildRegion   string
	SourceOMI     string
	SourceOMIName string
	SourceOMITags map[string]string
}

func extractBuildInfo(region string, state multistep.StateBag) *BuildInfoTemplate {
	rawSourceOMI, hasSourceOMI := state.GetOk("source_image")
	if !hasSourceOMI {
		return &BuildInfoTemplate{
			BuildRegion: region,
		}
	}

	sourceOMI := rawSourceOMI.(oscgo.Image)
	sourceOMITags := make(map[string]string, len(*sourceOMI.Tags))
	for _, tag := range *sourceOMI.Tags {
		sourceOMITags[tag.Key] = tag.Value
	}

	return &BuildInfoTemplate{
		BuildRegion:   region,
		SourceOMI:     *sourceOMI.ImageId,
		SourceOMIName: *sourceOMI.ImageName,
		SourceOMITags: sourceOMITags,
	}
}
