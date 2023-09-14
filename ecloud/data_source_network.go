package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRead,

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

func dataSourceNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error retrieving active networks: %s", err)
	}

	if len(networks) < 1 {
		return diag.Errorf("No networks found with provided arguments")
	}

	if len(networks) > 1 {
		return diag.Errorf("More than 1 network found with provided arguments")
	}

	d.SetId(networks[0].ID)
	d.Set("router_id", networks[0].RouterID)
	d.Set("subnet", networks[0].Subnet)
	d.Set("name", networks[0].Name)

	return nil
}
