package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNProfileGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNProfileGroupRead,

		Schema: map[string]*schema.Schema{
			"vpn_profile_group_id": {
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNProfileGroupRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_profile_group_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if availabilityZoneID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{availabilityZoneID.(string)}))
	}

	vpnProfileGroups, err := service.GetVPNProfileGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active VPN services: %s", err)
	}

	if len(vpnProfileGroups) < 1 {
		return errors.New("No VPN services found with provided arguments")
	}

	if len(vpnProfileGroups) > 1 {
		return errors.New("More than 1 VPN service found with provided arguments")
	}

	d.SetId(vpnProfileGroups[0].ID)
	d.Set("name", vpnProfileGroups[0].Name)
	d.Set("availability_zone_id", vpnProfileGroups[0].AvailabilityZoneID)
	d.Set("description", vpnProfileGroups[0].Description)

	return nil
}
