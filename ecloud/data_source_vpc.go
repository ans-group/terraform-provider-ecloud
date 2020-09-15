package ecloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceVPC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPCRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	name := d.Get("name").(string)

	params := connection.APIRequestParameters{}
	params.WithFilter(connection.APIRequestFiltering{
		Property: "name",
		Operator: connection.EQOperator,
		Value:    []string{name},
	})

	vpcs, err := service.GetVPCs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active VPCs: %s", err)
	}

	if len(vpcs) < 1 {
		return fmt.Errorf("No VPCs found with name [%s]", name)
	}

	if len(vpcs) > 1 {
		return fmt.Errorf("More than 1 VPC found with name [%s]", name)
	}

	d.SetId(vpcs[0].ID)
	d.Set("name", vpcs[0].Name)

	return nil
}
