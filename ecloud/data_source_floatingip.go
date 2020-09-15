package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceFloatingIP() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFloatingIPRead,

		Schema: map[string]*schema.Schema{
			"floatingip_id": {
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceFloatingIPRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("floatingip_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}

	azs, err := service.GetFloatingIPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active floating IPs: %s", err)
	}

	if len(azs) < 1 {
		return errors.New("No floating IPs found with provided arguments")
	}

	if len(azs) > 1 {
		return errors.New("More than 1 floating IP found with provided arguments")
	}

	d.SetId(azs[0].ID)

	return nil
}
