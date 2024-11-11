package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPNGatewayUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVPNGatewayUserRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVPNGatewayUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("vpn_gateway_user_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if vpnGatewayID, ok := d.GetOk("vpn_gateway_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpn_gateway_id", connection.EQOperator, []string{vpnGatewayID.(string)}))
	}

	users, err := service.GetVPNGatewayUsers(params)
	if err != nil {
		return diag.Errorf("Error retrieving VPN gateway users: %s", err)
	}

	if len(users) < 1 {
		return diag.Errorf("No VPN gateway users found with provided arguments")
	}

	if len(users) > 1 {
		return diag.Errorf("More than 1 VPN gateway user found with provided arguments")
	}

	d.SetId(users[0].ID)
	d.Set("name", users[0].Name)
	d.Set("vpn_gateway_id", users[0].VPNGatewayID)
	d.Set("username", users[0].Username)

	return nil
}
