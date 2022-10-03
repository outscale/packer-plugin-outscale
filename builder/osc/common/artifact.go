package common

import (
	"fmt"
	"log"
	"sort"
	"strings"

	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
	oscgo "github.com/outscale/osc-sdk-go/v2"
)

// Artifact is an artifact implementation that contains built OMIs.
type Artifact struct {
	// A map of regions to OMI IDs.
	Omis map[string]string

	// BuilderId is the unique ID for the builder that created this OMI
	BuilderIdValue string

	// StateData should store data such as GeneratedData
	// to be shared with post-processors
	StateData map[string]interface{}
}

func (a *Artifact) BuilderId() string {
	return a.BuilderIdValue
}

func (*Artifact) Files() []string {
	// We have no files
	return nil
}

func (a *Artifact) Id() string {
	parts := make([]string, 0, len(a.Omis))
	for region, omiId := range a.Omis {
		parts = append(parts, fmt.Sprintf("%s:%s", region, omiId))
	}

	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func (a *Artifact) String() string {
	omiStrings := make([]string, 0, len(a.Omis))
	for region, id := range a.Omis {
		single := fmt.Sprintf("%s: %s", region, id)
		omiStrings = append(omiStrings, single)
	}

	sort.Strings(omiStrings)
	return fmt.Sprintf("OMIs were created:\n%s\n", strings.Join(omiStrings, "\n"))
}

func (a *Artifact) State(name string) interface{} {
	if _, ok := a.StateData[name]; ok {
		return a.StateData[name]
	}

	switch name {
	case "atlas.artifact.metadata":
		return a.stateAtlasMetadata()
	case registryimage.ArtifactStateURI:
		return a.stateHCPPackerRegistryMetadata()
	default:
		return nil
	}
}

func (a *Artifact) Destroy() error {
	errors := make([]error, 0)

	config := a.State("accessConfig").(*AccessConfig)

	for region, imageId := range a.Omis {
		log.Printf("Deregistering image ID (%s) from region (%s)", imageId, region)

		regionConn := config.NewOSCClientByRegion(region)

		// Get image metadata
		imageResp, _, err := regionConn.Api.ImageApi.ReadImages(regionConn.Auth).ReadImagesRequest(oscgo.ReadImagesRequest{
			Filters: &oscgo.FiltersImage{
				ImageIds: &[]string{imageId},
			},
		}).Execute()
		if err != nil {
			errors = append(errors, err)
		}
		if len(imageResp.GetImages()) == 0 {
			err := fmt.Errorf("Error retrieving details for OMI (%s), no images found", imageId)
			errors = append(errors, err)
		}

		// Deregister ami
		input := oscgo.DeleteImageRequest{
			ImageId: imageId,
		}
		_, _, err = regionConn.Api.ImageApi.DeleteImage(regionConn.Auth).DeleteImageRequest(input).Execute()
		if err != nil {
			errors = append(errors, err)
		}

		// TODO: Delete the snapshots associated with an OMI too
	}

	if len(errors) > 0 {
		if len(errors) == 1 {
			return errors[0]
		} else {
			return &packersdk.MultiError{Errors: errors}
		}
	}

	return nil
}

func (a *Artifact) stateAtlasMetadata() interface{} {
	metadata := make(map[string]string)
	for region, imageId := range a.Omis {
		k := fmt.Sprintf("region.%s", region)
		metadata[k] = imageId
	}

	return metadata
}

// stateHCPPackerRegistryMetadata will write the metadata as an hcpRegistryImage for each of the OMIs
// present in this artifact.
func (a *Artifact) stateHCPPackerRegistryMetadata() interface{} {

	f := func(k, v interface{}) (*registryimage.Image, error) {

		region, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type of key in OMIs map")
		}
		imageId, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for value in OMIs map")
		}
		image := registryimage.Image{
			ImageID:        imageId,
			ProviderRegion: region,
			ProviderName:   "osc",
		}

		return &image, nil

	}

	images, err := registryimage.FromMappedData(a.Omis, f)
	if err != nil {
		log.Printf("[TRACE] error encountered when creating HCP Packer registry image for artifact.Omis: %s", err)
		return nil
	}

	if a.StateData == nil {
		return images
	}

	data, ok := a.StateData["generated_data"].(map[string]interface{})
	if !ok {
		return images
	}

	for _, image := range images {
		image.SourceImageID = data["SourceOMI"].(string)
	}

	return images
}
