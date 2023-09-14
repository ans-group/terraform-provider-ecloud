package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLoadBalancerSpec() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerSpecRead,

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

func dataSourceLoadBalancerSpecRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("Error retrieving host specs: %s", err)
	}

	if len(lbSpecs) < 1 {
		return diag.Errorf("No loadbalancer specs found with provided arguments")
	}

	if len(lbSpecs) > 1 {
		return diag.Errorf("More than 1 loadbalancer spec found with provided arguments")
	}

	d.SetId(lbSpecs[0].ID)
	d.Set("name", lbSpecs[0].Name)

	return nil
}
