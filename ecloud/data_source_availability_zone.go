package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceAvailabilityZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAvailabilityZoneRead,

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

func dataSourceAvailabilityZoneRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving active availability zones: %s", err)
	}

	if len(azs) < 1 {
		return errors.New("No availability zones found with provided arguments")
	}

	if len(azs) > 1 {
		return errors.New("More than 1 availability zone found with provided arguments")
	}

	d.SetId(azs[0].ID)
	d.Set("name", azs[0].Name)
	d.Set("region_id", azs[0].RegionID)
	d.Set("code", azs[0].Code)

	return nil
}
