package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	networks, err := service.GetNetworks(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active networks: %s", err)
	}

	if len(networks) < 1 {
		return fmt.Errorf("No networks found with name [%s]", name)
	}

	if len(networks) > 1 {
		return fmt.Errorf("More than 1 network found with name [%s]", name)
	}

	d.SetId(networks[0].ID)
	d.Set("name", networks[0].Name)
	d.Set("router_id", networks[0].RouterID)

	return nil
}
