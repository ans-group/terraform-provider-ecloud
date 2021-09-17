package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceVPNService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNServiceRead,

		Schema: map[string]*schema.Schema{
			"vpn_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVPNServiceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_service_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if routerID, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_id", connection.EQOperator, []string{routerID.(string)}))
	}

	vpnServices, err := service.GetVPNServices(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active VPN services: %s", err)
	}

	if len(vpnServices) < 1 {
		return errors.New("No VPN services found with provided arguments")
	}

	if len(vpnServices) > 1 {
		return errors.New("More than 1 VPN service found with provided arguments")
	}

	d.SetId(vpnServices[0].ID)
	d.Set("name", vpnServices[0].Name)
	d.Set("vpc_id", vpnServices[0].VPCID)
	d.Set("router_id", vpnServices[0].RouterID)

	return nil
}
