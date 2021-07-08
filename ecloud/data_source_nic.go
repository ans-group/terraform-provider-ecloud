package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceNic() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNicRead,

		Schema: map[string]*schema.Schema{
			"nic_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNicRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("nic_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if networkID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{networkID.(string)}))
	}
	if instanceID, ok := d.GetOk("instance_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("instance_id", connection.EQOperator, []string{instanceID.(string)}))
	}

	nics, err := service.GetNICs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active nics: %s", err)
	}

	if len(nics) < 1 {
		return errors.New("No nics found with provided arguments")
	}

	if len(nics) > 1 {
		return errors.New("More than 1 network found with provided arguments")
	}

	d.SetId(nics[0].ID)
	d.Set("network_id", nics[0].NetworkID)
	d.Set("instance_id", nics[0].InstanceID)
	d.Set("ip_address", nics[0].IPAddress)

	return nil
}
