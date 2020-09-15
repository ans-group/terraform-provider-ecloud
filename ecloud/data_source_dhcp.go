package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceDHCP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDHCPRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceDHCPRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	azID := d.Get("availability_zone_id").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "availability_zone_id",
		Operator: connection.EQOperator,
		Value:    []string{azID},
	})

	dhcps, err := service.GetDHCPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active DHCP servers/profiles: %s", err)
	}

	if len(dhcps) < 1 {
		return fmt.Errorf("No DHCP servers/profiles found with availability zone [%s]", azID)
	}

	if len(dhcps) > 1 {
		return fmt.Errorf("More than 1 DHCP server/profile found for availability zone [%s]", azID)
	}

	d.SetId(dhcps[0].ID)
	d.Set("vpc_id", dhcps[0].VPCID)
	d.Set("availability_zone_id", dhcps[0].AvailabilityZoneID)

	return nil
}
