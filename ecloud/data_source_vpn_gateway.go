package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPNGatewayRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"specification_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_gateway_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if routerID, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_id", connection.EQOperator, []string{routerID.(string)}))
	}
	if specificationID, ok := d.GetOk("specification_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("specification_id", connection.EQOperator, []string{specificationID.(string)}))
	}

	gateways, err := service.GetVPNGateways(params)
	if err != nil {
		return diag.Errorf("Error retrieving VPN gateways: %s", err)
	}

	if len(gateways) < 1 {
		return diag.Errorf("No VPN gateways found with provided arguments")
	}

	if len(gateways) > 1 {
		return diag.Errorf("More than 1 VPN gateway found with provided arguments")
	}

	d.SetId(gateways[0].ID)
	d.Set("name", gateways[0].Name)
	d.Set("router_id", gateways[0].RouterID)
	d.Set("specification_id", gateways[0].SpecificationID)
	d.Set("fqdn", gateways[0].FQDN)

	return nil
}
