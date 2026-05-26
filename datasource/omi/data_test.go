package omi_test

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/outscale/packer-plugin-outscale/builder/common"
	"github.com/outscale/packer-plugin-outscale/datasource/omi"
)

func TestDatasourceConfigure_FilterBlank(t *testing.T) {
	datasource := omi.Datasource{
		Config: omi.Config{
			OmiFilterOptions: common.OmiFilterOptions{},
		},
	}
	if err := datasource.Configure(nil); err == nil {
		t.Fatalf("Should error if filters map is empty or not specified")
	}
}

func TestRunConfigPrepare_SourceOmiFilterOwnersBlank(t *testing.T) {
	datasource := omi.Datasource{
		Config: omi.Config{
			OmiFilterOptions: common.OmiFilterOptions{
				NameValueFilter: config.NameValueFilter{
					Filters: map[string]string{"foo": "bar"},
				},
			},
		},
	}
	if err := datasource.Configure(nil); err == nil {
		t.Fatalf("Should error if Owners is not specified)")
	}
}

func TestRunConfigPrepare_SourceOmiFilterGood(t *testing.T) {
	filter_key := "name"
	filter_value := "foo"
	datasource := omi.Datasource{
		Config: omi.Config{
			OmiFilterOptions: common.OmiFilterOptions{
				NameValueFilter: config.NameValueFilter{
					Filters: map[string]string{filter_key: filter_value},
				},
				Owners: []string{"1234567"},
			},
		},
	}
	if err := datasource.Configure(nil); err != nil {
		t.Fatalf("err: %s", err)
	}
}
