package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceVPC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPCRead,

		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
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

func dataSourceVPCRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if regionID, ok := d.GetOk("region_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("region_id", connection.EQOperator, []string{regionID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	vpcs, err := service.GetVPCs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active VPCs: %s", err)
	}

	if len(vpcs) < 1 {
		return errors.New("No VPCs found with provided arguments")
	}

	if len(vpcs) > 1 {
		return errors.New("More than 1 VPC found with provided arguments")
	}

	d.SetId(vpcs[0].ID)
	d.Set("region_id", vpcs[0].RegionID)
	d.Set("name", vpcs[0].Name)

	return nil
}
