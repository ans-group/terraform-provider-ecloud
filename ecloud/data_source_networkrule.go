package ecloud

import (
	"context"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkRuleRead,

		Schema: map[string]*schema.Schema{
			"network_rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_policy_id": {
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
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func dataSourceNetworkRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(ecloudservice.ECloudService)

	params := connection.APIRequestParameters{}

	if id, ok := d.GetOk("network_rule_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{id.(string)}))
	}
	if networkPolicyID, ok := d.GetOk("network_policy_id"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("network_policy_id", connection.EQOperator, []string{networkPolicyID.(string)}))
	}
	if sequence, ok := d.GetOk("sequence"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("sequence", connection.EQOperator, []string{strconv.Itoa(sequence.(int))}))
	}
	if name, ok := d.GetOk("name"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("name", connection.EQOperator, []string{name.(string)}))
	}
	if source, ok := d.GetOk("source"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("source", connection.EQOperator, []string{source.(string)}))
	}
	if destination, ok := d.GetOk("destination"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("destination", connection.EQOperator, []string{destination.(string)}))
	}
	if action, ok := d.GetOk("action"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("action", connection.EQOperator, []string{action.(string)}))
	}
	if direction, ok := d.GetOk("direction"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("direction", connection.EQOperator, []string{direction.(string)}))
	}
	if enabled, ok := d.GetOk("enabled"); ok {
		params.WithFilter(*connection.NewAPIRequestFiltering("enabled", connection.EQOperator, []string{strconv.FormatBool(enabled.(bool))}))
	}

	tflog.Debug(ctx, "Retrieving network rules", map[string]interface{}{
		"parameters": params,
	})

	rules, err := service.GetNetworkRules(params)
	if err != nil {
		return diag.Errorf("Error retrieving network rules: %s", err)
	}

	if len(rules) < 1 {
		return diag.Errorf("No network rules found with provided arguments")
	}

	if len(rules) > 1 {
		return diag.Errorf("More than 1 network rule found with provided arguments")
	}

	d.SetId(rules[0].ID)
	d.Set("network_policy_id", rules[0].NetworkPolicyID)
	d.Set("sequence", rules[0].Sequence)
	d.Set("name", rules[0].Name)
	d.Set("source", rules[0].Source)
	d.Set("destination", rules[0].Destination)
	d.Set("action", rules[0].Action)
	d.Set("direction", rules[0].Direction)
	d.Set("enabled", rules[0].Enabled)

	return nil
}
