package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRouterRead,

		Schema: map[string]*schema.Schema{
			"router_id": {
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
			"router_throughput_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("router_id"); ok {
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
	if throughputID, ok := d.GetOk("router_throughput_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_throughput_id", connection.EQOperator, []string{throughputID.(string)}))
	}

	routers, err := service.GetRouters(params)
	if err != nil {
		return diag.Errorf("Error retrieving active routers: %s", err)
	}

	if len(routers) < 1 {
		return diag.Errorf("No routers found with provided arguments")
	}

	if len(routers) > 1 {
		return diag.Errorf("More than 1 router found with provided arguments")
	}

	d.SetId(routers[0].ID)
	d.Set("vpc_id", routers[0].VPCID)
	d.Set("name", routers[0].Name)
	d.Set("availability_zone_id", routers[0].AvailabilityZoneID)
	d.Set("router_throughput_id", routers[0].RouterThroughputID)

	return nil
}
