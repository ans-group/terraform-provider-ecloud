package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInstanceRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
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

func dataSourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("instance_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	instances, err := service.GetInstances(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active instances: %s", err)
	}

	if len(instances) < 1 {
		return errors.New("No instances found with provided arguments")
	}

	if len(instances) > 1 {
		return errors.New("More than 1 instance found with provided arguments")
	}

	d.SetId(instances[0].ID)
	d.Set("vpc_id", instances[0].VPCID)
	d.Set("name", instances[0].Name)

	return nil
}
