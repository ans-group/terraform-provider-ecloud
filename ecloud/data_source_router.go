package ecloud

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/connection"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func dataSourceRouter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRouterRead,

		Schema: map[string]*schema.Schema{
			"router_id": {
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

func dataSourceRouterRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	routers, err := service.GetRouters(params)
	if err != nil {
		return fmt.Errorf("Error retrieving active routers: %s", err)
	}

	if len(routers) < 1 {
		return errors.New("No routers found with provided arguments")
	}

	if len(routers) > 1 {
		return errors.New("More than 1 router found with provided arguments")
	}

	d.SetId(routers[0].ID)
	d.Set("vpc_id", routers[0].VPCID)
	d.Set("name", routers[0].Name)

	return nil
}
