package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPNEndpointRead,

		Schema: map[string]*schema.Schema{
			"vpn_endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceVPNEndpointRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_endpoint_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if vpnServiceID, ok := d.GetOk("vpn_service_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpn_service_id", connection.EQOperator, []string{vpnServiceID.(string)}))
	}
	if floatingIPID, ok := d.GetOk("floating_ip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("floating_ip_id", connection.EQOperator, []string{floatingIPID.(string)}))
	}

	vpnEndpoints, err := service.GetVPNEndpoints(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active VPN services: %s", err)
	}

	if len(vpnEndpoints) < 1 {
		return errors.New("No VPN services found with provided arguments")
	}

	if len(vpnEndpoints) > 1 {
		return errors.New("More than 1 VPN service found with provided arguments")
	}

	d.SetId(vpnEndpoints[0].ID)
	d.Set("name", vpnEndpoints[0].Name)
	d.Set("vpn_service_id", vpnEndpoints[0].VPNServiceID)
	d.Set("floating_ip_id", vpnEndpoints[0].FloatingIPID)

	return nil
}
