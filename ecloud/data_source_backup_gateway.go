package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBackupGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupGatewayRead,

		Schema: map[string]*schema.Schema{
			"backup_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
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
			"gateway_spec_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBackupGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("backup_gateway_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{azID.(string)}))
	}

	tflog.Debug(ctx, "Retrieving backup gateways", map[string]interface{}{
		"parameters": params,
	})
	gateways, err := service.GetBackupGateways(params)
	if err != nil {
		return diag.Errorf("Error retrieving backup gateways: %s", err)
	}

	if len(gateways) < 1 {
		return diag.Errorf("No backup gateways found with provided arguments")
	}

	if len(gateways) > 1 {
		return diag.Errorf("More than 1 backup gateway found with provided arguments")
	}

	d.SetId(gateways[0].ID)
	d.Set("vpc_id", gateways[0].VPCID)
	d.Set("name", gateways[0].Name)
	d.Set("availability_zone_id", gateways[0].AvailabilityZoneID)
	d.Set("gateway_spec_id", gateways[0].GatewaySpecID)

	return nil
}
