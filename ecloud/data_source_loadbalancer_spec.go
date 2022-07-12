package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLoadBalancerSpec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLoadBalancerSpecRead,

		Schema: map[string]*schema.Schema{
			"loadbalancer_spec_id": {
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

func dataSourceLoadBalancerSpecRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("loadbalancer_spec_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	lbSpecs, err := service.GetLoadBalancerSpecs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving host specs: %s", err)
	}

	if len(lbSpecs) < 1 {
		return errors.New("No loadbalancer specs found with provided arguments")
	}

	if len(lbSpecs) > 1 {
		return errors.New("More than 1 loadbalancer spec found with provided arguments")
	}

	d.SetId(lbSpecs[0].ID)
	d.Set("name", lbSpecs[0].Name)

	return nil
}
