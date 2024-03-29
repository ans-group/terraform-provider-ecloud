package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPNServiceRead,

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

func dataSourceVPNServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error retrieving active VPN services: %s", err)
	}

	if len(vpnServices) < 1 {
		return diag.Errorf("No VPN services found with provided arguments")
	}

	if len(vpnServices) > 1 {
		return diag.Errorf("More than 1 VPN service found with provided arguments")
	}

	d.SetId(vpnServices[0].ID)
	d.Set("name", vpnServices[0].Name)
	d.Set("vpc_id", vpnServices[0].VPCID)
	d.Set("router_id", vpnServices[0].RouterID)

	return nil
}
