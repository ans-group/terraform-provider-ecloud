package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAvailabilityZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAvailabilityZoneRead,

		Schema: map[string]*schema.Schema{
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAvailabilityZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if regionID, ok := d.GetOk("region_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("region_id", connection.EQOperator, []string{regionID.(string)}))
	}
	if code, ok := d.GetOk("code"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("code", connection.EQOperator, []string{code.(string)}))
	}

	azs, err := service.GetAvailabilityZones(params)
	if err != nil {
		return diag.Errorf("Error retrieving active availability zones: %s", err)
	}

	if len(azs) < 1 {
		return diag.Errorf("No availability zones found with provided arguments")
	}

	if len(azs) > 1 {
		return diag.Errorf("More than 1 availability zone found with provided arguments")
	}

	d.SetId(azs[0].ID)
	d.Set("name", azs[0].Name)
	d.Set("region_id", azs[0].RegionID)
	d.Set("code", azs[0].Code)

	return nil
}
