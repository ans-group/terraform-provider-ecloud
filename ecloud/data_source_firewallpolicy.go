package ecloud

import (
	"context"
	"log"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFirewallPolicyRead,

		Schema: map[string]*schema.Schema{
			"firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sequence": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceFirewallPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("firewall_policy_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if routerID, ok := d.GetOk("router_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("router_id", connection.EQOperator, []string{routerID.(string)}))
	}
	if sequence, ok := d.GetOk("sequence"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("sequence", connection.EQOperator, []string{strconv.Itoa(sequence.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}

	log.Printf("[DEBUG] Retrieving firewall policies with parameters: %+v", params)
	policies, err := service.GetFirewallPolicies(params)
	if err != nil {
		return diag.Errorf("Error retrieving active firewall policies: %s", err)
	}

	if len(policies) < 1 {
		return diag.Errorf("No firewall policies found with provided arguments")
	}

	if len(policies) > 1 {
		return diag.Errorf("More than 1 firewall policy found with provided arguments")
	}

	d.SetId(policies[0].ID)
	d.Set("router_id", policies[0].RouterID)
	d.Set("sequence", policies[0].Sequence)
	d.Set("name", policies[0].Name)

	return nil
}
