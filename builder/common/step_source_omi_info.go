package common

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	oscgo "github.com/outscale/osc-sdk-go/v2"
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
	itime, _ := time.Parse(time.RFC3339, *a[i].CreationDate)
	jtime, _ := time.Parse(time.RFC3339, *a[j].CreationDate)
	return itime.Unix() < jtime.Unix()
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

func (s *StepSourceOMIInfo) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	oscconn := state.Get("osc").(*OscClient)
	ui := state.Get("ui").(packersdk.Ui)

	params := oscgo.ReadImagesRequest{
		Filters: &oscgo.FiltersImage{},
	}
	if s.SourceOmi != "" {
		params.Filters.SetImageIds([]string{s.SourceOmi})
	}

	// We have filters to apply
	if len(s.OmiFilters.Filters) > 0 {
		omiFilter := buildOSCOMIFilters(s.OmiFilters.Filters)
		params.Filters = &omiFilter
	}
	//TODO:Check if AccountIds correspond to Owners.
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
	log.Printf("filters to pass to API are %#v", params.GetFilters().ImageIds)
	log.Printf("Using OMI Filters %#v", params)

	imageResp, _, err := oscconn.Api.ImageApi.ReadImages(oscconn.Auth).ReadImagesRequest(params).Execute()
	if err != nil {
		err := fmt.Errorf("Error querying OMI: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(*imageResp.Images) == 0 {
		err := fmt.Errorf("No OMI was found matching filters: %v", params.Filters.GetImageNames())
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	if len(*imageResp.Images) > 1 && !s.OmiFilters.MostRecent {
		err := fmt.Errorf("your query returned more than one result. Please try a more specific search, or set most_recent to true")
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	var image oscgo.Image
	if s.OmiFilters.MostRecent {
		image = mostRecentOscOmi(*imageResp.Images)
	} else {
		image = imageResp.GetImages()[0]
	}

	ui.Message(fmt.Sprintf("Found Image ID: %s", image.GetImageId()))

	state.Put("source_image", image)
	return multistep.ActionContinue
}

func (s *StepSourceOMIInfo) Cleanup(multistep.StateBag) {}
