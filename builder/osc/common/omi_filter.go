//go:generate packer-sdc struct-markdown
package common

import (
	"fmt"
	"log"

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

	log.Printf("filters to pass to API are %#v", params.GetFilters().ImageIds)
	log.Printf("Using OMI Filters %#v", params)

	imageResp, _, err := oscconn.Api.ImageApi.ReadImages(oscconn.Auth).ReadImagesRequest(params).Execute()
	if err != nil {
		err := fmt.Errorf("Error querying OMI: %s", err)
		return nil, err
	}

	if len(imageResp.GetImages()) == 0 {
		err := fmt.Errorf("No OMI was found matching filters: %#v", params)
		return nil, err
	}

	if len(imageResp.GetImages()) > 1 && !d.MostRecent {
		err := fmt.Errorf("your query returned more than one result. Please try a more specific search, or set most_recent to true")
		return nil, err
	}

	var image oscgo.Image
	if d.MostRecent {
		image = mostRecentOscOmi(imageResp.GetImages())
	} else {
		image = imageResp.GetImages()[0]
	}
	return &image, nil
}
