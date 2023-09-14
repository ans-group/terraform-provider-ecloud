package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRegionRead,

		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("region_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	regions, err := service.GetRegions(params)
	if err != nil {
		return diag.Errorf("Error retrieving active regions: %s", err)
	}

	if len(regions) < 1 {
		return diag.Errorf("No regions found with provided arguments")
	}

	if len(regions) > 1 {
		return diag.Errorf("More than 1 region found with provided arguments")
	}

	d.SetId(regions[0].ID)
	d.Set("name", regions[0].Name)

	return nil
}
