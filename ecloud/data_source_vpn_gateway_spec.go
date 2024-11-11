package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNGatewaySpecification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPNGatewaySpecificationRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_specification_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNGatewaySpecificationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_gateway_specification_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	specs, err := service.GetVPNGatewaySpecifications(params)
	if err != nil {
		return diag.Errorf("Error retrieving VPN gateway specifications: %s", err)
	}

	if len(specs) < 1 {
		return diag.Errorf("No VPN gateway specifications found with provided arguments")
	}

	if len(specs) > 1 {
		return diag.Errorf("More than 1 VPN gateway specification found with provided arguments")
	}

	d.SetId(specs[0].ID)
	d.Set("name", specs[0].Name)
	d.Set("description", specs[0].Description)

	return nil
}
