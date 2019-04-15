package ecloud

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceSolution() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSolutionRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSolutionRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	solutions, err := service.GetSolutions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving solutions: %s", err)
	}

	if len(solutions) < 1 {
		return fmt.Errorf("No solution found with name [%s]", name)
	}

	if len(solutions) > 1 {
		return fmt.Errorf("More than one solution found with name [%s]", name)
	}

	d.SetId(strconv.Itoa(solutions[0].ID))
	d.Set("name", solutions[0].Name)
	d.Set("environment", solutions[0].Environment.String())

	return nil
}
