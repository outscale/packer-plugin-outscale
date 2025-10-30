package common

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

// StepSourceOMIInfo extracts critical information from the source OMI
// that is used throughout the OMI creation process.
//
// Produces:
//
//	source_image *osc.Image - the source OMI info
type StepSourceOMIInfo struct {
	SourceOmi  string
	OmiFilters OmiFilterOptions
}

type imageOscSort []oscgo.Image

func (a imageOscSort) Len() int      { return len(a) }
func (a imageOscSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageOscSort) Less(i, j int) bool {
	return a[j].CreationDate.Unix() < a[i].CreationDate.Unix()
}

// Returns the most recent OMI out of a slice of images.
func mostRecentOscOmi(images []oscgo.Image) oscgo.Image {
	sortedImages := images
	sort.Sort(imageOscSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 0)
	return err == nil
}

func (s *StepSourceOMIInfo) Run(
	ctx context.Context,
	state multistep.StateBag,
) multistep.StepAction {
	oscconn := state.Get("osc").(*OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	params := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{},
	}
	if s.SourceOmi != "" {
		params.Filters.ImageIds = &[]string{s.SourceOmi}
	}

	// We have filters to apply
	if len(s.OmiFilters.Filters) > 0 {
		omiFilter := buildOSCOMIFilters(s.OmiFilters.Filters)
		params.Filters = &omiFilter
	}
	// TODO:Check if AccountIds correspond to Owners.
	if len(s.OmiFilters.Owners) > 0 {
		var oid []string
		var oali []string

		for _, o := range s.OmiFilters.Owners {
			if isNumeric(o) {
				oid = append(oid, o)
			} else {
				oali = append(oali, o)
			}
		}
		params.Filters.AccountIds = &oid
		params.Filters.AccountAliases = &oali
	}
	log.Printf("filters to pass to API are %#v", params.Filters.ImageIds)
	log.Printf("Using OMI Filters %#v", params)

	imageResp, err := oscconn.ReadImages(ctx, params)
	if err != nil {
		err := fmt.Errorf("error querying OMI: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(*imageResp.Images) == 0 {
		err := fmt.Errorf("no OMI was found matching filters: %v", params.Filters.ImageNames)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(*imageResp.Images) > 1 && !s.OmiFilters.MostRecent {
		err := errors.New(
			"your query returned more than one result. Please try a more specific search, or set most_recent to true",
		)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	var image oscgo.Image
	if s.OmiFilters.MostRecent {
		image = mostRecentOscOmi(*imageResp.Images)
	} else {
		image = (*imageResp.Images)[0]
	}

	ui.Message(fmt.Sprintf("Found Image ID: %s", image.ImageId))

	state.Put("source_image", image)
	return multistep.ActionContinue
}

func (s *StepSourceOMIInfo) Cleanup(multistep.StateBag) {}
