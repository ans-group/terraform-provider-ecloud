package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceAvailabilityZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAvailabilityZoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAvailabilityZoneRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	azs, err := service.GetAvailabilityZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active availability zones: %s", err)
	}

	if len(azs) < 1 {
		return fmt.Errorf("No availability zones found with name [%s]", name)
	}

	if len(azs) > 1 {
		return fmt.Errorf("More than 1 availability zone found with name [%s]", name)
	}

	d.SetId(azs[0].ID)
	d.Set("name", azs[0].Name)
	d.Set("datacentre_site_id", azs[0].DatacentreSiteID)
	d.Set("code", azs[0].Code)

	return nil
}
