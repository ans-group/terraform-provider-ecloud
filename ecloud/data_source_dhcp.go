package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceDHCP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDHCPRead,

		Schema: map[string]*schema.Schema{
			"dhcp_id": {
				Type: schema.TypeString,
			},
			"name": {
				Type: schema.TypeString,
			},
			"filters": dataSourceAPIRequestFiltersSchema(),
		},
	}
}

func dataSourceDHCPRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("dhcp_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if azID, ok := d.GetOk("availabilityzone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availabilityzone_id", connection.EQOperator, []string{azID.(string)}))
	}
	if filters, ok := d.GetOk("filters"); ok {
		params.WithFilter(buildDataSourceAPIRequestFilters(filters.(*schema.Set))...)
	}

	dhcps, err := service.GetDHCPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active DHCP servers/profiles: %s", err)
	}

	if len(dhcps) < 1 {
		return errors.New("No DHCP servers/profiles found with provided arguments")
	}

	if len(dhcps) > 1 {
		return errors.New("More than 1 DHCP server/profile found with provided arguments")
	}

	d.SetId(dhcps[0].ID)
	d.Set("vpc_id", dhcps[0].VPCID)
	d.Set("availability_zone_id", dhcps[0].AvailabilityZoneID)

	return nil
}
