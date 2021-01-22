package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet": {
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

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if routerID, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_id", connection.EQOperator, []string{routerID.(string)}))
	}
	if subnet, ok := d.GetOk("subnet"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("subnet", connection.EQOperator, []string{subnet.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	networks, err := service.GetNetworks(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active networks: %s", err)
	}

	if len(networks) < 1 {
		return errors.New("No networks found with provided arguments")
	}

	if len(networks) > 1 {
		return errors.New("More than 1 network found with provided arguments")
	}

	d.SetId(networks[0].ID)
	d.Set("router_id", networks[0].RouterID)
	d.Set("subnet", networks[0].Subnet)
	d.Set("name", networks[0].Name)

	return nil
}
