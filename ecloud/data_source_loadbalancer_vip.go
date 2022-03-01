package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceLoadBalancerVip() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLoadBalancerVipRead,

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

func dataSourceLoadBalancerVipRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error retrieving loadbalancer vips: %s", err)
	}

	if len(lbNetworks) < 1 {
		return errors.New("No loadbalancer vips found with provided arguments")
	}

	if len(lbNetworks) > 1 {
		return errors.New("More than 1 loadbalancer vip found with provided arguments")
	}

	d.SetId(lbNetworks[0].ID)
	d.Set("name", lbNetworks[0].Name)
	d.Set("load_balancer_id", lbNetworks[0].LoadBalancerID)

	return nil
}
