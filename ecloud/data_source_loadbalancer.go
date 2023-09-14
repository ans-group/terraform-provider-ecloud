package ecloud

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"load_balancer_spec_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("load_balancer_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if azID, ok := d.GetOk("availability_zone_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("availability_zone_id", connection.EQOperator, []string{azID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if configID, ok := d.GetOk("config_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("config_id", connection.EQOperator, []string{strconv.Itoa(configID.(int))}))
	}
	if lbSpec, ok := d.GetOk("load_balancer_spec_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("load_balancer_spec_id", connection.EQOperator, []string{lbSpec.(string)}))
	}
	if lbNetworkID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{lbNetworkID.(string)}))
	}

	lbs, err := service.GetLoadBalancers(params)
	if err != nil {
		return diag.Errorf("Error retrieving loadbalancers: %s", err)
	}

	if len(lbs) < 1 {
		return diag.Errorf("No loadbalancers found with provided arguments")
	}

	if len(lbs) > 1 {
		return diag.Errorf("More than 1 loadbalancer found with provided arguments")
	}

	d.SetId(lbs[0].ID)
	d.Set("vpc_id", lbs[0].VPCID)
	d.Set("availability_zone_id", lbs[0].AvailabilityZoneID)
	d.Set("name", lbs[0].Name)
	d.Set("config_id", lbs[0].ConfigID)
	d.Set("load_balancer_spec_id", lbs[0].LoadBalancerSpecID)
	d.Set("network_id", lbs[0].NetworkID)

	return nil
}
