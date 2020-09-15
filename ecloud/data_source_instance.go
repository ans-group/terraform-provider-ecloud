package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstanceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	instances, err := service.GetInstances(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active instances: %s", err)
	}

	if len(instances) < 1 {
		return fmt.Errorf("No instances found with name [%s]", name)
	}

	if len(instances) > 1 {
		return fmt.Errorf("More than 1 instance found with name [%s]", name)
	}

	d.SetId(instances[0].ID)
	d.Set("name", instances[0].Name)

	return nil
}
