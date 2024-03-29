package ecloud

import (
	"context"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkPolicyRead,

		Schema: map[string]*schema.Schema{
			"network_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_id": {
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

func dataSourceNetworkPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("network_policy_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if networkID, ok := d.GetOk("network_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_id", connection.EQOperator, []string{networkID.(string)}))
	}
	if vpcID, ok := d.GetOk("vpc_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("vpc_id", connection.EQOperator, []string{vpcID.(string)}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	tflog.Debug(ctx, "Retrieving network policies", map[string]interface{}{
		"parameters": params,
	})
	policies, err := service.GetNetworkPolicies(params)
	if err != nil {
		return diag.Errorf("Error retrieving active network policies: %s", err)
	}

	if len(policies) < 1 {
		return diag.Errorf("No network policies found with provided arguments")
	}

	if len(policies) > 1 {
		return diag.Errorf("More than 1 network policy found with provided arguments")
	}

	d.SetId(policies[0].ID)
	d.Set("network_id", policies[0].NetworkID)
	d.Set("vpc_id", policies[0].VPCID)
	d.Set("name", policies[0].Name)

	return nil
}
