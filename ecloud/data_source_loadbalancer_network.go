package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceLoadBalancerNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLoadBalancerNetworkRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_network_id": {
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
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceLoadBalancerNetworkRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if lbNetworkID, ok := d.GetOk("load_balancer_network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{lbNetworkID.(string)}))
	}
	if lbID, ok := d.GetOk("load_balancer_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("load_balancer_id", connection.EQOperator, []string{lbID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if networkID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{networkID.(string)}))
	}

	lbNetworks, err := service.GetLoadBalancerNetworks(params)
	if err != nil {
		return fmt.Errorf("Error retrieving loadbalancer networks: %s", err)
	}

	if len(lbNetworks) < 1 {
		return errors.New("No loadbalancer networks found with provided arguments")
	}

	if len(lbNetworks) > 1 {
		return errors.New("More than 1 loadbalancer network found with provided arguments")
	}

	d.SetId(lbNetworks[0].ID)
	d.Set("name", lbNetworks[0].Name)
	d.Set("config_id", lbNetworks[0].NetworkID)
	d.Set("load_balancer_id", lbNetworks[0].LoadBalancerID)
	d.Set("network_id", lbNetworks[0].NetworkID)

	return nil
}
