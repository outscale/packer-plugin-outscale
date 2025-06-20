//go:generate packer-sdc struct-markdown
package common

import (
	"errors"
	"fmt"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func (d *OmiFilterOptions) GetOwners() []*string {
	res := make([]*string, 0, len(d.Owners))
	for _, owner := range d.Owners {
		i := owner
		res = append(res, &i)
	}
	return res
}

func (d *OmiFilterOptions) GetFilteredImage(params oscgo.ReadImagesRequest, oscconn *OscClient) (*oscgo.Image, error) {

	// We have filters to apply
	if len(d.Filters) > 0 {
		omiFilters := buildOSCOMIFilters(d.Filters)
		// chack if we can parse omi filters
		params.Filters = &omiFilters
	}
	//TODO:Check if AccountIds correspond to Owners.
	if len(d.Owners) > 0 {
		var oid []string
		var oali []string

		for _, o := range d.Owners {
			if isNumeric(o) {
				oid = append(oid, o)
			} else {
				oali = append(oali, o)
			}
		}
		params.Filters.AccountIds = &oid
		params.Filters.AccountAliases = &oali
	}

	imageResp, _, err := oscconn.Api.ImageApi.ReadImages(oscconn.Auth).ReadImagesRequest(params).Execute()
	if err != nil {
		return nil, fmt.Errorf("error querying OMI: %w", err)
	}

	if len(imageResp.GetImages()) == 0 {
		return nil, fmt.Errorf("no OMI was found matching filters: %v", params.Filters.GetImageNames())
	}

	if len(imageResp.GetImages()) > 1 && !d.MostRecent {
		return nil, errors.New("your query returned more than one result. Please try a more specific search, or set most_recent to true")
	}

	var image oscgo.Image
	if d.MostRecent {
		image = mostRecentOscOmi(imageResp.GetImages())
	} else {
		image = imageResp.GetImages()[0]
	}
	return &image, nil
}
