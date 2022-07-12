package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostRead,

		Schema: map[string]*schema.Schema{
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceHostRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("host_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	hosts, err := service.GetHosts(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active host: %s", err)
	}

	if len(hosts) < 1 {
		return errors.New("No hosts found with provided arguments")
	}

	if len(hosts) > 1 {
		return errors.New("More than 1 host found with provided arguments")
	}

	d.SetId(hosts[0].ID)
	d.Set("name", hosts[0].Name)
	d.Set("host_group_id", hosts[0].HostGroupID)

	return nil
}
