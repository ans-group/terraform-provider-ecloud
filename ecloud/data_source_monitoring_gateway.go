package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMonitoringGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMonitoringGatewayRead,

		Schema: map[string]*schema.Schema{
			"monitoring_gateway_id": {
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
			"specification_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMonitoringGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("monitoring_gateway_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if azID, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_id", connection.EQOperator, []string{azID.(string)}))
	}

	tflog.Debug(ctx, "Retrieving monitoring gateways", map[string]interface{}{
		"parameters": params,
	})
	gateways, err := service.GetMonitoringGateways(params)
	if err != nil {
		return diag.Errorf("Error retrieving monitoring gateways: %s", err)
	}

	if len(gateways) < 1 {
		return diag.Errorf("No monitoring gateways found with provided arguments")
	}

	if len(gateways) > 1 {
		return diag.Errorf("More than 1 monitoring gateway found with provided arguments")
	}

	d.SetId(gateways[0].ID)
	d.Set("vpc_id", gateways[0].VPCID)
	d.Set("name", gateways[0].Name)
	d.Set("router_id", gateways[0].RouterID)
	d.Set("specification_id", gateways[0].SpecificationID)

	return nil
}
