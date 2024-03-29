package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLoadBalancerVip() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerVipRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_vip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceLoadBalancerVipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if lbVipID, ok := d.GetOk("load_balancer_vip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{lbVipID.(string)}))
	}
	if lbID, ok := d.GetOk("load_balancer_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("load_balancer_id", connection.EQOperator, []string{lbID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	lbNetworks, err := service.GetVIPs(params)
	if err != nil {
		return diag.Errorf("Error retrieving loadbalancer vips: %s", err)
	}

	if len(lbNetworks) < 1 {
		return diag.Errorf("No loadbalancer vips found with provided arguments")
	}

	if len(lbNetworks) > 1 {
		return diag.Errorf("More than 1 loadbalancer vip found with provided arguments")
	}

	d.SetId(lbNetworks[0].ID)
	d.Set("name", lbNetworks[0].Name)
	d.Set("load_balancer_id", lbNetworks[0].LoadBalancerID)

	return nil
}
