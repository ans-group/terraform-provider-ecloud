package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceResourceTier() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTierRead,

		Schema: map[string]*schema.Schema{
			"resource_tier_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceResourceTierRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("resource_tier_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{azID.(string)}))
	}

	rts, err := service.GetResourceTiers(params)
	if err != nil {
		return diag.Errorf("Error retrieving active resource tiers: %s", err)
	}

	if len(rts) < 1 {
		return diag.Errorf("No resource tiers found with provided arguments")
	}

	if len(rts) > 1 {
		return diag.Errorf("More than 1 resource tier found with provided arguments")
	}

	d.SetId(rts[0].ID)
	d.Set("name", rts[0].Name)
	d.Set("availability_zone_id", rts[0].AvailabilityZoneID)

	return nil
}
