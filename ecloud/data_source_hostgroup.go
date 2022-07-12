package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostGroupRead,

		Schema: map[string]*schema.Schema{
			"host_group_id": {
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
			"host_spec_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHostGroupRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("host_group_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	hostGroups, err := service.GetHostGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active host groups: %s", err)
	}

	if len(hostGroups) < 1 {
		return errors.New("No host groups found with provided arguments")
	}

	if len(hostGroups) > 1 {
		return errors.New("More than 1 host group found with provided arguments")
	}

	d.SetId(hostGroups[0].ID)
	d.Set("vpc_id", hostGroups[0].VPCID)
	d.Set("name", hostGroups[0].Name)
	d.Set("host_spec_id", hostGroups[0].HostSpecID)

	return nil
}
