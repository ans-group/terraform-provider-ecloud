package ecloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"solution_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func dataSourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)
	solutionID := d.Get("solution_id").(int)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	networks, err := service.GetSolutionNetworks(solutionID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving networks: %s", err)
	}

	if len(networks) < 1 {
		return fmt.Errorf("No network found with specified name [%s] and solution_id [%d]", name, solutionID)
	}

	if len(networks) > 1 {
		return fmt.Errorf("More than one network found with specified name [%s] and solution_id [%d]", name, solutionID)
	}

	d.SetId(strconv.Itoa(networks[0].ID))
	d.Set("name", networks[0].Name)

	return nil
}
