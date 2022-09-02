package bsuvolume

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/antihax/optional"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	registryimage "github.com/hashicorp/packer-plugin-sdk/packer/registry/image"
	"github.com/outscale/osc-sdk-go/osc"
)

// map of region to list of volume IDs
type BsuVolumes map[string][]string

// map of region to list of snapshot IDs
type BsuSnapshots map[string][]string

// Artifact is an artifact implementation that contains built AMIs.
type Artifact struct {
	// A map of regions to BSU Volume IDs.
	Volumes BsuVolumes

	// A map of regions to BSU Snapshot IDs.
	Snapshots BsuSnapshots

	// BuilderId is the unique ID for the builder that created this AMI
	BuilderIdValue string

	// Client connection for performing API stuff.
	Conn *osc.APIClient

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

// returns a sorted list of region:ID pairs
func (a *Artifact) idList() []string {
	parts := make([]string, 0, len(a.Volumes))
	for region, volumeIDs := range a.Volumes {
		for _, volumeID := range volumeIDs {
			parts = append(parts, fmt.Sprintf("%s:%s", region, volumeID))
		}
	}

	sort.Strings(parts)
	return parts
}

func (a *Artifact) Id() string {
	return strings.Join(a.idList(), ",")
}

func (a *Artifact) String() string {
	return fmt.Sprintf("BSU Volumes were created:\n\n%s", strings.Join(a.idList(), "\n"))
}

func (a *Artifact) State(name string) interface{} {
	// To be able to push metadata to HCP Packer Registry, Packer will read the 'par.artifact.metadata'
	// state from artifacts to get a build's metadata.
	if name == registryimage.ArtifactStateURI {
		return a.stateHCPPackerRegistryMetadata()
	}
	return a.StateData[name]

}

func (a *Artifact) Destroy() error {
	errors := make([]error, 0)

	for region, volumeIDs := range a.Volumes {
		for _, volumeID := range volumeIDs {
			log.Printf("Deregistering Volume ID (%s) from region (%s)", volumeID, region)

			input := osc.DeleteVolumeRequest{
				VolumeId: volumeID,
			}
			if _, _, err := a.Conn.VolumeApi.DeleteVolume(context.Background(), &osc.DeleteVolumeOpts{
				DeleteVolumeRequest: optional.NewInterface(optional.NewInterface(input)),
			}); err != nil {
				errors = append(errors, err)
			}
		}
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

// stateHCPPackerRegistryMetadata will write the metadata as an hcpRegistryImage for each of the OMIs
// present in this artifact.
func (a *Artifact) stateHCPPackerRegistryMetadata() interface{} {

	images := make([]*registryimage.Image, 0, len(a.Volumes)+len(a.Snapshots))
	for region, volumeIDs := range a.Volumes {
		for _, volumeID := range volumeIDs {
			volumeID := volumeID
			image := registryimage.Image{
				ImageID:        volumeID,
				ProviderRegion: region,
				ProviderName:   "osc",
			}
			images = append(images, &image)
		}
	}

	for region, snapshotIDs := range a.Snapshots {
		for _, snapshotID := range snapshotIDs {
			snapshotID := snapshotID
			image := registryimage.Image{
				ImageID:        snapshotID,
				ProviderRegion: region,
				ProviderName:   "osc",
			}
			images = append(images, &image)
		}
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
